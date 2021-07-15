package scan

import (
	"errors"
	"go/ast"
	"log"
	"path/filepath"
	"strings"

	"github.com/pjoc-team/ast/astutil"
	"github.com/pjoc-team/ast/path"
	"golang.org/x/tools/go/packages"
)

// TypeT type of the "type" key word
type TypeT string

const (
	// TypeArray 数组类型
	TypeArray TypeT = "array"

	// TypeStruct 结构体类型
	TypeStruct TypeT = "struct"

	// TypeFunc 函数类型
	TypeFunc TypeT = "func"

	// TypeInterface 接口类型
	TypeInterface TypeT = "interface"

	// TypeMap map类型
	TypeMap TypeT = "map"

	// TypeChan chan类型
	TypeChan TypeT = "chan"
)

// Scanner 扫描器
type Scanner struct {
	pkg     *Pkg
	options *options
}

// Pkg 包解析器
type Pkg struct {
	// 包名
	Name string

	// 导入时可用的ID
	ID string

	// 文档
	Doc string

	// 文件列表
	Files []*File

	// 扫描过程中的报错
	Errors []error

	// 路径和对象类型映射表，例如 github.com/pjoc-team/ast/scan
	// -> scan.go -> File 对应Type: File
	PathAndTypes map[string]interface{} `json:"-"`

	p *packages.Package
}

// File 源文件
type File struct {
	// Path 查找该函数的路径，一般是从 Pkg -> File -> Func/Struct
	Path Path

	// Name 文件名
	Name string

	// Imports 当前文件的导入列表，虽然是同个package，但有可能相同的导入在不同的文件是不同的name
	Imports []*Import

	// Types 当前文件定义的类型
	Types []*Type

	// Funcs 当前文件定义的函数
	Funcs []*Func

	// Source 源文件
	Source string

	// Values 变量
	Values []*Value
}

// Import 文件内的导入
type Import struct {
	// Path 查找该导入的路径，一般是从 Pkg -> File -> Import
	Path Path

	// Name 命名，可能为空
	Name string

	// Value 导入值
	Value string
}

// AliasName 命名
func (i *Import) AliasName() string {
	if i.Name != "" {
		return i.Name
	}
	index := strings.LastIndex(i.Value, ".")
	if index < 0 {
		return i.Value
	}
	return i.Value[index+1:]
}

// Value 变量类型
type Value struct {
	// Path 查找该字段的路径，一般是从 Pkg -> File -> Value
	Path Path `json:"path" yaml:"path"`

	// Name 变量名
	Name string
	// Type 变量类型
	Type string
	// Doc 文档
	Doc string
	// Value 变量值
	Value string
}

// Type 类型定义，可能是Array/Struct/Operation/Interface/Map/Chan等
type Type struct {
	// Path 查找该类型定义的路径，一般是从 Pkg -> File -> Func/Struct
	Path Path

	// Type 基础类型，例如Array/Struct/Operation/Interface/Map/Chan
	Type TypeT

	// Name 类型名称
	Name string

	// Fields 如果是struct类型，则会有多个Fields
	Fields []*Field

	// Doc 文档说明
	Doc string
}

// Func 函数
type Func struct {
	// Path 查找该函数的路径，一般是从 Pkg -> File -> Func/Struct
	Path Path `json:"path" yaml:"path"`

	// Receiver 接收者，如果函数是属于某个类型的，则会有接收者。
	// 接收者应该是在同个package，但有可能在不同的file
	Receiver *Field `json:"receiver" yaml:"receiver"`

	// Name 函数名
	Name string `json:"name" yaml:"name"`

	// Params 参数列表
	Params []*Field `json:"params" yaml:"params"`

	// Results 响应列表
	Results []*Field `json:"results" yaml:"results"`

	// Doc 文档
	Doc string `json:"doc" yaml:"doc"`
}

// Field 字段
type Field struct {
	// Path 查找该字段的路径，一般是从 Pkg -> File -> Func/Struct -> Field
	Path Path `json:"path" yaml:"path"`

	// Name 字段名
	Name string `json:"name" yaml:"name"`

	// Type 字段类型
	Type string `json:"type" yaml:"type"`

	// Doc 文档
	Doc string `json:"doc" yaml:"doc"`
}

