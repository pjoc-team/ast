package astutil

import (
	"go/ast"
	"reflect"
	"strings"
)

// NewUnsupportedTypeError 未支持的类型
func NewUnsupportedTypeError(t ast.Expr) error {
	return &UnsupportedTypeError{t: t}
}

// NewIllegalFieldTypeError 未支持的字段
func NewIllegalFieldTypeError(t ast.Expr) error {
	return &IllegalFieldTypeError{t: t}
}

// NewIllegalFuncTypeError 未支持的函数类型
func NewIllegalFuncTypeError(t ast.FuncType) error {
	return &IllegalFuncTypeError{t: t}
}

// UnsupportedTypeError 不支持的Type
type UnsupportedTypeError struct {
	t ast.Expr
}

func (u *UnsupportedTypeError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("unsupported type: ")
	sb.WriteString(reflect.TypeOf(u.t).String())
	return sb.String()
}

// IllegalFieldTypeError 不支持的Type
type IllegalFieldTypeError struct {
	t ast.Expr
}

func (u *IllegalFieldTypeError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("unsupported field type: ")
	sb.WriteString(reflect.TypeOf(u.t).String())
	return sb.String()
}

// IllegalFuncTypeError func类型
type IllegalFuncTypeError struct {
	t ast.FuncType
}

func (u *IllegalFuncTypeError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("unsupported func type: ")
	sb.WriteString(reflect.TypeOf(u.t).String())
	return sb.String()
}
