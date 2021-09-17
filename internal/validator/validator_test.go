package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ErrEmptyValue = errors.New("some of input values is empty")

func TestRequire(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Positive test",
			args: args{
				values: []int{1, 2, 3, 4, 5, 6, 7},
			},
			wantErr: nil,
		},
		{
			name: "Empty ID",
			args: args{
				values: []int{0, 1, 2, 3, 4, 5, 6},
			},
			wantErr: ErrEmptyValue,
		},
		{
			name: "Empty Value",
			args: args{
				values: []int{1, 2, 3, 0, 5, 6, 7},
			},
			wantErr: ErrEmptyValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, Require(tt.args.values...))
		})
	}
}
