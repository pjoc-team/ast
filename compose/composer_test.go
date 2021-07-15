package compose

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/pjoc-team/ast/jsonutil"
	"github.com/pjoc-team/ast/scan"
)

var builder *Builder

func init() {
	// astutil.ParsePackage()
	builder = NewBuilder(
		&Codes{
			Predefines: []*Object{
				{
					Name: "a1",
					Type: "string",
					Path: nil,
				},
				{
					Name: "a2",
					Type: "Person",
					Path: nil,
				},
			},
			Packages: []*scan.Pkg{
				{
					Name:   "",
					ID:     "",
					Doc:    "",
					Files:  nil,
					Errors: nil,
					PathAndTypes: map[string]interface{}{
						"compose -> composer_test.go -> Person": &scan.Func{
							Path: []string{"compose", "composer_test.go", "Person"},
							Receiver: &scan.Field{
								Path: nil,
								Name: "p",
								Type: "Person",
								Doc:  "",
							},
							Name: "t",
							Params: []*scan.Field{
								{
									Path: nil,
									Name: "name",
									Type: "string",
									Doc:  "",
								},
							},
							Results: []*scan.Field{
								{
									Path: nil,
									Name: "",
									Type: "",
									Doc:  "",
								},
							},
							Doc: "",
						},
						"compose -> composer_test.go -> t": &scan.Func{
							Path: []string{"compose", "composer_test.go", "Person"},
							Receiver: &scan.Field{
								Path: nil,
								Name: "p",
								Type: "Person",
								Doc:  "",
							},
							Name: "t",
							Params: []*scan.Field{
								{
									Path: nil,
									Name: "name",
									Type: "string",
									Doc:  "",
								},
							},
							Results: []*scan.Field{
								{
									Path: nil,
									Name: "",
									Type: "string",
									Doc:  "",
								},
							},
							Doc: "",
						},
					},
				},
			},
		},
	)
}

// TestBuilder_buildStep test
func TestBuilder_buildStep(t *testing.T) {
	b, err2 := builder.NewActionBuilder(
		WithVars(
			[]*Object{
				{
					Name: "a1",
					Type: "string",
					Path: nil,
				},
				{
					Name: "b1",
					Type: "Person",
					Path: nil,
				},
			},
		),
	)

	if err2 != nil {
		t.Fatal(err2.Error())
	}

	// b := &ActionBuilder{
	//	Builder: NewBuilder(
	//		&Codes{
	//			Predefines: make([]*Object, 0),
	//		},
	//	),
	//	CodeContext: &CodeContext{
	//		Used: make(map[string]*Object),
	//		Values: map[string]*Object{
	//			"a1": {
	//				Name: "a1",
	//				Type: "string",
	//				Path: nil,
	//			},
	//			"b1": {
	//				Name: "b1",
	//				Type: "Person",
	//				Path: nil,
	//			},
	//		},
	//	},
	// }
	st := &Step{
		Results: []*scan.Field{
			{
				Name: "rs",
				Type: "string",
				Doc:  "",
			},
		},
		Operation: &Operation{
			Func: &scan.Func{
				Path: []string{"compose", "composer_test.go", "Person"},
				Receiver: &scan.Field{
					Path: nil,
					Name: "p",
					Type: "Person",
					Doc:  "",
				},
				Name: "t",
				Params: []*scan.Field{
					{
						Path: nil,
						Name: "name",
						Type: "string",
						Doc:  "",
					},
				},
				Results: []*scan.Field{
					{
						Path: nil,
						Name: "",
						Type: "",
						Doc:  "",
					},
				},
				Doc: "",
			},
		},
		Args: []*Param{
			{
				Value: "a1",
				Field: scan.Field{
					Name: "a1",
					Type: "",
					Doc:  "",
				},
			},
		},
	}
	step, err := b.buildStep(st)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(step)
	for s, object := range b.CodeContext.Vars {
		fmt.Printf("Values %v %v\n", s, object.Type)
	}
	for s, object := range b.CodeContext.Used {
		fmt.Printf("Used %v %v\n", s, object.Type)
	}
	for s, v := range b.CodeContext.RequiredNewType {
		fmt.Printf("required: %v %v\n", s, v)
	}
}

