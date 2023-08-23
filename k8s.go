package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"runtime"
)

func getNames(s, resource string) ([]string, error) {

	var namespace = ""

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(getKubeConfigDirectory(), kubeConfigFile))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	var names []string

	switch resource {
	case "resources":
		dc := discovery.NewDiscoveryClientForConfigOrDie(config)
		apiResourceList, err := dc.ServerPreferredResources()
		if err != nil {
			return nil, err
		}
		for _, list := range apiResourceList {
			for _, resource := range list.APIResources {
				names = append(names, resource.Name)
			}
		}
	case "pods":
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, pod := range pods.Items {
			names = append(names, pod.Name)
		}
	case "namespaces":
		if v, find := ICLConfig.Completer["-n"]; find {
			names = append(names, v...)
		}
		namespaces, _ := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		for _, ns := range namespaces.Items {
			if Contains(names, ns.Name) {
				continue
			}
			names = append(names, ns.Name)
		}
	case "services":
		services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, svc := range services.Items {
			names = append(names, svc.Name)
		}
		// ... 其他资源类型
	}

	return names, nil
}

func getKubeContextList() func(string) []string {
	return func(s string) []string {
		var contexts []string
		for k, _ := range ICLConfig.KubeConfigFileMap {
			contexts = append(contexts, k)
		}
		return contexts
	}
}

const KubeConfigDirectoryWindows = "D:\\kubeconfig\\"
const KubeConfigDirectoryLinux = "/var/kubeconfig/"
const KubeConfigDirectoryDarwin = "/var/kubeconfig/"

func getKubeConfigDirectory() string {
	switch os := runtime.GOOS; os {
	case "windows":
		return KubeConfigDirectoryWindows
	case "linux":
		return KubeConfigDirectoryLinux
	case "darwin":
		return KubeConfigDirectoryDarwin
	default:
		return KubeConfigDirectoryWindows
	}
}
