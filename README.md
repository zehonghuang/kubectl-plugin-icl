# Kubernetes Command Line Tool

一个集成了 Helm、Kubectl、文件系统操作及自定义命令的命令行工具。

## 主要特性

1. **命令补全**：自动补全功能支持 Helm 子命令，以及 YAML 文件补全。
2. **多命令执行**：使用 `&&` 符号可以一次执行多个命令，例如：`cd .. && install xxx`。
3. **配置文件管理**：自动检测操作系统并选择合适的配置文件目录。
4. **颜色区分**：`ls` 命令可以通过颜色区分文件和文件夹。
5. **上下文切换**：使用 `sw` 命令快速切换 Kubernetes 上下文。

## 快速开始

### 安装依赖库

项目依赖 `github.com/chzyer/readline` 和 `github.com/fatih/color`，可以使用以下命令安装：

```shell
go build -o kubectl-icl
sudo mv kubectl-icl /usr/local/bin
kubectl icl <context>
```

### 配置文件示例

在以下目录中创建配置文件`config.yaml`：

- Windows: `D:\kubeconfig\`
- Linux: `/var/kubeconfig/`
- Darwin: `/var/kubeconfig/`

```yaml
remote_repository: <远程仓库URL>
local_repository: <本地仓库路径>
kube_config_map:
  context1: <上下文1>
  context2: <上下文2>
kube_config_file_map:
  context1: context1.yaml
  context2: context2.yaml
```

### 使用说明

启动命令行工具后，以下命令可用：

- `help`：显示帮助信息。
- `sw <context>`：切换到指定的 Kubernetes 上下文。
- `ls`：列出当前目录的文件和文件夹，用颜色区分。
- `k`：运行任何 Kubectl 命令。
- `helm`：运行 Helm 命令，支持 `install` 和 `upgrade` 子命令补全。
- `cd <目录>`：切换当前工作目录。

## 贡献

如果你有任何想法或建议，欢迎提交 Issues 或 Pull Requests。
