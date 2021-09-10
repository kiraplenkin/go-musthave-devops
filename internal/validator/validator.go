package validator

import (
	"fmt"
	"strconv"
)

// Require func for validate input requested values
func Require(values ...string) error {
	for _, v := range values {
		if v == "" {
			return fmt.Errorf("some of input values is empty")
		}
		_, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("can't parse int from string value")
		}
	}
	return nil
}

// Transform func for transforming values from string to unt
func Transform(value string) int {
	v, _ := strconv.Atoi(value)
	return v
}
