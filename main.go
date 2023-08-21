package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var config *Config
var currentContext string
var kubeConfigFile string
var safelyDirectory string

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: kubectl icl <context>")
		return
	}

	if os.Args[1] == "help" {
		printHelp()
		return
	}

	configFile := filepath.Join(getKubeConfigDirectory(), "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("配置文件 %s 不存在，请确保它位于正确的位置。\n", configFile)
		return
	} else {
		config, err = loadConfig(configFile)
		if err != nil {
			fmt.Println("加载文件异常。")
			return
		}
	}

	if err := cloneOrPull(); err != nil {
		fmt.Printf("无法克隆或拉取仓库: %s\n", err)
		return
	}

	//TODO 换名字
	if err := switchContext(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          getPrompt(),
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    GetCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		continueLoop, err := handleBaseCommand(args, rl)
		if continueLoop {
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		switch args[0] {
		case "k":
			if i := index("-f", args); i > -1 && i < len(args) {
				if err := validSafely(args[i+1], safelyDirectory); err != nil {
					fmt.Println(err)
				}
			}
			runCommand("kubectl", append(args[1:], []string{"--kubeconfig", filepath.Join(getKubeConfigDirectory(), kubeConfigFile)}...)...)
			break
		case "helm":
			runCommand("helm", append(args[1:], []string{"--kubeconfig", filepath.Join(getKubeConfigDirectory(), kubeConfigFile)}...)...)
		default:
			fmt.Printf("无效命令: %s\n", args[0])
		}

	}
}

func runCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		//fmt.Println("命令执行出错:", err)
	}
}

func getKubeConfigFile(inputContext string) error {
	if v, ok := config.KubeConfigFileMap[inputContext]; !ok {
		return fmt.Errorf("找不到对应 %s 配置.\n", inputContext)
	} else {
		kubeConfigFile = v + ".yaml"
	}
	return nil
}

func handleBaseCommand(args []string, rl *readline.Instance) (bool, error) {
	switch args[0] {
	case "help":
		printHelp()
		return true, nil
	case "sw":
		if len(args) > 1 {
			if err := switchContext(args[1]); err != nil {
				return true, err
			}
			rl.SetPrompt(getPrompt())
			fmt.Printf("上下文已切换至 %s\n", args[1])
		} else {
			fmt.Println("请提供上下文名称")
		}
		return true, nil
	case "ls":
		return true, printfFilesOnWorkingDir()
	case "cd":
		isSafely := changeDirSafely(args[1], safelyDirectory)
		if isSafely != nil {
			return true, isSafely
		}
		rl.SetPrompt(getPrompt())
		return true, nil
	case "pull":
		return true, cloneOrPull()
	default:
		return false, nil
	}
}

func cloneOrPull() error {
	repoPath := config.LocalRepository
	cmdArgs := []string{"git"}
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		cmdArgs = append(cmdArgs, "clone", config.RemoteRepository, repoPath)
	} else {
		cmdArgs = append(cmdArgs, "-C", repoPath, "pull")
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func switchContext(inputContext string) error {

	to := filepath.Join(config.LocalRepository, config.KubeConfigFileMap[inputContext])

	if err := getKubeConfigFile(inputContext); err != nil {
		return err
	}
	tgDir := filepath.Join(to, baseDirName())
	if _, err := os.Stat(tgDir); os.IsNotExist(err) {
		tgDir = to
	}
	if err := os.Chdir(tgDir); err != nil {
		return fmt.Errorf("无法切换到目录 %s: %w", tgDir, err)
	}
	currentContext = inputContext
	safelyDirectory = to
	return printFiles(tgDir)
}

func printfFilesOnWorkingDir() error {
	return printFiles(workingDir())
}

func printFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("无法读取目录 %s: %w", dir, err)
	}
	fmt.Printf("项目 in %s:\n", dir)
	for _, file := range files {
		info, _ := file.Info()
		green := color.New(color.FgGreen).SprintFunc()
		if file.IsDir() {
			fmt.Printf("  %s \t\t%s\n", green(file.Name()), info.ModTime().Format(time.RFC822))
		} else {
			fmt.Printf("  %s\t\t%s\n", file.Name(), info.ModTime().Format(time.RFC822))
		}

	}
	return nil
}

