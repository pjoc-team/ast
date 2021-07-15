package compose

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/pjoc-team/ast/jsonutil"
	"github.com/pjoc-team/ast/scan"
)

// ValueType 字段值的取值方式
type ValueType int

const (
	// ObjectValue 从对象取值
	ObjectValue ValueType = iota

	// ConstValue 固定值
	ConstValue
)

// Builder 构建服务
type Builder struct {
	// Codes 提供对象
	Codes *Codes

	// 预定义对象，跟 Codes.Predefines 是一致的
	predefines map[string]*Object
}

// NewBuilder 新建对象
func NewBuilder(codes *Codes) *Builder {
	b := &Builder{
		Codes:      codes,
		predefines: make(map[string]*Object),
	}
	if codes != nil {
		for _, predefine := range codes.Predefines {
			exists, ok := b.predefines[predefine.Name]
			if ok {
				log.Fatalf(
					"failed to put vars: %#v because key: %v is duplicate, exists is: %#v",
					predefine,
					predefine.Name, exists,
				)
			}
			b.predefines[predefine.Name] = predefine
		}
	}
	return b
}

// ActionBuilder 根据 Action 编译代码
type ActionBuilder struct {
	Builder *Builder

	// Action 代码步骤
	// Action *Action

	// CodeContext 上下文依赖
	CodeContext *CodeContext
}

// CodeContext 代码上下文需要依赖的对象
type CodeContext struct {
	// 代码结果，即拼装CodeLines
	Code string
	// 代码行
	CodeLines []string

	// 需要新建对象的类型，key是对象名称，value是类型
	RequiredNewType map[string]*scan.Type
	// Vars 本次代码初始化的变量，用于防重
	Vars map[string]*Object
	// 在steps中用到的对象
	Used map[string]*Object
	// Instances 需要初始化的对象
	Instances map[string]*Object
	// Predefines 预定义对象
	Predefines map[string]*Object
	// Imports 需要导入的包
	Imports map[string]*scan.Pkg
}

// Codes 已提供的对象
type Codes struct {
	// Packages 提供的代码包
	Packages []*scan.Pkg

	// Predefines 预定义对象集合，可以被引用
	Predefines []*Object
}

// Object 预定义对象，可以被代码引用
type Object struct {
	// Name 对象变量名
	Name string

	// Type 类型名
	Type string

	// Doc 文档
	Doc string

	// Scope 可见域
	Scope ObjectScope

	// Path 查找路径
	Path scan.Path
}

func (o *Object) String() string {
	b := &strings.Builder{}

	b.WriteString("Name= ")
	b.WriteString(o.Name)
	b.WriteString(", ")
	b.WriteString("Type= ")
	b.WriteString(o.Type)
	b.WriteString(", ")
	b.WriteString("Doc= ")
	b.WriteString(o.Doc)
	b.WriteString(", ")
	b.WriteString("Scope= ")
	b.WriteString(strconv.Itoa(int(o.Scope)))
	b.WriteString(", ")
	b.WriteString("Path= ")
	b.WriteString(o.Path.String())
	return b.String()
}

// Step 步骤
type Step struct {
	// Results 返回类型，已经返回的变量名
	Results []*scan.Field `json:"results" yaml:"results"`

	// Operation 执行的函数
	Operation *Operation `json:"operation" yaml:"operation"`

	// Args 参数
	Args []*Param `json:"args" yaml:"args"`
}

// Param 参数
type Param struct {
	scan.Field

	ValueType ValueType `json:"value_type" yaml:"valueType"`

	Value string `json:"value" yaml:"value"`
}

// Action 动作
type Action struct {
	// Steps 执行步骤
	Steps []*Step `json:"steps" yaml:"steps"`
}

// NewActionBuilder 创建action编译器
func (b *Builder) NewActionBuilder(opts ...Option) (*ActionBuilder, error) {
	o := newOptions()
	o.apply(opts...)

	codeContext := NewCodeContext(o)

	err2 := putPredefines(b.predefines, codeContext.Predefines)
	if err2 != nil {
		return nil, err2
	}
	err2 = putPredefines(o.predefines, codeContext.Predefines)
	if err2 != nil {
		return nil, err2
	}

	ab := &ActionBuilder{
		Builder:     b,
		CodeContext: codeContext,
	}
	return ab, nil
}

// NewCodeContext 创建CodeContext对象
func NewCodeContext(o *options) *CodeContext {
	codeContext := &CodeContext{
		Predefines:      make(map[string]*Object),
		Vars:            o.vars,
		Instances:       make(map[string]*Object),
		RequiredNewType: make(map[string]*scan.Type),
		Used:            make(map[string]*Object),
		Imports:         make(map[string]*scan.Pkg),
	}
	return codeContext
}

