package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	ErrEmptyInputValue = errors.New("some of input values is empty")
	ErrBadValue        = errors.New("can't parse int from string value")
)

func TestRequire(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Positive test",
			args: args{
				values: []string{"1", "100", "101", "102", "103", "104", "105", "106", "107", "108"},
			},
			wantErr: nil,
		},
		{
			name: "Empty ID",
			args: args{
				values: []string{"", "100", "101", "102", "103", "104", "105", "106", "107", "108"},
			},
			wantErr: ErrEmptyInputValue,
		},
		{
			name: "Empty Value",
			args: args{
				values: []string{"1", "100", "101", "102", "103", "", "105", "106", "107", "108"},
			},
			wantErr: ErrEmptyInputValue,
		},
		{
			name: "Bad Value",
			args: args{
				values: []string{"1", "100", "101", "102", "103", " 104", "105", "106", "107", "108"},
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
