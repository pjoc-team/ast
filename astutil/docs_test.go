package astutil

import (
	"go/ast"
	"testing"
)

func TestParseComment(t *testing.T) {
	type args struct {
		comment *ast.CommentGroup
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				comment: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "test",
						},
						{
							Text: "hello",
						},
					},
				},
			},
			want: "test hello",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := ParseComment(tt.args.comment); got != tt.want {
					t.Errorf("ParseComment() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
