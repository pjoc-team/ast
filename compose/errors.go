package compose

import "errors"

var (
	// UnsupportedFuncTypeError 不支持的类型
	UnsupportedFuncTypeError = errors.New("unsupported function type")
	// ArgsLengthNotEqualTwoError 参数长度必须等于2
	ArgsLengthNotEqualTwoError = errors.New("args' length must be 2")
	// ArgsLengthNotEqualOneError 参数长度必须等于1
	ArgsLengthNotEqualOneError = errors.New("args' length must be 2")
	// ResultsLengthNotEqualOneError 结果长度必须等于1
	ResultsLengthNotEqualOneError = errors.New("results' length must be 1")
	// ResultsLengthNotEqualTwoError 结果长度必须等于1
	ResultsLengthNotEqualTwoError = errors.New("results' length must be 2")
	// UnarySymbolNilError 二元操作符为空
	UnarySymbolNilError = errors.New("unary symbol is nil")
)
