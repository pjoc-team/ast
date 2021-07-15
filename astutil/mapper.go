package astutil

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"reflect"
	"strings"
)

// FuncToken 函数
type FuncToken struct {
	Ident         *ast.Ident
	FuncName      string
	InArgNames    []string
	InTypes       []string
	InArgAndTypes []string
	OutTypes      []string
}

// BuildSignature 唯一标识
func (t FuncToken) BuildSignature() string {
	return fmt.Sprintf("%s(%s) (%s)", t.FuncName, buildArray(t.InTypes), buildArray(t.OutTypes))
}

func buildArray(s []string) string {
	sb := strings.Builder{}
	delimiter := ""
	for _, inType := range s {
		sb.WriteString(delimiter)
		sb.WriteString(inType)
		delimiter = ","
	}
	return sb.String()
}

// NotFoundError 未找到interface
var NotFoundError = fmt.Errorf("not found interface")

// ParseInterfaceTypeFuncs 解析接口类型函数
func ParseInterfaceTypeFuncs(decl *ast.GenDecl, interfaceType string) ([]*FuncToken, error) {
	if decl.Tok != token.TYPE {
		return nil, NotFoundError
	}
	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is Type.
		itype, ok := tspec.Type.(*ast.InterfaceType)
		if !ok || tspec.Name.Name != interfaceType {
			continue
		}
		return ParseInterfaceFunc(itype)
	}
	return nil, NotFoundError
}

// ParseInterfaceFunc 解析接口函数
func ParseInterfaceFunc(itype *ast.InterfaceType) ([]*FuncToken, error) {
	ft := make([]*FuncToken, 0, len(itype.Methods.List))
	for _, field := range itype.Methods.List {
		fType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		code, err := BuildFuncCode(fType)
		if err != nil {
			return nil, err
		}
		code.Ident = field.Names[0]
		code.FuncName = field.Names[0].Name // settings func name
		ft = append(ft, code)
	}
	return ft, nil
}

// BuildFuncCode 构建函数代码
func BuildFuncCode(fType *ast.FuncType) (token *FuncToken, err error) {
	token = &FuncToken{}
	for i, f := range fType.Params.List {
		fieldType, err := FieldType(f.Type)
		if err != nil {
			return nil, err
		}
		var argName string
		if len(f.Names) == 0 {
			argName = fmt.Sprintf("a%d", i)
			for isFieldNameExists(fType, argName) { // generate the unique arg name
				argName = fmt.Sprintf("%s%d", argName, i)
			}
		} else {
			argName = f.Names[0].Name
		}
		fieldToken := argName + " " + fieldType
		token.InArgNames = append(token.InArgNames, argName)
		token.InArgAndTypes = append(token.InArgAndTypes, fieldToken)
		token.InTypes = append(token.InTypes, fieldType)
	}

	for _, f := range fType.Results.List {
		fieldToken, err := FieldType(f.Type)
		if err != nil {
			return nil, err
		}
		token.OutTypes = append(token.OutTypes, fieldToken)
	}
	return
}

func isFieldNameExists(fType *ast.FuncType, argName string) bool {
	for _, field := range fType.Params.List {
		for _, name := range field.Names {
			if argName == name.Name {
				return true
			}
		}
	}
	for _, field := range fType.Results.List {
		for _, name := range field.Names {
			if argName == name.Name {
				return true
			}
		}
	}
	return false
}

type ValueSpec struct {
	Value string
	Type  string
}

func (v *ValueSpec) String() string {
	return fmt.Sprintf("Value: %v Type: %v", v.Value, v.Type)
}

func ValueType(node ast.Node) (string, error) {
	builder := &strings.Builder{}
	switch tp := node.(type) {
	case *ast.BasicLit:
		// Identifiers and basic type literals
		// (these tokens stand for classes of literals)
		// IDENT  // main
		// INT    // 12345
		// FLOAT  // 123.45
		// IMAG   // 123.45i
		// CHAR   // 'a'
		// STRING // "abc"
		// token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		switch tp.Kind {
		case token.STRING:
			return "string", nil
		case token.INT:
			return "int", nil
		case token.FLOAT:
			return "float", nil
		case token.IMAG:
			return "image", nil
		case token.CHAR:
			return "char", nil
		default:
			return "", fmt.Errorf("unknown kind: %s", tp.Kind)
		}
	case *ast.UnaryExpr:
		switch tp.Op {
		case token.AND:
			builder.WriteString("*")
		}
		xType, err := ValueType(tp.X)
		if err != nil {
			return "", err
		}
		builder.WriteString(xType)
		return builder.String(), err
	case *ast.FuncLit: // ignore
		fieldType, err := FieldType(tp.Type)
		if err != nil {
			return "", err
		}
		return fieldType, nil
	case *ast.CompositeLit:
		fieldType, err := FieldType(tp.Type)
		if err != nil {
			return "", err
		}
		return fieldType, nil
	case *ast.SelectorExpr:
		return FieldType(tp)
	case *ast.CallExpr: // _ = errors.New("")
		return "", nil
	case *ast.ParenExpr: // _ = (*regexp.Regexp)(nil)
		return "", nil
	case *ast.Ident:
		return FieldType(tp)
	default:
		log.Printf("unknown type: %v when parse value type", reflect.TypeOf(node))
		// return "", fmt.Errorf("unknown type: %v when parse value type", reflect.TypeOf(node))
	}
	return "", nil
}