// ScanPkg 扫描包
func ScanPkg(pkg *packages.Package, opts ...Option) (*Pkg, error) {
	p := &Pkg{
		Name:         pkg.Name,
		ID:           pkg.ID,
		PathAndTypes: make(map[string]interface{}),
		p:            pkg,
	}

	o := &options{}
	o.apply(opts...)

	s := &Scanner{
		pkg:     p,
		options: o,
	}
	for i, file := range pkg.Syntax {
		goFile := pkg.GoFiles[i]
		codeFile, errs := s.processFile(goFile, file)
		p.Errors = append(p.Errors, errs...)
		if codeFile != nil {
			p.Files = append(p.Files, codeFile)
		}
	}
	s.paths()
	return p, nil
}

// FindPath 查找路径是否存在
func (p *Pkg) FindPath(packagePath Path) (interface{}, bool) {
	pathStr := packagePath.String()
	if pathStr == p.ID {
		return p, true
	}
	object, ok := p.PathAndTypes[pathStr]
	return object, ok
}

func (s *Scanner) processFile(goFile string, file *ast.File) (codeFile *File, errs []error) {
	codeFile = &File{
		Source: path.SourcePath(goFile),
		Name:   filepath.Base(goFile),
	}
	if !s.isExported(file) {
		return nil, nil
	}
	wf := s.walk(s.pkg, codeFile, errs)
	ast.Inspect(file, wf)
	err := s.values(s.pkg, codeFile, file)
	if err != nil {
		return nil, []error{err}
	}
	return
}

func (s *Scanner) values(pkg *Pkg, codeFile *File, node *ast.File) error {
	for _, decl := range node.Decls {
		switch dt := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range dt.Specs {
				switch st := spec.(type) {
				case *ast.ValueSpec:
					var v *Value
					v, err := s.parseValue(st)
					if err != nil {
						return err
					}
					if v == nil {
						continue
					}
					codeFile.Values = append(codeFile.Values, v)
				}
			}
		}

	}
	return nil
}

func (s *Scanner) walk(pkg *Pkg, codeFile *File, errs []error) func(ast.Node) bool {
	return func(node ast.Node) bool {
		if node == nil {
			return false // 停止遍历
		}
		var err error
		defer func() {
			if err != nil {
				errs = append(errs, err)
			}
		}()

		// 判断是否需要导出
		if !s.isExported(node) {
			return true
		}

		switch n := node.(type) {
		case *ast.File:
			if n.Doc != nil {
				pkg.Doc = astutil.ParseComment(n.Doc)
			}
		case *ast.ImportSpec:
			var imports *Import
			imports, err = s.parseImport(n)
			if err != nil {
				log.Printf("failed to parse import: %#v, error: %v", n, err.Error())
				return true
			} else if imports == nil {
				return true
			}
			codeFile.Imports = append(codeFile.Imports, imports)
		case *ast.FuncDecl:
			var f *Func
			f, err = s.parseFunc(n)
			if err != nil {
				log.Printf("failed to parse func: %#v, error: %v", n, err.Error())
				return true
			} else if f == nil {
				return true
			}
			codeFile.Funcs = append(codeFile.Funcs, f)
			return true
		case *ast.TypeSpec:
			var t *Type
			t, err = s.parseType(n)
			if err != nil {
				log.Printf("failed to parse type: %#v, error: %v", n, err.Error())
				return true
			} else if t == nil {
				return true
			}
			codeFile.Types = append(codeFile.Types, t)
		default:
			return true
		}
		return true
	}
}

func (s *Scanner) isExported(node ast.Node) bool {
	if !s.options.onlyExported {
		return true
	}

	var name string
	switch n := node.(type) {
	case *ast.File:
		return ast.FileExports(n)
	case *ast.FuncDecl:
		name = n.Name.Name
	case *ast.Field:
		if len(n.Names) > 0 {
			name = n.Names[0].Name
		}
	case *ast.TypeSpec:
		name = n.Name.Name
	default:
		return true
	}
	if name == "" {
		return true
	}
	return ast.IsExported(name)
}

