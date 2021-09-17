package validator

import (
	"fmt"
)

// Require func for validate input requested values
func Require(values ...int) error {
	for _, v := range values {
		if v == 0 {
			return fmt.Errorf("some of input values is empty")
		}
	}
	return nil
}
