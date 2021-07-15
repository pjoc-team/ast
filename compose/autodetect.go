package compose

import (
	"regexp"
	"strings"
)

var numberPattern *regexp.Regexp

func init() {
	var err error
	numberPattern, err = regexp.Compile(`^[\d]+$`)
	if err != nil {
		panic(err.Error())
	}
}

// detectArgValueType 探测变量的类型
func detectArgValueType(arg *Param) ValueType {
	value := arg.Value
	if strings.Contains(value, "\"") {
		return ConstValue
	} else if value == "true" || value == "false" {
		return ConstValue
	} else if numberPattern.MatchString(value) {
		return ConstValue
	}
	return ObjectValue
}
