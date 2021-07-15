package astutil

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strings"
)

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.

	owner         string
	interfaceType string
}

// Package 包
type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

// Parser 解析器
type Parser struct {
	debug   bool
	pkg     *Package
	p       *packages.Package
	Imports []string
}

// InterfaceParser 接口解析器
type InterfaceParser struct {
	Parser             *Parser
	typeName           string // Name of the constant type.
	tp                 *ast.TypeSpec
	interfaceName      string
	interfaceType      *ast.InterfaceType
	InterfaceFuncIdent map[string]*ast.Ident
	TypeFuncIdent      map[string]*ast.Ident
	FuncCodes          []*FuncToken
}

// NewParser 初始化解析器
func NewParser(pkg *packages.Package) *Parser {
	p := &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	parser := &Parser{
		pkg:     p,
		p:       pkg,
		debug:   true,
		Imports: make([]string, 0),
	}

	for i, file := range pkg.Syntax {
		p.files[i] = &File{
			file: file,
			pkg:  p,
		}
	}
	parser.parseImports()

	return parser
}

func (p *Parser) parseImports() {
	for _, file := range p.p.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			is, ok := node.(*ast.ImportSpec)
			if !ok {
				return true
			}
			sb := strings.Builder{}
			if is.Name != nil {
				sb.WriteString(is.Name.Name)
				sb.WriteString(" ")
			}
			sb.WriteString(is.Path.Value)
			p.Imports = append(p.Imports, sb.String())
			return true
		})
	}
}

// Printf 打印
func (p Parser) Printf(format string, args ...interface{}) {
	if p.debug {
		fmt.Printf(format+"\n", args...)
	}
}

func (p Parser) findInterfaceType(interfaceName string) (interfaceType *ast.InterfaceType) {
	for _, file := range p.p.Syntax {
		p.Printf("find interface: %v in file: %v", interfaceName, file.Name.Name)
		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				return true
			}
			for _, spec := range decl.Specs {
				tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is CONST.
				if tspec.Name.Name != interfaceName {
					continue
				}
				itype, ok := tspec.Type.(*ast.InterfaceType)
				if ok {
					interfaceType = itype
					return false
				}
			}
			return true
		})
	}
	return
}

func (p Parser) findType(typeName string) (tp *ast.TypeSpec) {
	for _, file := range p.p.Syntax {
		if tp != nil {
			return tp
		}
		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				// We only care about const declarations.
				return true
			}
			// The name of the type of the constants we are declaring.
			// Can change if this is a multi-element declaration.
			//typ := ""
			// Loop over the elements of the declaration. Each element is a ValueSpec:
			// a list of names possibly followed by a type, possibly followed by values.
			// If the type and value are both missing, we carry down the type (and value,
			// but the "go/types" package takes care of that).
			for _, spec := range decl.Specs {
				tspec, ok := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is Type.
				if !ok {
					continue
				}
				obj, ok := p.pkg.defs[tspec.Name]
				if !ok {
					p.Printf("not found ident: %v in this package: %v\n", tspec.Name, p.p)
					continue
				}
				p.Printf("find type: %v in this package: %v", obj.Id(), p.p)
				if tspec.Name.Name == typeName {
					tp = tspec
					return false
				}
			}
			return true
		})
	}
	return tp
}

// ParseTypeAndInterface 解析类型和接口
func (p Parser) ParseTypeAndInterface(typeName string, interfaceName string) (*InterfaceParser, error) {
	interfaceType := p.findInterfaceType(interfaceName)
	if interfaceType == nil {
		return nil, fmt.Errorf("not found interface: %v package: %v", interfaceName, p.p.Name)
	}

	tp := p.findType(typeName)
	if tp == nil {
		return nil, fmt.Errorf("not found type: %v", typeName)
	}

	ip := &InterfaceParser{
		typeName:           typeName,
		tp:                 tp,
		interfaceName:      interfaceName,
		interfaceType:      interfaceType,
		InterfaceFuncIdent: make(map[string]*ast.Ident),
		TypeFuncIdent:      make(map[string]*ast.Ident),
		Parser:             &p,
	}
	err := ip.Parser.ParseInterface(ip)
	if err != nil {
		return nil, err
	}

	for _, file := range ip.Parser.pkg.files {
		ip.Parser.parseFile(ip, file)
	}

	return ip, nil
}

// IsTypeImplementsInterface 判断对应类型是否实现了接口
func (ip *InterfaceParser) IsTypeImplementsInterface() (bool, error) {
	if len(ip.interfaceType.Methods.List) <= 0 {
		return true, nil
	}

	return ip.Parser.matchFuncs(ip), nil
}

func (p Parser) matchFuncs(ip *InterfaceParser) bool {
	for s, ident := range ip.InterfaceFuncIdent { // loop interface funcs
		a, ok := ip.TypeFuncIdent[s]
		if !ok {
			return false
		}
		p.Printf("found func: %v on type: %v of interface: %v", s, a, ident)
	}
	return true
}

// ParseInterface 解析接口
func (p Parser) ParseInterface(ip *InterfaceParser) error {
	tokens, err := ParseInterfaceFunc(ip.interfaceType)
	if err != nil {
		return err
	}
	for _, funcToken := range tokens {
		ip.InterfaceFuncIdent[funcToken.BuildSignature()] = funcToken.Ident
	}
	return nil
}

func (p Parser) parseFile(ip *InterfaceParser, file *File) {
	if file.file != nil {
		ast.Inspect(file.file, p.walkFunc(ip))
	}
}

func (p Parser) walkFunc(ip *InterfaceParser) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		switch node.(type) {
		case *ast.FuncDecl:
			fd := node.(*ast.FuncDecl)
			return p.funcDecl(ip, fd)
		}
		return true
	}
}

func (p Parser) funcDecl(ip *InterfaceParser, decl *ast.FuncDecl) bool {
	recv := decl.Recv
	if recv == nil || len(recv.List) == 0 {
		return true
	}
	tp, ok := recv.List[0].Type.(*ast.Ident)
	if !ok {
		return true
	}

	if ip.tp.Name != nil && tp.Name != ip.tp.Name.Name {
		return true
	}
	ft := decl.Type
	code, err := BuildFuncCode(ft)
	if err != nil {
		return true
	}
	code.FuncName = decl.Name.Name
	code.Ident = decl.Name
	ip.FuncCodes = append(ip.FuncCodes, code)
	ip.TypeFuncIdent[code.BuildSignature()] = code.Ident
	return true
}

// IsFuncImplementsInterface 检测函数是否实现了接口
func IsFuncImplementsInterface(decl ast.FuncDecl, receiveType *ast.Ident) bool {
	recv := decl.Recv
	if recv == nil || len(recv.List) <= 0 || len(recv.List[0].Names) <= 0 {
		return false
	}
	recvType := recv.List[0].Names[0]
	return recvType == receiveType
}