func ParseValue(node ast.Node) (*ValueSpec, error) {

	v := &ValueSpec{}
	valueType, err := ValueType(node)
	if err != nil {
		return nil, err
	}
	v.Type = valueType
	sb := &strings.Builder{}
	switch tp := node.(type) {
	case *ast.BasicLit:
		v.Value = tp.Value
	case *ast.UnaryExpr:
		sb := &strings.Builder{}
		sb.WriteString(tp.Op.String())
		value, err := ParseValue(tp.X)
		if err != nil {
			return nil, err
		}
		sb.WriteString(value.Value)
		v.Value = sb.String()
	// case *ast.FuncLit: // ignore
	// 	fieldType, err := FieldType(tp.Type)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	v.Type = fieldType
	case *ast.CompositeLit:
		sb.WriteString(valueType)
		sb.WriteString("{}")
		v.Value = sb.String()
	case *ast.SelectorExpr:
		v.Value = valueType
	case *ast.CallExpr:
		s, err := ValueType(tp.Fun)
		if err != nil {
			log.Printf("failed to build value type: %T error: %v", tp.Fun, err.Error())
			return nil, err
		}
		sb.WriteString(s)
		sb.WriteString("(")
		for i, arg := range tp.Args {
			if i != 0 {
				sb.WriteString(", ")
			}
			as, err := ParseValue(arg)
			if err != nil {
				log.Printf("failed to build value type: %T error: %v", arg, err.Error())
				return nil, err
			}

			sb.WriteString(as.Value)
		}
		sb.WriteString(")")
		v.Value = sb.String()
	default:
		log.Printf("unknown type: %v when parse value", reflect.TypeOf(node))
		// return nil, fmt.Errorf("unknown type: %v when parse value", reflect.TypeOf(node))
	}
	return v, nil
}

// FieldType 根据类型生成字段类型代码
func FieldType(node ast.Node) (string, error) {
	b := strings.Builder{}
	switch tp := node.(type) {
	case *ast.ArrayType:
		b.WriteString("[]")
		childType, err := FieldType(tp.Elt)
		if err != nil {
			return "", err
		}
		b.WriteString(childType)
		return b.String(), nil
	case *ast.InterfaceType:
		return "interface{}", nil
	case *ast.Ident:
		return tp.Name, nil
	case *ast.SelectorExpr:
		prefix, ok := tp.X.(*ast.Ident)
		if !ok {
			return "", fmt.Errorf("unknown type: %T when parse X", tp.X)
		}
		return fmt.Sprintf("%s.%s", prefix.Name, tp.Sel.Name), nil
	case *ast.StarExpr:
		elem, err := FieldType(tp.X)
		if err != nil {
			return "", err
		}
		return "*" + elem, nil
	case *ast.FuncType:
		pb := strings.Builder{}
		delimiter := ""
		for _, field := range tp.Params.List {
			fieldType, err := FieldType(field.Type)
			if err != nil {
				return "", err
			}
			pb.WriteString(delimiter)
			pb.WriteString(fieldType)
			delimiter = ","
		}
		params := pb.String()

		rb := strings.Builder{}
		delimiter = ""
		for _, field := range tp.Results.List {
			fieldType, err := FieldType(field.Type)
			if err != nil {
				return "", err
			}
			rb.WriteString(delimiter)
			rb.WriteString(fieldType)
			delimiter = ","
		}
		results := rb.String()
		if len(tp.Results.List) > 1 {
			results = "(" + results + ")"
		}
		return fmt.Sprintf("func(%s) %s", params, results), nil
	case *ast.Ellipsis:
		elt, err := FieldType(tp.Elt)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("...%s", elt), nil
	case *ast.MapType:
		kt, err := FieldType(tp.Key)
		if err != nil {
			return "", err
		}
		vt, err := FieldType(tp.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", kt, vt), nil
	case *ast.StructType:
		sb := strings.Builder{}
		for _, field := range tp.Fields.List {
			sb.WriteString("\n")
			ft, err := FieldType(field)
			if err != nil {
				return "", err
			}
			if len(field.Names) > 0 {
				sb.WriteString(field.Names[0].Name)
				sb.WriteString(" ")
			}
			sb.WriteString(ft)
		}
		return fmt.Sprintf("struct {%s}", sb.String()), nil
	default:
		return "", fmt.Errorf("unknown type: %v when parse field token", reflect.TypeOf(node))
	}
}
