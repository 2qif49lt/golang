// xhbext包补充一些常见操作
package xhbext

import (
	"os"
	"path/filepath"
)

// GetProcAbsDir 获取当前进程文件的所在目录的绝对路径
func GetProcAbsDir() (string, error) {
	abs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", nil
	}
	return filepath.Dir(abs), nil
}
func IsPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
