package path

import (
	"testing"

	"github.com/blademainer/commons/pkg/path"
)

// Test_getMod test
func Test_getMod(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				gomodPath,
			},
			want:    "github.com/pjoc-team/ast",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := GetMod(tt.args.path)
				if (err != nil) != tt.wantErr {
					t.Errorf("getMod() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("getMod() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSourcePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				path: path.NeighborPath("."),
			},
			want: "github.com/pjoc-team/ast/path",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := SourcePath(tt.args.path); got != tt.want {
					t.Errorf("SourcePath() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
