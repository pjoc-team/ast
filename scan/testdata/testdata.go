// Package testdata for test
//
package testdata

var (
	// StringVar string var
	StringVar = "string var"
	// StructVar struct var
	StructVar = &StructType{}
)

const (
	// StringConst const
	StringConst = "string const"
)

// StructType struct type
type StructType struct {
}

// StringType string type
type StringType string

// ChanType chan type
type ChanType chan string

// ChanReceiverType chan type
type ChanReceiverType <-chan string

// FuncType func type
type FuncType func(param1 string, variable ...int) (ret string, err error)

// FuncDeclare func type
func FuncDeclare(param1 string, variable ...int) (string, error) {
	return "", nil
}

// FuncType1 func type
func FuncType1(param1 string, variable ...int) (ret string, err error) {
	return "", nil
}

// FuncVar func var
var FuncVar = func(param1 string, variable ...int) (ret string, err error) {
	return "", nil
}
