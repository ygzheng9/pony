package base

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

// SetRoot make up level folder the running folder
func SetRoot() {
	// 设定目录: 当前文件的上级目录为运行时的目录
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	fmt.Printf("exe root: %s\n", dir)

	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