func TestBuilder_Build(t *testing.T) {
	st := &Step{
		Results: []*scan.Field{
			{
				Name: "rs",
				Type: "string",
			},
		},
		Operation: &Operation{
			Func: &scan.Func{
				Receiver: &scan.Field{
					Path: nil,
					Name: "p",
					Type: "Person",
				},
				Name: "t",
				Params: []*scan.Field{
					{
						Path: nil,
						Name: "name",
						Type: "string",
					},
				},
				Results: []*scan.Field{
					{
						Path: nil,
						Name: "",
						Type: "",
					},
				},
				Doc: "",
			},
		},
		Args: []*Param{
			{
				Value: "a1",
				Field: scan.Field{
					Name: "a1",
					Type: "",
				},
			},
		},
	}
	st2 := &Step{
		Results: []*scan.Field{
			{
				Name: "rs2",
				Type: "string",
			},
		},
		Operation: &Operation{
			Func: &scan.Func{
				Receiver: &scan.Field{
					Path: nil,
					Name: "p",
					Type: "Person",
				},
				Name: "t",
				Params: []*scan.Field{
					{
						Path: nil,
						Name: "name",
						Type: "string",
					},
				},
				Results: []*scan.Field{
					{
						Path: nil,
						Name: "",
						Type: "",
					},
				},
				Doc: "",
			},
		},
		Args: []*Param{
			{
				Value: "rs",
				Field: scan.Field{
					Name: "rs",
					Type: "",
				},
			},
		},
	}

	// builder := NewBuilder(
	//	&Codes{
	//		Packages: nil,
	//		Predefines: []*Object{
	//			{
	//				Name: "a1",
	//				Type: "Person",
	//				Path: nil,
	//			},
	//		},
	//	},
	// )

	action := &Action{
		Steps: []*Step{st, st2},
	}

	_, err := builder.BuildAction(action)
	if err == nil {
		t.FailNow()
	} else {
		log.Printf("err: %v", err.Error())
	}
	// fmt.Println(step)
	// fmt.Println(step.Code)
	// for s, object := range step.Values {
	// 	fmt.Printf("Values %v %v\n", s, object.Type)
	// }
	// for s, object := range step.Used {
	// 	fmt.Printf("Used %v %v\n", s, object.Type)
	// }
}

