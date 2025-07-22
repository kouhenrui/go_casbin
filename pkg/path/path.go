package path

import (
	"path/filepath"
)

// GetProjectRoot 获取项目根目录的绝对路径
func getProjectRoot() (string, error) {
	// 获取当前工作目录
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	
	// 如果当前在 tmp 目录，回到项目根目录
	if filepath.Base(currentDir) == "tmp" {
		currentDir = filepath.Dir(currentDir)
	}
	
	// 如果当前在 cmd 目录，回到项目根目录
	if filepath.Base(currentDir) == "cmd" {
		currentDir = filepath.Dir(currentDir)
	}
	
	return currentDir, nil
}

// GetAbsolutePath 获取相对于项目根目录的绝对路径
func GetAbsolutePath(relativePath string) (string, error) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		return "", err
	}
	
	return filepath.Join(projectRoot, relativePath), nil
} 