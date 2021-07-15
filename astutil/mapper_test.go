package astutil

import (
	"fmt"
	"go/ast"
	"testing"

	"golang.org/x/tools/go/packages"
)

func a() (*string, error) {
	return nil, nil
}

func TestFieldType(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, "./mapper_test.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			decls := file.Decls
			for _, decl := range decls {
				switch tp := decl.(type) {
				case *ast.FuncDecl:
					if tp.Name.Name != "a" {
						continue
					}
					funcType := tp.Type
					for _, field := range funcType.Results.List {
						fieldType, err := FieldType(field.Type)
						if err != nil {
							t.Fatal(err.Error())
						}
						fmt.Println(fieldType)
					}
					field := funcType.Results.List[0]
					fieldType, err := FieldType(field.Type)
					if err != nil || fieldType != "*string" {
						t.FailNow()
					}
				}
			}
		}
	}
}

func TestValueType(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, "./testdata/testdata.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			decls := file.Decls
			for _, decl := range decls {
				switch tp := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range tp.Specs {
						switch st := spec.(type) {
						case *ast.ValueSpec:
							fmt.Printf("name: %v\n", st.Names[0].Name)
							for _, value := range st.Values {
								valueType, err := ParseValue(value)
								if err != nil {
									t.Fatal(err.Error())
								}
								fmt.Println(valueType)
							}
						}
					}
				}
			}
		}
	}
}

func TestPbValueType(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, "github.com/pjoc-team/ast")
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			decls := file.Decls
			for _, decl := range decls {
				switch tp := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range tp.Specs {
						switch st := spec.(type) {
						case *ast.ValueSpec:
							fmt.Printf("name: %v\n", st.Names[0].Name)
							for _, value := range st.Values {
								valueType, err := ParseValue(value)
								if err != nil {
									t.Fatal(err.Error())
								}
								fmt.Println(valueType)
							}
						}
					}
				}
			}
		}
	}
}
