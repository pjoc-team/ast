package scan

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/pjoc-team/ast/astutil"
	"github.com/pjoc-team/ast/jsonutil"
)

// TestScanPkg test
func TestScanPkg(t *testing.T) {
	packages := astutil.ParsePackage(
		[]string{
			"strconv",
		}, nil,
	)
	if len(packages) == 0 {
		t.Fatal("no package")
	}
	for _, p := range packages {
		pkg, err := ScanPkg(p, WithOnlyExported(true))
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(pkg.Errors) > 0 {
			for _, err2 := range pkg.Errors {
				os.Stderr.WriteString(err2.Error() + "\n")
			}
		}
		prettyJson, err := jsonutil.PrettyJson(pkg)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(prettyJson)
	}
}

// TestScanPkgComponent test
func TestScanPkgComponent(t *testing.T) {
	packages := astutil.ParsePackage(
		[]string{
			"pattern=../astutil/...",
		}, nil,
	)
	if len(packages) == 0 {
		t.Fatal("no package")
	}
	for _, p := range packages {
		pkg, err := ScanPkg(p, WithOnlyExported(true))
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(pkg.Errors) > 0 {
			for _, err2 := range pkg.Errors {
				os.Stderr.WriteString(err2.Error() + "\n")
			}
		}
		prettyJson, err := jsonutil.PrettyJson(pkg)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(prettyJson)
	}
}

// TestPath test
func TestPath(t *testing.T) {
	packages := astutil.ParsePackage([]string{"."}, nil)
	if len(packages) == 0 {
		t.Fatal("no package")
	}
	for _, p := range packages {
		pkg, err := ScanPkg(p, WithOnlyExported(true))
		if err != nil {
			t.Fatal(err.Error())
		}
		for s, i := range pkg.PathAndTypes {
			prettyJson, err := jsonutil.PrettyJson(i)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("%v = %s\n", s, prettyJson)
		}
	}
}

// ExampleScanPkg example for scan pkg
func ExampleScanPkg() {
	packages := astutil.ParsePackage([]string{".", "../component", "fmt", "ast"}, nil)
	if len(packages) == 0 {
		return
	}
	for _, p := range packages {
		pkg, err := ScanPkg(p, WithOnlyExported(true))
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(pkg.Errors) > 0 {
			for _, err2 := range pkg.Errors {
				os.Stderr.WriteString(err2.Error() + "\n")
			}
		}
		prettyJson, err := jsonutil.PrettyJson(pkg)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(prettyJson)
	}
}

func TestImport_AliasName(t *testing.T) {
	type fields struct {
		Path  Path
		Name  string
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "t1",
			fields: fields{
				Path:  nil,
				Value: "string",
			},
			want: "string",
		},
		{
			name: "t2",
			fields: fields{
				Path:  nil,
				Name:  "string",
				Value: "string",
			},
			want: "string",
		},
		{
			name: "t3",
			fields: fields{
				Path:  nil,
				Name:  "s",
				Value: "string",
			},
			want: "s",
		},
		{
			name: "t4",
			fields: fields{
				Path:  nil,
				Value: "hello.world.cc",
			},
			want: "cc",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				i := &Import{
					Path:  tt.fields.Path,
					Name:  tt.fields.Name,
					Value: tt.fields.Value,
				}
				if got := i.AliasName(); got != tt.want {
					t.Errorf("AliasName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestScanPkgTestData(t *testing.T) {
	packages := astutil.ParsePackage(
		[]string{
			"pattern=./testdata",
		}, nil,
	)
	if len(packages) == 0 {
		t.Fatal("no package")
	}
	for _, p := range packages {
		pkg, err := ScanPkg(p, WithOnlyExported(true))
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(pkg.Errors) > 0 {
			for _, err2 := range pkg.Errors {
				os.Stderr.WriteString(err2.Error() + "\n")
			}
		}
		prettyJson, err := jsonutil.PrettyJson(pkg)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(prettyJson)
	}
}
