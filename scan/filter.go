package scan

import (
	"context"
	"go/ast"
)

// Filter 过滤器。如果通过则返回true，否则屏蔽请返回false
type Filter func(ctx context.Context, expr ast.Expr) bool
