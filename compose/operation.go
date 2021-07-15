package compose

import (
	"github.com/pjoc-team/ast/scan"
)

// OperationType 执行的操作
type OperationType int

const (
	// Func 执行函数
	Func OperationType = iota

	// Assign 赋值操作，比如 channelRequest.AppID = 1
	Assign

	// Unary 二元操作，比如使用比值操作：`==`、`>`、`!=`等
	Unary

	// New 实例化对象
	New
)

// Operation 执行操作，可以是赋值操作，可以是执行函数
type Operation struct {
	// OperationType 操作类型
	Type OperationType `json:"type" yaml:"type"`

	// Func 需要执行的函数
	Func *scan.Func `json:"func" yaml:"func"`

	// UnarySymbol 二元操作符，比如 `==`、`>`、`!=`等。如果是Type=Unary，那么args的长度必须是2
	UnarySymbol string `json:"unary_symbol" yaml:"unarySymbol"`
}
