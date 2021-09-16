package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	ErrEmptyInputValue = errors.New("some of input values is empty")
	ErrBadValue        = errors.New("can't parse int values")
)

func TestRequire(t *testing.T) {
	type args struct {
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Positive test",
			args: args{
				values: []interface{}{1, 2, 3, 4, 5, 6, 7},
			},
			wantErr: nil,
		},
		{
			name: "Empty ID",
			args: args{
				values: []interface{}{nil, 1, 2, 3, 4},
			},
			wantErr: ErrEmptyInputValue,
		},
		//{
		//	name: "Empty Value",
		//	args: args{
		//		values: []interface{}{1, 2, 3, nil, 5, 6},
		//	},
		//	wantErr: ErrEmptyInputValue,
		//},
		{
			name: "Bad Value",
			args: args{
				values: []interface{}{1, 2, 3, "4", 5, 6},
			},
			wantErr: ErrBadValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, Require(tt.args.values...))
		})
	}
}
