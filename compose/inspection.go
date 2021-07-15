package compose

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pjoc-team/ast/scan"
)

// findObject 查找引用的参数是否存在
func (b *ActionBuilder) findObject(name string) ([]*Object, error) {
	if name == "" {
		return nil, fmt.Errorf("name is nil")
	}

	return b.recurseFindObject(nil, name, 0)
}

// recurseFindObject 递归查找引用的参数是否存在
func (b *ActionBuilder) recurseFindObject(obj *Object, name string, depth int) ([]*Object, error) {
	if name == "" {
		return nil, fmt.Errorf("name is nil")
	}

	var objectName string
	var least string

	lastOne := false
	index := strings.Index(name, ".")
	if index < 0 {
		objectName = name
		lastOne = true
	} else {
		objectName = name[:index]
		least = name[index+1:]
	}

	os := make([]*Object, 0, 8)

	// 如果depth>0，则查找字段
	if depth > 0 {
		found, err := b.findField(obj, name)
		if err != nil {
			return nil, err
		}
		os = append(os, found)
		if lastOne {
			return os, nil
		}

		object, err := b.recurseFindObject(found, objectName, depth+1)
		if err != nil {
			return nil, err
		} else if object != nil {
			os = append(os, object...)
		}
		return os, nil
	}

	found, ok := b.requiredObjectName(objectName)
	if !ok {
		err := fmt.Errorf("not foud obj: %v", objectName)
		// log.Errorf(err.Error())
		return nil, err
	}
	os = append(os, found)
	if lastOne {
		return os, nil
	}

	object, err := b.recurseFindObject(found, least, depth+1)
	if err != nil {
		return nil, err
	} else if object != nil {
		os = append(os, object...)
	}
	return os, nil
}

func (b *ActionBuilder) findField(obj *Object, name string) (*Object, error) {
	if obj != nil {
		strct, err := b.findPath(obj.Path)
		if err != nil {
			return nil, err
		}
		if ts, ok := (strct).(*scan.Type); strct != nil && ok && ts.Type == scan.TypeStruct {
			for _, field := range ts.Fields {
				if field.Name == name {
					obj := &Object{
						Name:  name,
						Type:  field.Type,
						Doc:   field.Doc,
						Scope: 0,
						Path:  field.Path,
					}
					return obj, nil
				}
			}
		}
	}
	return nil, nil
}

func (b *ActionBuilder) requiredObjectName(objectName string) (*Object, bool) {
	if object, ok := b.CodeContext.Vars[objectName]; ok {
		return object, ok
	} else if object, ok := b.CodeContext.Predefines[objectName]; ok {
		return object, ok
	}
	return nil, false
}

// genIdentName 生成唯一标识
func (b *ActionBuilder) genIdentName(pkg string, typeName string) string {
	for i := 0; ; i++ {
		name := strcase.ToLowerCamel(typeName)
		obj, err := b.findObject(name)
		if err != nil || obj == nil {
			return name
		}
		if pkg != "" {
			name := pkg + "." + typeName
			_, err := b.findObject(name)
			if err != nil || obj == nil {
				return name
			}
		}
		name = fmt.Sprintf("%s%d", name, i)
		_, err = b.findObject(name)
		if err != nil || obj == nil {
			return name
		}
	}

}

// depth 计算深度
func depth(arg string) int {
	return strings.Count(arg, `.`)
}

// finalField 计算最后一层的变量名
func finalField(arg string) string {
	index := strings.LastIndex(arg, `.`)
	if index < 0 {
		return arg
	}
	return arg[index+1:]
}

// checkArgsType 检查两个arg的参数类型是否匹配，如果一个包含了package，另外一个没包含package，则默认按相等处理
func checkArgsType(object string, another string) bool {
	if depth(object) == depth(another) {
		return object == another
	}
	return finalField(object) == finalField(another)
}
