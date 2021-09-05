package validator

import "fmt"

// Require - func for validate input requested values
func Require(values ...string) error {
	for _, v := range values {
		if v == "" {
			return fmt.Errorf("%s is empty", v)
		}
	}
	return nil
}