// BuildAction 构建代码
func (b *Builder) BuildAction(action *Action, opts ...Option) (*CodeContext, error) {
	ab, err := b.NewActionBuilder(opts...)
	if err != nil {
		return nil, err
	}
	return ab.BuildAction(action)
}

// BuildAction 构建action
func (b *ActionBuilder) BuildAction(action *Action) (*CodeContext, error) {
	steps, err := b.buildSteps(action.Steps)
	if err != nil {
		return nil, err
	}
	b.CodeContext.CodeLines = steps
	rs := &strings.Builder{}
	for i, step := range steps {
		if i != 0 {
			rs.WriteString("\n")
		}
		rs.WriteString(step)
	}
	b.CodeContext.Code = rs.String()
	return b.CodeContext, nil
}

// BuildStep 构建步骤
func (b *ActionBuilder) BuildStep(step *Step) (str string, err error) {
	defer func() {
		if i := recover(); i != nil {
			err = fmt.Errorf("failed to build step: %#v error: %v", *step, i)
			debug.PrintStack()
		}
	}()
	return b.buildStep(step)
}

// HasErrorResults 该步骤是否会返回error类型参数，有则返回error参数名，否则为空
func (b *ActionBuilder) HasErrorResults(step *Step) (str string, err error) {
	if step.Operation == nil || step.Operation.Type != Func {
		return "", nil
	}

	function, err := b.findFunc(step.Operation.Func)
	if err != nil {
		log.Printf("not found func: %v", step.Operation.Func)
		return "", err
	}

	for i, result := range function.Results {
		if result.Type == "error" {
			return step.Results[i].Name, nil
		}
	}

	return "", nil
}

func putPredefines(from map[string]*Object, to map[string]*Object) error {
	for k, predefine := range from {
		object, ok := to[k]
		if ok {
			err := fmt.Errorf(
				"failed to put object: %#v because key: %v is already exists: %#v",
				predefine, k, object,
			)
			log.Print(err.Error())
			return err
		}
		to[k] = predefine
		log.Printf("put predefine: %v of key: %v", predefine, k)
	}
	return nil
}

// buildSteps 编译多个步骤
func (b *ActionBuilder) buildSteps(steps []*Step) ([]string, error) {
	rs := make([]string, 0, len(steps))
	for _, step := range steps {
		stepCode, err := b.buildStep(step)
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		rs = append(rs, stepCode)
	}
	return rs, nil
}

// buildStep 单独编译一个步骤
func (b *ActionBuilder) buildStep(step *Step) (string, error) {
	if step.Operation == nil && len(step.Args) == 0 && len(step.Results) == 0 {
		return "", nil
	}

	if step.Operation == nil {
		return b.buildAssign(step)
	}

	switch step.Operation.Type {
	case Func:
		return b.buildFuncType(step)
	case Assign:
		return b.buildAssign(step)
	case Unary:
		return b.buildUnary(step)
	case New:
		return b.buildCreate(step)
	default:
		return "", UnsupportedFuncTypeError
	}
}

func (b *ActionBuilder) buildAssign(step *Step) (string, error) {
	sb := &strings.Builder{}

	if len(step.Results) <= 0 {
		return "", errors.New("assign operation's results mustn't be null")
	} else if len(step.Args) <= 0 {
		return "", errors.New("assign operation's args mustn't be null")
	}

	sb.WriteString(b.buildResultsAssign(step))

	for i, arg := range step.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.Value)
	}

	return sb.String(), nil
}

func (b *ActionBuilder) buildUnary(step *Step) (string, error) {
	if len(step.Args) != 2 {
		return "", ArgsLengthNotEqualTwoError
	} else if step.Operation.UnarySymbol == "" {
		return "", UnarySymbolNilError
	}

	sb := &strings.Builder{}

	sb.WriteString(b.buildResultsAssign(step))

	sb.WriteString(step.Args[0].Value)
	sb.WriteString(step.Operation.UnarySymbol)
	sb.WriteString(step.Args[1].Value)

	return sb.String(), nil
}

func (b *ActionBuilder) buildCreate(step *Step) (string, error) {
	if len(step.Args) != 1 {
		return "", ArgsLengthNotEqualOneError
	} else if len(step.Results) != 1 {
		return "", ArgsLengthNotEqualOneError
	}
	sb := &strings.Builder{}
	sb.WriteString(b.buildResultsAssign(step))
	sb.WriteString("&")
	sb.WriteString(step.Args[0].Value)
	sb.WriteString("{}")
	return sb.String(), nil
}

