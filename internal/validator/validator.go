package validator

import (
	"fmt"
	"reflect"
)

// Require func for validate input requested values
func Require(values ...interface{}) error {
	for _, v := range values {
		if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
			return fmt.Errorf("some of input values is empty")
		}
		if _, ok := v.(uint); !ok {
			return fmt.Errorf("can't parse int values")
		}
	}
	return nil
}
