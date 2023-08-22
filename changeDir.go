package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func changeDirSafely(targetDir string, baseDir string) error {

	if err := validSafely(targetDir, baseDir); err == nil {
		// 进行安全的目录切换
		err = os.Chdir(targetDir)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func validSafely(targetDir string, baseDir string) error {
	// 将 targetDir 和 baseDir 解析为绝对路径
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}

	// 使用 filepath.Rel 计算 targetDir 相对于 baseDir 的相对路径
	relPath, err := filepath.Rel(absBaseDir, absTargetDir)
	if err != nil {
		return err
	}

	// 检查相对路径是否包含 ".."，如果是，则 targetDir 不在 baseDir 的子目录中
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("无法越过合法目录：%s", baseDir)
	}
	return nil
}

func workingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	return workingDir
}

func baseDirName() string {
	return filepath.Base(workingDir())
}

func dirEntryList(path string) []os.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	return files
}
