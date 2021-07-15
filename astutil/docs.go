package astutil

import (
	"go/ast"
	"strings"
)

// ParseComment 解析文档，生成可阅读的格式
func ParseComment(comment *ast.CommentGroup) string {
	if comment == nil {
		return ""
	}
	b := strings.Builder{}
	delimiter := ""
	for _, c := range comment.List {
		b.WriteString(delimiter)
		text := strings.TrimPrefix(c.Text, "//")
		text = strings.TrimPrefix(text, " ")
		if text == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString(text)
		delimiter = " "
	}
	return b.String()
}
