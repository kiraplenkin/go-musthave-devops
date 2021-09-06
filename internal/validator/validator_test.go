package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
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
				values: []string{"1", "test_type", "test_value"},
			},
			wantErr: nil,
		},
		{
			name: "Empty ID",
			args: args{
				values: []string{"", "test_type", "test_value"},
			},
			wantErr: errors.New("some of input values is empty"),
		},
		{
			name: "Empty Type",
			args: args{
				values: []string{"1", "", "test_value"},
			},
			wantErr: errors.New("some of input values is empty"),
		},
		{
			name: "Empty Value",
			args: args{
				values: []string{"1", "test_type", ""},
			},
			wantErr: errors.New("some of input values is empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, Require(tt.args.values...))
		})
	}
}