func TestBuilder_BuildAction(t *testing.T) {
	st1 := &Step{
		Results: []*scan.Field{
			{
				Name: "rs",
				Type: "",
			},
		},
		Operation: &Operation{
			Type: Func,
			Func: &scan.Func{
				Receiver: &scan.Field{
					Path: nil,
					Name: "p",
					Type: "Person",
				},
				Path: []string{"compose", "composer_test.go", "Person"},
				Name: "t",
				Params: []*scan.Field{
					{
						Path: nil,
						Name: "name",
						Type: "string",
					},
				},
				Results: []*scan.Field{
					{
						Path: nil,
						Name: "",
						Type: "",
					},
				},
				Doc: "",
			},
		},
		Args: []*Param{
			{
				Value: "a1",
				Field: scan.Field{
					Name: "a1",
					Type: "",
				},
			},
		},
	}
	st2 := &Step{
		Results: []*scan.Field{
			{
				Name: "rs2",
				Type: "string",
			},
		},
		Operation: &Operation{
			Func: &scan.Func{
				Receiver: &scan.Field{
					Path: nil,
					Name: "p",
					Type: "Person",
				},
				Path: []string{"compose", "composer_test.go", "t"},
				Name: "t",
				Params: []*scan.Field{
					{
						Path: nil,
						Name: "name",
						Type: "string",
					},
				},
				Results: []*scan.Field{
					{
						Path: nil,
						Name: "",
						Type: "",
					},
				},
				Doc: "",
			},
		},
		Args: []*Param{
			{
				Value: "rs",
				Field: scan.Field{
					Name: "rs",
					Type: "",
				},
			},
		},
	}
	st3 := &Step{
		Results: []*scan.Field{
			{
				Name: "rs",
				Type: "string",
			},
		},
		Operation: &Operation{
			Func: &scan.Func{
				Receiver: &scan.Field{
					Path: nil,
					Name: "p",
					Type: "Person",
				},
				Path: []string{"compose", "composer_test.go", "t"},
				Name: "t",
				Params: []*scan.Field{
					{
						Path: nil,
						Name: "name",
						Type: "string",
					},
				},
				Results: []*scan.Field{
					{
						Path: nil,
						Name: "",
						Type: "",
					},
				},
				Doc: "",
			},
		},
		Args: []*Param{
			{
				Value: "rs",
				Field: scan.Field{
					Name: "rs",
					Type: "",
				},
			},
		},
	}

	type args struct {
		action *Action
	}
	tests := []struct {
		name    string
		fields  *Builder
		args    args
		want    *CodeContext
		wantErr bool
	}{
		{
			name:   "t1",
			fields: builder,
			args: args{
				action: &Action{Steps: []*Step{st1, st2}},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "t2",
			fields: builder,
			args: args{
				action: &Action{
					Steps: []*Step{
						st1,
						st2,
					},
				},
			},
			want: &CodeContext{
				Code: "rs := person.t(a1)\nrs2 := person.t(rs)",
				CodeLines: []string{
					"rs := person.t(a1)",
					"rs2 := person.t(rs)",
				},
				RequiredNewType: make(map[string]*scan.Type),
				Predefines: map[string]*Object{
					"a1": {
						Name: "a1",
						Type: "string",
					},
					"a2": {
						Name: "a2",
						Type: "Person",
					},
				},
				Vars: map[string]*Object{
					"rs": {
						Name: "rs",
						Type: "",
						Doc:  "",
						Path: nil,
					},
					"rs2": {
						Name: "rs2",
						Type: "string",
						Doc:  "",
						Path: nil,
					},
				},
				Used: map[string]*Object{
					"a1": {
						Name: "a1",
						Type: "string",
						Doc:  "",
						Path: nil,
					},
					"person": {
						Name: "person",
						Type: "Person",
						Doc:  "",
						Path: nil,
					},
				},
				Instances: make(map[string]*Object),
				Imports: map[string]*scan.Pkg{
				},
			},
			wantErr: false,
		},
		{
			name: "t3",
			fields: NewBuilder(
				&Codes{
					Packages: builder.Codes.Packages,
					Predefines: []*Object{
						{
							Name: "a1",
							Type: "string",
							Path: nil,
						},
						{
							Name: "a2",
							Type: "Person",
							Path: nil,
						},
					},
				},
			),
			args: args{
				action: &Action{
					Steps: []*Step{
						st1,
						st3,
					},
				},
			},
			want: &CodeContext{
				Code: "rs := person.t(a1)\nrs = person.t(rs)",
				CodeLines: []string{
					"rs := person.t(a1)",
					"rs = person.t(rs)",
				},
				RequiredNewType: make(map[string]*scan.Type),
				Predefines: map[string]*Object{
					"a1": {
						Name: "a1",
						Type: "string",
					},
					"a2": {
						Name: "a2",
						Type: "Person",
					},
				},
				Vars: map[string]*Object{
					"rs": {
						Name: "rs",
						Type: "string",
						Doc:  "",
						Path: nil,
					},
				},
				Used: map[string]*Object{
					"a1": {
						Name: "a1",
						Type: "string",
						Doc:  "",
						Path: nil,
					},
					"person": {
						Name: "person",
						Type: "Person",
						Doc:  "",
						Path: nil,
					},
					"rs": {
						Name: "rs",
						Type: "string",
						Doc:  "",
						Path: nil,
					},
				},
				Instances: make(map[string]*Object),
				Imports: map[string]*scan.Pkg{
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				b := &Builder{
					Codes:      tt.fields.Codes,
					predefines: tt.fields.predefines,
				}
				got, err := b.BuildAction(tt.args.action)
				if (err != nil) != tt.wantErr {
					t.Errorf("BuildAction() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					wantJson, _ := jsonutil.PrettyJson(tt.want)
					gotJson, _ := jsonutil.PrettyJson(got)
					t.Errorf("BuildAction() got = %v, want %v", gotJson, wantJson)
				}
			},
		)
	}
}

func TestActionBuilder_buildAssign(t *testing.T) {
	type fields struct {
		Builder     *Builder
		CodeContext *CodeContext
	}
	type args struct {
		step *Step
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Builder:     nil,
				CodeContext: nil,
			},
			args: args{
				step: &Step{
					Results: []*scan.Field{
						{
							Name: "A.b",
							Type: "",
							Doc:  "",
						},
					},
					Operation: nil,
					Args: []*Param{
						{
							Field:     scan.Field{},
							ValueType: 1,
							Value:     "\"b\"",
						},
					},
				},
			},
			want:    "A.b = \"b\"",
			wantErr: false,
		},
		{
			name: "t2",
			fields: fields{
				Builder:     nil,
				CodeContext: nil,
			},
			args: args{
				step: &Step{
					Results: []*scan.Field{
						{
							Name: "A.b",
						},
						{
							Name: "c.d",
						},
					},
					Operation: &Operation{
						Type: Assign,
					},
					Args: []*Param{
						{
							Value: "b.c",
						},
						{
							Value: "b.c",
						},
					},
				},
			},
			want:    "A.b, c.d = b.c, b.c",
			wantErr: false,
		},
		{
			name: "t3",
			fields: fields{
				Builder:     nil,
				CodeContext: nil,
			},
			args: args{
				step: &Step{
					Results: []*scan.Field{
						{
							Name: "A.b",
						},
					},
					Operation: &Operation{
						Type: Assign,
					},
					Args: []*Param{
						{
							Value: "b.c",
						},
					},
				},
			},
			want:    "A.b = b.c",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				b := &ActionBuilder{
					Builder:     tt.fields.Builder,
					CodeContext: tt.fields.CodeContext,
				}
				got, err := b.buildAssign(tt.args.step)
				if (err != nil) != tt.wantErr {
					t.Errorf("buildAssign() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("buildAssign() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestActionBuilder_buildUnary(t *testing.T) {
	type fields struct {
		Builder     *Builder
		CodeContext *CodeContext
	}
	type args struct {
		step *Step
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
			},
			args: args{
				step: &Step{
					Results: []*scan.Field{
						{
							Path: nil,
							Name: "ok",
							Type: "",
							Doc:  "",
						},
					},
					Operation: &Operation{
						Type:        Unary,
						Func:        nil,
						UnarySymbol: "==",
					},
					Args: []*Param{
						{
							Field:     scan.Field{},
							ValueType: 0,
							Value:     "response.Code",
						},
						{
							Field:     scan.Field{},
							ValueType: 1,
							Value:     "100",
						},
					},
				},
			},
			want:    "ok = response.Code==100",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				b := &ActionBuilder{
					Builder:     tt.fields.Builder,
					CodeContext: tt.fields.CodeContext,
				}
				got, err := b.buildUnary(tt.args.step)
				if (err != nil) != tt.wantErr {
					t.Errorf("buildUnary() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("buildUnary() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestActionBuilder_buildCreate(t *testing.T) {
	type fields struct {
		Builder     *Builder
		CodeContext *CodeContext
	}
	type args struct {
		step *Step
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				Builder:     &Builder{},
				CodeContext: NewCodeContext(newOptions()),
			},
			args: args{
				step: &Step{
					Results: []*scan.Field{
						{
							Path: nil,
							Name: "a",
							Type: "",
							Doc:  "",
						},
					},
					Operation: &Operation{
						Type:        New,
						Func:        nil,
						UnarySymbol: "",
					},
					Args: []*Param{
						{
							Value: "strings.Builder",
						},
					},
				},
			},
			want:    "a := &strings.Builder{}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				b := &ActionBuilder{
					Builder:     tt.fields.Builder,
					CodeContext: tt.fields.CodeContext,
				}
				got, err := b.buildCreate(tt.args.step)
				if (err != nil) != tt.wantErr {
					t.Errorf("buildCreate() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("buildCreate() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
