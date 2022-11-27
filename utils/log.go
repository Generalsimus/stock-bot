package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func LogStruct(a ...any) (n int, err error) {
	var logs []interface{}
	for _, value := range a {

		if reflect.ValueOf(value).Kind() == reflect.Struct {
			empJSON, _ := json.MarshalIndent(value, "", "  ")
			logs = append(logs, string(empJSON))
		} else {
			logs = append(logs, value)
		}
	}
	return fmt.Println(logs...)
}
