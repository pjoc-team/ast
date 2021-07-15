package compose

import (
	"log"

	"github.com/pjoc-team/ast/scan"
)

func (b *ActionBuilder) putVar(result *scan.Field, resultFieldDeclare *scan.Field) {
	log.Printf("put var: %v type: %v", result.Name, resultFieldDeclare.Type)
	b.CodeContext.Vars[result.Name] = &Object{
		Name: result.Name,
		Type: resultFieldDeclare.Type,
	}
}
