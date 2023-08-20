package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func changeDirSafely(targetDir string, baseDir string) (bool, error) {
	// 将 targetDir 和 baseDir 解析为绝对路径
	absContextDir, err := filepath.Abs(targetDir)
	if err != nil {
		return false, err
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false, err
	}

	// 使用 filepath.Rel 计算 targetDir 相对于 baseDir 的相对路径
	relPath, err := filepath.Rel(absBaseDir, absContextDir)
	if err != nil {
		return false, err
	}

	// 检查相对路径是否包含 ".."，如果是，则 targetDir 不在 baseDir 的子目录中
	if strings.HasPrefix(relPath, "..") {
		return false, nil
	}

	// 进行安全的目录切换
	err = os.Chdir(absContextDir)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func workingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("程序异常")
		panic(1)
	}
	return workingDir
}

func baseDirName() string {
	return filepath.Base(workingDir())
}
