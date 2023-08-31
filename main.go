package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var ICLConfig *Config
var currentContext string
var kubeConfigFile string
var safelyDirectory string
var allowCommand = []string{"pwd", "clear", "ps", "ls"}

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
		ICLConfig, err = loadConfig(configFile)
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
	if err := switchContext(os.Args[1], ""); err != nil {
		fmt.Println(err)
		return
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt: getPrompt(),
		HistoryFile: func() string {
			switch os := runtime.GOOS; os {
			case "windows":
				return "D:\\readline.tmp"
			case "linux":
				return "/tmp/readline.tmp"
			case "darwin":
				return "/tmp/readline.tmp"
			default:
				return "D:\\readline.tmp"
			}
		}(),
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

		//if strings.Index(line, "&&") > -1 {
		//	commands := strings.Split(line, "&&")
		//	for _, command := range commands {
		//
		//	}
		//}

		//swCommand := ""
		//kCommand := ""
		//
		//andIndex := strings.Index(line, "&&")
		//if andIndex > 0 && andIndex <  {
		//
		//}
		//commands := strings.Split(line, "&&")
		//if len(commands) > 2 || true {
		//	fmt.Println("'&&'指令仅且支持'k'和'sw'同时用.")
		//	continue
		//}
		//
		//for _, command := range commands {
		//	command = strings.TrimSpace(command)
		//	if strings.HasPrefix(command, "sw ") {
		//		swCommand = command
		//	} else if strings.HasPrefix(command, "k ") {
		//		kCommand = command
		//	} else {
		//		fmt.Println("'&&'指令仅且支持'k'和'sw'同时用.")
		//		continue
		//	}
		//}

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
			if i := indexOf("-f", args); i > -1 && i < len(args) {
				if err := validSafely(args[i+1], safelyDirectory); err != nil {
					fmt.Println(err)
					break
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

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", append([]string{"-Command"}, append([]string{name}, arg...)...)...)
	} else {
		cmd = exec.Command(name, arg...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("命令执行出错:", err)
	}
}

func getKubeConfigFile(inputContext string) error {
	if v, ok := ICLConfig.KubeConfigFileMap[inputContext]; !ok {
		return fmt.Errorf("找不到对应 %s 配置", inputContext)
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
			tgDir := ""
			if len(args) == 3 {
				tgDir = args[2]
			}
			if err := switchContext(args[1], tgDir); err != nil {
				return true, err
			}
			rl.SetPrompt(getPrompt())
			fmt.Printf("上下文已切换至 %s\n", args[1])
		} else {
			fmt.Println("请提供上下文名称")
		}
		return true, nil
	//case "ls":
	//	return true, printfFilesOnWorkingDir()
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
		if Contains(allowCommand, args[0]) {
			runCommand(args[0])
			return true, nil
		}
		return false, nil
	}
}

func cloneOrPull() error {
	repoPath := ICLConfig.LocalRepository
	cmdArgs := []string{"git"}
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		cmdArgs = append(cmdArgs, "clone", ICLConfig.RemoteRepository, repoPath)
	} else {
		cmdArgs = append(cmdArgs, "-C", repoPath, "pull")
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func switchContext(inputContext, dir string) error {

	to := filepath.Join(ICLConfig.LocalRepository, ICLConfig.KubeConfigFileMap[inputContext])

	if err := getKubeConfigFile(inputContext); err != nil {
		return err
	}
	tgDir := filepath.Join(to, baseDirName())
	if len(dir) > 0 {
		tgDir = filepath.Join(to, dir)
	}
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
		s := ""
		line = strings.TrimPrefix(line, "cd ")
		if strings.HasPrefix(line, "..") || strings.HasPrefix(line, "../") {
			s = "../"
		} else if strings.HasPrefix(line, ".") || strings.HasPrefix(line, "./") {
			s = "./"
		}
		var files []os.DirEntry

		if len(s) > 0 {
			files, _ = os.ReadDir(s)
		} else {
			files, _ = os.ReadDir(workingDir())
		}

		var names []string
		for _, file := range files {
			if file.IsDir() {
				names = append(names, s+file.Name()) // 为目录添加 "./" 前缀
			} else {
				if len(s) == 0 {
					names = append(names, file.Name())
				}
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

	return fmt.Sprintf("%s %s > ", hiRed(ICLConfig.KubeConfigFileMap[currentContext]), hiRed("("+filepath.Base(workingDir())+")"))
}

type Config struct {
	RemoteRepository  string              `yaml:"remote_repository"`
	LocalRepository   string              `yaml:"local_repository"`
	KubeConfigMap     map[string]string   `yaml:"kube_config_map"`
	KubeConfigFileMap map[string]string   `yaml:"kube_config_file_map"`
	Completer         map[string][]string `yaml:"completer"`
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
