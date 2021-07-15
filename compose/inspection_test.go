package compose

import (
	"reflect"
	"testing"

	"github.com/pjoc-team/ast/scan"
)

// TestBuilder_recurseFindObject test
func TestBuilder_recurseFindObject(t *testing.T) {
	type fields struct {
		Codes        *Codes
		Action       *Action
		vars         map[string]*Object
		predefines   map[string]*Object
		instances    map[string]*Object
		previousStep *Step
		currentStep  *Step
	}
	type args struct {
		name  string
		depth int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Object
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action: nil,
				vars: map[string]*Object{
					"a": {
						Name: "a",
						Type: "Person",
					},
					"b": {
						Name: "b",
						Type: "Friend",
					},
				},
				predefines:   nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				name:  "a.b.c",
				depth: 0,
			},
			want: []*Object{
				{
					Name: "a",
					Type: "Person",
				},
				{
					Name: "b",
					Type: "Friend",
				},
			},
			wantErr: false,
		},
		{
			name: "t2",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action: nil,
				vars: map[string]*Object{
					"a": {
						Name: "a",
						Type: "Person",
					},
					"b": {
						Name: "b",
						Type: "Friend",
					},
				},
				predefines:   nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				name:  "a.b",
				depth: 0,
			},
			want: []*Object{
				{
					Name: "a",
					Type: "Person",
				},
			},
			wantErr: false,
		},
		{
			name: "t3",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action: nil,
				vars: map[string]*Object{
					"a": {
						Name: "a",
						Type: "Person",
					},
					"b": {
						Name: "b",
						Type: "Friend",
					},
				},
				predefines:   nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				name:  "a",
				depth: 0,
			},
			want: []*Object{
				{
					Name: "a",
					Type: "Person",
				},
			},
			wantErr: false,
		},
		{
			name: "t4",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action: nil,
				vars: map[string]*Object{
					"a": {
						Name: "a",
						Type: "Person",
					},
					"b": {
						Name: "b",
						Type: "Friend",
					},
				},
				predefines:   nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				name:  "c",
				depth: 0,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &ActionBuilder{
					Builder: NewBuilder(tt.fields.Codes),
					CodeContext: &CodeContext{
						Vars:      tt.fields.vars,
						Instances: tt.fields.instances,
					},
				}
				got, err := c.recurseFindObject(nil, tt.args.name, tt.args.depth)
				if (err != nil) != tt.wantErr {
					t.Errorf("findObject() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("findObject() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBuilder_identName(t *testing.T) {
	type fields struct {
		Codes           *Codes
		Action          *Action
		requiredNewType map[string]*scan.Type
		vars            map[string]*Object
		predefines      map[string]*Object
		used            map[string]*Object
		instances       map[string]*Object
		previousStep    *Step
		currentStep     *Step
	}
	type args struct {
		pkg      string
		typeName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "t1",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action:          nil,
				requiredNewType: nil,
				vars: map[string]*Object{
					"alice": {
						Name: "alice",
						Type: "Student",
						Path: nil,
					},
					"bob": {
						Name: "bob",
						Type: "Student",
						Path: nil,
					},
				},
				predefines:   nil,
				used:         nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				pkg:      "t",
				typeName: "Alice",
			},
			want: "t.Alice",
		},
		{
			name: "t2",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action:          nil,
				requiredNewType: nil,
				vars: map[string]*Object{
					"alice": {
						Name: "alice",
						Type: "Student",
						Path: nil,
					},
					"tAlice": {
						Name: "tAlice",
						Type: "Student",
						Path: nil,
					},
					"bob": {
						Name: "bob",
						Type: "Student",
						Path: nil,
					},
				},
				predefines:   nil,
				used:         nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				pkg:      "t",
				typeName: "Alice",
			},
			want: "t.Alice",
		},
		{
			name: "t3",
			fields: fields{
				Codes: &Codes{
					Predefines: make([]*Object, 0),
				},
				Action:          nil,
				requiredNewType: nil,
				vars: map[string]*Object{
					"alice": {
						Name: "alice",
						Type: "Student",
						Path: nil,
					},
					"tAlice": {
						Name: "tAlice",
						Type: "Student",
						Path: nil,
					},
					"alice0": {
						Name: "alice0",
						Type: "Student",
						Path: nil,
					},
					"bob": {
						Name: "bob",
						Type: "Student",
						Path: nil,
					},
				},
				predefines:   nil,
				used:         nil,
				instances:    nil,
				previousStep: nil,
				currentStep:  nil,
			},
			args: args{
				pkg:      "t",
				typeName: "Alice",
			},
			want: "t.Alice",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				b := &ActionBuilder{
					Builder: NewBuilder(tt.fields.Codes),
					CodeContext: &CodeContext{
						RequiredNewType: tt.fields.requiredNewType,
						Vars:            tt.fields.vars,
						Used:            tt.fields.used,
						Instances:       tt.fields.instances,
					},
				}
				if got := b.genIdentName(tt.args.pkg, tt.args.typeName); got != tt.want {
					t.Errorf("genIdentName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_finalDepth(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				"",
			},
			want: "",
		},
		{
			name: "t1",
			args: args{
				"t1",
			},
			want: "t1",
		},
		{
			name: "t1",
			args: args{
				"foo.bar",
			},
			want: "bar",
		},
		{
			name: "t1",
			args: args{
				"foo.bar.c",
			},
			want: "c",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := finalField(tt.args.arg); got != tt.want {
					t.Errorf("finalField() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_checkArgsType(t *testing.T) {
	type args struct {
		object  string
		another string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t1",
			args: args{
				object:  "foo.bar",
				another: "bar",
			},
			want: true,
		},
		{
			name: "t2",
			args: args{
				object:  "foo.bar",
				another: "test.bar",
			},
			want: false,
		},
		{
			name: "t3",
			args: args{
				object:  "foo",
				another: "bar",
			},
			want: false,
		},
		{
			name: "t4",
			args: args{
				object:  "bar.foo",
				another: "foo.bar",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := checkArgsType(tt.args.object, tt.args.another); got != tt.want {
					t.Errorf("checkArgsType() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}