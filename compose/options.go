package compose

import (
	"log"
)

// options 组合代码选项
type options struct {
	Vars       []*Object
	vars       map[string]*Object
	Predefines []*Object
	predefines map[string]*Object
}

// newOptions 新建默认options
func newOptions() *options {
	o := &options{
		vars:       make(map[string]*Object),
		predefines: make(map[string]*Object),
	}
	return o
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option 选项函数
type Option func(o *options)

// WithVars 设置本地变量
func WithVars(vars []*Object) Option {
	return func(o *options) {
		o.Vars = append(o.Vars, vars...)
		for _, object := range vars {
			exists, ok := o.vars[object.Name]
			if ok {
				log.Fatalf(
					"failed to put vars: %#v because key: %v is exists: %#v", vars,
					object.Name, exists,
				)
			}
			o.vars[object.Name] = object
		}
	}
}

// WithPredefines 设置预定义对象
func WithPredefines(predefines []*Object) Option {
	return func(o *options) {
		o.Predefines = append(o.Predefines, predefines...)
		for _, object := range predefines {
			exists, ok := o.predefines[object.Name]
			if ok {
				log.Fatalf(
					"failed to put object: %#v because key: %v is exists: %#v", object,
					object.Name, exists,
				)
			}
			o.predefines[object.Name] = object
		}
	}
}
