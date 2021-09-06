package validator

import "errors"

// Require - func for validate input requested values
func Require(values ...string) error {
	for _, v := range values {
		if v == "" {
			return errors.New("some of input values is empty")
		}
	}
	return nil
}
