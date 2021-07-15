package compose

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/pjoc-team/ast/scan"
)

// HasNoError  HasErrorResults的返回
var HasNoError = errors.New("not exists")

func (b *ActionBuilder) findPkg(function *scan.Func) (*scan.Pkg, error) {
	packagePath := function.Path
	if len(packagePath) < 2 {
		return nil, fmt.Errorf(
			"function: %v's package is illegal, package: %v", function.Name, packagePath,
		)
	}
	packagePath = packagePath[:len(packagePath)-2]
	log.Printf("func: %v package path: %v", function.Path, packagePath)
	for _, pkg := range b.Builder.Codes.Packages {
		obj, ok := pkg.FindPath(packagePath)
		if !ok {
			continue
		}
		rs, ok := obj.(*scan.Pkg)
		if !ok {
			return nil, fmt.Errorf(
				"found path, but type is not Pkg, real type: %v",
				reflect.TypeOf(obj),
			)
		}
		return rs, nil
	}
	return nil, nil
}

func (b *ActionBuilder) findFunc(function *scan.Func) (*scan.Func, error) {
	packagePath := function.Path
	for _, pkg := range b.Builder.Codes.Packages {
		obj, ok := pkg.FindPath(packagePath)
		if !ok {
			continue
		}
		rs, ok := obj.(*scan.Func)
		if !ok {
			return nil, fmt.Errorf(
				"found path, but type is not Func, real type: %v",
				reflect.TypeOf(obj),
			)
		}
		return rs, nil
	}
	return nil, fmt.Errorf("not found func of path: %v", function.Path)
}

// ArgNames 获取当前函数的所有参数名
func (s *Step) ArgNames() ([]string, error) {
	args := make([]string, 0, len(s.Args))
	for _, arg := range s.Args {
		args = append(args, arg.Value)
	}
	return args, nil
}

// ResultNames 获取当前函数的所有返回参数名
func (s *Step) ResultNames() ([]string, error) {
	results := make([]string, 0, len(s.Results))
	for _, result := range s.Results {
		results = append(results, result.Name)
	}
	return results, nil
}
