package compose

// ObjectScope 对象可见域
type ObjectScope int

const (
	// Local 本地变量，只能在特定func内使用
	Local ObjectScope = iota

	// Struct 结构体变量，只能在struct内引用
	Struct

	// Package 包变量
	Package
)

