package main

import (
	"github.com/chzyer/readline"
)

func GetCompleter() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("k",
			readline.PcItem("get",
				readline.PcItemDynamic(func(s string) []string {
					return []string{"pods", "ingresses", "deployments", "services"}
				},
					readline.PcItem("-n",
						readline.PcItemDynamic(func(s string) []string {
							names, _ := getNames(s, "namespaces")
							return names
						})), // 指定命名空间
					readline.PcItem("--namespace",
						readline.PcItemDynamic(func(s string) []string {
							names, _ := getNames(s, "namespaces")
							return names
						})), // 指定命名空间
					readline.PcItem("-o", // 输出格式
						readline.PcItem("json"),
						readline.PcItem("yaml"),
						readline.PcItem("wide"),
					),
					readline.PcItem("--show-labels"), // 显示所有标签
					readline.PcItem("--watch"),       // 监控Pod的实时状态
					readline.PcItem("--selector"),
				),
			),
			readline.PcItem("describe",
				readline.PcItem("pods",
					readline.PcItem("-n"),
				),
				readline.PcItem("nodes"),
				readline.PcItem("services",
					readline.PcItem("-n",
						readline.PcItemDynamic(func(s string) []string {
							names, _ := getNames(s, "namespaces")
							return names
						}))),
			),
			readline.PcItem("create",
				readline.PcItem("-f",
					readline.PcItemDynamic(listYamlFiles())),
				readline.PcItem("-n",
					readline.PcItemDynamic(func(s string) []string {
						names, _ := getNames(s, "namespaces")
						return names
					})),
			),
			readline.PcItem("delete",
				readline.PcItem("pods",
					readline.PcItem("-n",
						readline.PcItemDynamic(func(s string) []string {
							names, _ := getNames(s, "namespaces")
							return names
						})),
				),
				readline.PcItem("nodes"),
				readline.PcItem("services",
					readline.PcItem("-n",
						readline.PcItemDynamic(func(s string) []string {
							names, _ := getNames(s, "namespaces")
							return names
						})),
				),
			),
			readline.PcItem("apply",
				readline.PcItem("-f",
					readline.PcItemDynamic(listYamlFiles())),
				readline.PcItem("-n",
					readline.PcItemDynamic(func(s string) []string {
						names, _ := getNames(s, "namespaces")
						return names
					})),
			),
			readline.PcItem("logs",
				readline.PcItem("-n"),
			),
			readline.PcItem("exec",
				readline.PcItem("-n"),
			)),

		readline.PcItem("cd",
			readline.PcItemDynamic(listFiles())),
		readline.PcItem("sw",
			readline.PcItemDynamic(getKubeContextList())),
	)
}
