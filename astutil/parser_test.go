package astutil

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"testing"
)

type ObjectType struct {

}

type Intface interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Hello(config packages.Config, t ObjectType, arr []ObjectType) error
}

type impl struct {
}

func (i impl) Hello(config packages.Config, t ObjectType, arr []ObjectType) error {
	panic("implement me")
}

func (i impl) Marshal(v interface{}) ([]byte, error) {
	panic("implement me")
}

func (i impl) Unmarshal(data []byte, v interface{}) error {
	panic("implement me")
}

func TestIsTypeImplementsInterface(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		t.Fatal(err)
	}
	//if len(pkgs) != 1 {
	//	t.Fatalf("error: %d packages found, packages: %#v", len(pkgs), pkgs)
	//}
	//pkg := pkgs[0]
	for _, pkg := range pkgs {
		fmt.Println(pkg.Name)
		p := NewParser(pkg)
		ip, err := p.ParseTypeAndInterface("impl", "Intface")

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println(ip)
		ok, err := ip.IsTypeImplementsInterface()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(ok)
		if ok {
			return
		}
	}

	t.Fail()

}
