package astutil

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ParsePackage analyzes the single package constructed from the patterns and tags.
// ParsePackage exits if there is an error.
func ParsePackage(patterns []string, tags []string) []*packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypesInfo,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) == 0 {
		return pkgs
	} else if len(pkgs[0].Errors) > 0 {
		log.Fatal(pkgs[0].Errors)
	}
	return pkgs
}