func (b *ActionBuilder) buildFuncType(step *Step) (string, error) {
	if step.Operation.Func == nil {
		return "", fmt.Errorf("func is nil")
	}

	function, err := b.findFunc(step.Operation.Func)
	if err != nil {
		log.Printf("failed to find func: %v error: %v", step.Operation.Func, err.Error())
		return "", err
	}

	sb := strings.Builder{}

	// 返回结果
	result := b.buildResult(step, function)
	sb.WriteString(result)

	// 对象和函数名
	// function := step.Operation

	functionStatement, err := b.buildFunc(function)
	if err != nil {
		log.Printf("failed to build function: %v error: %v", function, err.Error())
		return "", err
	}
	sb.WriteString(functionStatement)

	// 参数
	sb.WriteString("(")
	args, err := b.buildArgs(function, step.Args)
	if err != nil {
		log.Printf("failed to build args: %v error: %v", step.Args, err.Error())
		return "", err
	}
	sb.WriteString(args)
	sb.WriteString(")")
	return sb.String(), nil
}

func (b *ActionBuilder) buildResult(step *Step, function *scan.Func) string {
	sb := &strings.Builder{}

	// = 号左边的代码
	allExists := true
	for i, result := range step.Results {
		resultFieldDeclare := function.Results[i]
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(result.Name)

		// 判断对象是否预定义存在
		obj, err := b.findObject(result.Name)
		if err != nil && allExists {
			allExists = false
		} else if len(obj) > 0 {
			b.pubNewObjects(result.Name, obj[0])
		}
		// 如果是最后一个返回参数，则拼接=号
		if i == len(step.Results)-1 {
			if allExists {
				sb.WriteString(" = ")
			} else {
				sb.WriteString(" := ")
			}
		}
		b.putVar(result, resultFieldDeclare)
	}
	return sb.String()
}

func (b *ActionBuilder) buildArgs(function *scan.Func, args []*Param) (
	string,
	error,
) {
	sb := &strings.Builder{}
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		argName, err := b.buildArg(function, arg, i)
		if err != nil {
			json, _ := jsonutil.PrettyJson(function)
			log.Printf(
				"failed to build func: %v name: %v arg: %v error: %v", json,
				function.Name, arg.Value, err.Error(),
			)
			return "", err
		}
		sb.WriteString(argName)
	}
	return sb.String(), nil
}

func (b *ActionBuilder) buildFunc(function *scan.Func) (string, error) {
	sb := &strings.Builder{}
	pkg, err := b.findPkg(function)
	if err != nil {
		log.Printf("failed to find pkg of func: %v error: %v", function, err.Error())
		return "", err
	}
	if pkg != nil {
		b.CodeContext.Imports[pkg.Name] = pkg
		sb.WriteString(pkg.Name)
		sb.WriteString(".")
	}
	if function.Receiver != nil {
		t := function.Receiver.Type
		name := b.genIdentName("", t) // TODO found pkg
		sb.WriteString(name)
		sb.WriteString(".")
		// err := b.requiredObject(name, function.Receiver)
		// if err != nil{
		// 	return "", err
		// }
		// 实例化对象
		b.pubUsedObjects(
			name, &Object{
				Name: b.genIdentName("", t), // TODO found pkg
				Type: t,
				Path: function.Receiver.Path,
			},
		)
	}
	funcName := function.Name
	sb.WriteString(funcName)
	return sb.String(), nil
}

func (b *ActionBuilder) buildArg(function *scan.Func, arg *Param, i int) (string, error) {
	valueType := arg.ValueType
	if valueType == ObjectValue {
		valueType = detectArgValueType(arg)
	}

	switch valueType {
	case ConstValue:
		return arg.Value, nil
	case ObjectValue:
		obj, err := b.checkArg(function, arg, i)
		if err != nil {
			log.Printf("failed to check arg: %v error: %v", arg.Value, err.Error())
			// return "", err // TODO 支持组合代码时校验错误
		}
		return b.convertArg(function, arg, i, obj)
	default:
		err := fmt.Errorf("unknown arg value type: %v", arg.ValueType)
		log.Print(err.Error())
		return "", err
	}
}

// convertArg 判断参数是否需要转义
func (b *ActionBuilder) convertArg(
	function *scan.Func, arg *Param, i int, obj *Object,
) (string, error) {
	if obj == nil {
		return arg.Value, nil
	}
	var functionArgType *scan.Field
	lastArgType := function.Params[len(function.Params)-1]
	if i >= len(function.Params) {
		functionArgType = lastArgType
		// 可变变量
		if !strings.HasPrefix(functionArgType.Type, "...") {
			return "", fmt.Errorf("params size is too long")
		}
	} else {
		functionArgType = function.Params[i]
	}

	argType := functionArgType.Type
	// 如果是最后一个变量，则需要判断下是否是可变变量
	if i >= len(function.Params)-1 {
		functionArgType = lastArgType
		if strings.HasPrefix(functionArgType.Type, "...") {
			argType = strings.TrimPrefix(argType, "...")
		}
	}

	if obj.Type == "interface{}" {
		return fmt.Sprintf("%s.(%s)", arg.Value, argType), nil
	}

	return arg.Value, nil
}