func (s *Scanner) parseType(ts *ast.TypeSpec) (*Type, error) {
	t := &Type{}
	t.Name = ts.Name.Name
	t.Doc = astutil.ParseComment(ts.Doc)
	switch tp := ts.Type.(type) {
	case *ast.StructType:
		fields, err := s.parseStruct(tp)
		if err != nil {
			log.Printf("failed to parse type: %v error: %v", ts.Name.Name, err.Error())
			return nil, err
		}
		t.Fields = fields
		t.Type = TypeStruct
	case *ast.Ident:
		// 可能是继承其他类型
		t.Type = TypeT(tp.Name)
	default:
		// log.Printf("unsupported type %T", ts.Type)
		return nil, astutil.NewUnsupportedTypeError(ts.Type)
	}
	return t, nil
}

func (s *Scanner) parseStruct(st *ast.StructType) ([]*Field, error) {
	fields := make([]*Field, 0, len(st.Fields.List))
	for _, field := range st.Fields.List {
		f, err := s.parseField(field)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f...)
	}
	return fields, nil
}

func (s *Scanner) parseImport(is *ast.ImportSpec) (*Import, error) {
	i := &Import{}
	if is.Name != nil {
		i.Name = is.Name.Name
	}
	i.Value = is.Path.Value
	return i, nil
}

func (s *Scanner) parseFunc(fd *ast.FuncDecl) (*Func, error) {
	codeFunc := &Func{}
	codeFunc.Name = fd.Name.Name
	codeFunc.Doc = astutil.ParseComment(fd.Doc)

	if fd.Recv != nil {
		if len(fd.Recv.List) != 1 {
			log.Fatalf("receivers of func: %s is not equals 1 ", fd.Name.Name)
		}
		for _, field := range fd.Recv.List {
			f, err := s.parseField(field)
			if err != nil {
				log.Printf("failed to parse field of func: %v error: %v", fd.Name.Name, err.Error())
				return nil, err
			}
			if len(f) != 1 {
				err = errors.New("receive is not 1")
				log.Printf(err.Error())
				return nil, err
			}
			codeFunc.Receiver = f[0]
		}
	}

	funcType := fd.Type
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			codeField, err := s.parseField(field)
			if err != nil {
				log.Printf(
					"failed to parse params field: %#v of func: %v error: %v", field,
					fd.Name.Name, err.Error(),
				)
				return nil, err
			}
			codeFunc.Params = append(codeFunc.Params, codeField...)
		}
	}

	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			codeField, err := s.parseField(field)
			if err != nil {
				log.Printf(
					"failed to parse results field: %#v of func: %v error: %v", field,
					fd.Name.Name, err.Error(),
				)
				return nil, err
			}
			codeFunc.Results = append(codeFunc.Results, codeField...)
		}
	}

	return codeFunc, nil
}

func (s *Scanner) parseField(field *ast.Field) ([]*Field, error) {
	fields := make([]*Field, 0)

	f := &Field{}
	f.Doc = astutil.ParseComment(field.Doc)
	ft, err := astutil.FieldType(field.Type)
	if err != nil {
		if len(field.Names) > 0 {
			log.Printf("failed to parse field: %v error: %v", field.Names[0].Name, err.Error())
		} else {
			log.Printf("failed to parse field: %#v error: %v", field, err.Error())
		}
		return nil, err
	}
	f.Type = ft
	if len(field.Names) > 0 {
		for _, name := range field.Names {
			nf := *f
			nf.Name = name.Name
			fields = append(fields, &nf)
		}
		return fields, err
	}
	f.Name = ft
	fields = append(fields, f)
	return fields, nil
}

func (s *Scanner) parseValue(valueSpec *ast.ValueSpec) (*Value, error) {
	v := &Value{}
	if len(valueSpec.Names) == 0 {
		log.Fatal("value names must gt 0")
	}
	v.Name = valueSpec.Names[0].Name
	if !ast.IsExported(v.Name) {
		return nil, nil
	}

	v.Doc = astutil.ParseComment(valueSpec.Doc)

	var vs *astutil.ValueSpec
	var err error
	if valueSpec.Type != nil {
		vs, err = astutil.ParseValue(valueSpec.Type)
		if err != nil {
			log.Printf("failed to parse type, value_spec: %T, err: %v", valueSpec, err)
			return nil, err
		}
	} else {
		vv := valueSpec.Values[0]
		vs, err = astutil.ParseValue(vv)
		if err != nil {
			log.Printf("failed to parse type, value_spec: %T, err: %v", valueSpec, err)
			return nil, err
		}
	}

	v.Value = vs.Value
	v.Type = vs.Type
	return v, nil
}