func listFiles() func(string) []string {
	return func(line string) []string {
		files, err := os.ReadDir(workingDir())
		if err != nil {
			return nil
		}

		var names []string
		for _, file := range files {
			if file.IsDir() {
				names = append(names, file.Name()) // 为目录添加 "./" 前缀
			} else {
				names = append(names, file.Name())
			}
		}
		return names
	}
}

func listYamlFiles() func(string) []string {
	return func(line string) []string {
		files, err := os.ReadDir(workingDir())
		if err != nil {
			return nil
		}
		var names []string
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
				names = append(names, file.Name())
			}
		}
		return names
	}
}

func getPrompt() string {
	hiRed := color.New(color.FgRed).SprintFunc()

	return fmt.Sprintf("%s %s > ", hiRed(config.KubeConfigFileMap[currentContext]), hiRed("("+filepath.Base(workingDir())+")"))
}

type Config struct {
	RemoteRepository  string            `yaml:"remote_repository"`
	LocalRepository   string            `yaml:"local_repository"`
	KubeConfigMap     map[string]string `yaml:"kube_config_map"`
	KubeConfigFileMap map[string]string `yaml:"kube_config_file_map"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func index[T any](v T, array []T) int {
	if n := len(array); array != nil && n != 0 {
		i := 0
		for !reflect.DeepEqual(v, array[i]) {
			i++
		}
		if i != n {
			return i
		}
	}
	return -1
}

func printHelp() {
	fmt.Println("kubectl icl - 交互式 Kubernetes 上下文管理工具")
	fmt.Println("\n可用命令:")
	fmt.Println("  help                   打印帮助信息")
	fmt.Println("  sw <context>           切换 Kubernetes 上下文")
	fmt.Println("  ls                     列出当前目录下的项目")
	fmt.Println("  cd <directory>         更改目录，无法越过当前上下文的根目录")
	fmt.Println("  pull                   从远程仓库拉取最新版本")
	fmt.Println("  k <kubectl commands>   执行 kubectl 命令")
	fmt.Println("  helm <helm commands>   执行 helm 命令")
	fmt.Println("\n配置文件:")
	fmt.Printf("  配置文件名为 'config.yaml'，应位于目录: %s\n", getKubeConfigDirectory())
	fmt.Println("  配置文件结构:")
	fmt.Println("    remote_repository: <远程仓库URL>")
	fmt.Println("    local_repository:  <本地仓库路径>")
	//fmt.Println("    kube_config_map:   <上下文到Kubeconfig的映射>")
	fmt.Println("    kube_config_file_map: <上下文到Kubeconfig文件的映射>")

	fmt.Println("配置文件示例:")
	fmt.Println("```yaml")
	fmt.Println("remote_repository: git@github.com:user/repo.git")
	fmt.Println("local_repository: /path/to/local/repo")
	//fmt.Println("kube_config_map:")
	//fmt.Println("  cn-hz: cn-hangzhou")
	//fmt.Println("  us-sv: us-siliconvalley")
	fmt.Printf("## 集群cn-hangzhou,us-siliconvalley的配置文件均放在: %s\n", getKubeConfigDirectory())
	fmt.Println("kube_config_file_map:")
	fmt.Println("  cn-hz: cn-hangzhou")
	fmt.Println("  us-sv: us-siliconvalley")
	fmt.Println("```")

	fmt.Println("\n命令示例:")
	fmt.Println("  kubectl icl my-cluster         切换到名为 'my-cluster' 的上下文")
	fmt.Println("  k get pods -n kube-system      获取 'kube-system' 命名空间中的 pods")
	fmt.Println("  cd my-directory                更改到 'my-directory' 目录")
	fmt.Println("\n更多信息，请访问 <网址或文档链接>")
}