func (b *ActionBuilder) checkArg(function *scan.Func, arg *Param, i int) (*Object, error) {
	object, err := b.findObject(arg.Value)
	// if len(object) > 0 { // TODO 支持待生成代码的变量类型识别，比如channelRequest.Amount
	//	return object[len(object)-1], nil
	// }
	if err != nil {
		log.Printf("not found obj: %v", arg.Value)
		return nil, err
	} else if len(object) > 0 {
		o := object[len(object)-1]
		if o == nil {
			return nil, nil
		}

		var functionArgType *scan.Field
		lastArgType := function.Params[len(function.Params)-1]
		if i >= len(function.Params) {
			functionArgType = lastArgType
			// 可变变量
			if !strings.HasPrefix(functionArgType.Type, "...") {
				return nil, fmt.Errorf("params size is too long")
			}
		} else {
			functionArgType = function.Params[i]
		}

		argType := functionArgType.Type
		// 如果是最后一个变量，则需要判断下是否是可变变量
		if i >= len(function.Params)-1 {
			functionArgType = lastArgType
			if strings.HasPrefix(functionArgType.Type, "...") {
				argType = strings.TrimPrefix(argType, "...")
			}
		}

		// 如果是interface或者是预定义类型，则跳过。
		// TODO 支持待生成代码的变量类型识别，比如channelRequest.Amount
		if argType == "interface{}" {
			return o, nil
		} else if o.Type == "interface{}" {
			return o, nil
		} else if argType == "" {
			return o, nil
		}

		// 比较arg的类型和func的参数类型
		if !checkArgsType(o.Type, argType) {
			err = fmt.Errorf(
				"find a object: %v's type is: %v but the arg: %v's type is required"+
					": %v",
				o.Name, o.Type, arg.Value, argType,
			)
			log.Printf("error: %v", err.Error())
			return nil, err
		}
		used := object[0]
		b.pubUsedObjects(used.Name, used)
		return o, nil
	}
	return nil, nil
}

// pubUsedObjects 放置需要依赖的对象
func (b *ActionBuilder) requiredObject(name string, field *scan.Field) error {
	for _, pkg := range b.Builder.Codes.Packages {
		t, ok := pkg.FindPath(field.Path)
		if !ok {
			continue
		}
		tp, ok := t.(*scan.Type)
		if !ok {
			err := fmt.Errorf(
				"path: %v type is not *scan.Type, real type: %v",
				field.Path.String(), reflect.TypeOf(tp),
			)
			log.Print(err.Error())
			return err
		}
		b.CodeContext.RequiredNewType[name] = tp
		return nil
	}
	err := fmt.Errorf("not found type of path: %v", field.Path.String())
	log.Print(err.Error())
	return err
}

// pubUsedObjects 放置需要依赖的对象
func (b *ActionBuilder) pubUsedObjects(name string, objects ...*Object) {
	for _, object := range objects {
		b.CodeContext.Used[name] = object
	}
}

// pubNewObjects 放置新建的对象
func (b *ActionBuilder) pubNewObjects(name string, objects ...*Object) {
	for _, object := range objects {
		log.Printf("put object: %v type: %v", name, object.Type)
		b.CodeContext.Vars[name] = object
	}
}

func (b *ActionBuilder) findPath(path scan.Path) (interface{}, error) {
	if path == nil {
		return nil, nil
	}
	for _, pkg := range b.Builder.Codes.Packages {
		found, b2 := pkg.FindPath(path)
		if b2 {
			return found, nil
		}
	}

	return nil, nil
}

// buildResultsAssign 编译=号及左侧的关联变量
func (b *ActionBuilder) buildResultsAssign(step *Step) string {
	sb := strings.Builder{}
	for i, result := range step.Results {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(result.Name)
	}

	allExists := b.isAllExists(step)

	if allExists {
		sb.WriteString(" = ")
	} else {
		sb.WriteString(" := ")
	}

	return sb.String()
}

func (b *ActionBuilder) isAllExists(step *Step) bool {
	allExists := true
	for i, result := range step.Results {
		// 判断对象是否预定义存在
		obj, err := b.findObject(result.Name)
		if err != nil && allExists {
			allExists = false
			b.putVar(result, &step.Args[i].Field)
		} else if len(obj) > 0 {
			b.pubUsedObjects(result.Name, obj[0])
		}
	}
	return allExists
}
