package path

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/blademainer/commons/pkg/path"
)

var (
	rootPath  = path.NeighborPath("../")
	gomodPath = path.NeighborPath("../go.mod")
)

// SourcePath 移除机器路径，只保留项目路径
func SourcePath(path string) string {
	pkg, err := GetMod(gomodPath)
	if err != nil {
		return path
	}
	if rootIndex := strings.Index(path, rootPath); rootIndex > 0 {
		return path[rootIndex:]
	} else if pkgIndex := strings.Index(path, pkg); pkgIndex > 0 {
		return path[pkgIndex:]
	}
	return path
}

// GetMod 获取当前项目的模块
func GetMod(path string) (string, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("faile read file: %v error: %v", path, err.Error())
		return "", err
	}
	module := regexp.MustCompile(`module\s+([^\s]+/)+[^\s]+`).Find(file)
	if len(module) == 0 {
		log.Printf("faile get mod, file: %v ", string(file))
		return "", nil
	}
	pkg := regexp.MustCompile(`([^\s]+/)+[^\s]+$`).Find(module)
	return string(pkg), nil
}
