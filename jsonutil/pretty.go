package jsonutil

import (
	"bytes"
	"encoding/json"
)

// PrettyJson 打印json
func PrettyJson(i interface{}) (string, error) {
	bs := &bytes.Buffer{}
	encoder := json.NewEncoder(bs)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	if err != nil {
		return "", err
	}
	return bs.String(), nil
}
