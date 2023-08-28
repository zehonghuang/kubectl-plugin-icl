package main

import "fmt"

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
