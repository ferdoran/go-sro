package utils

import (
	"reflect"
	"testing"
)

func TestConcatTwoUint32ToByteArray(t *testing.T) {
	type args struct {
		secret1 uint32
		secret2 uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Given 2 numbers when ConcatTwoUintToByteArray is called Then it returns a byte array in correct order",
			args: args{
				secret1: 61221336,
				secret2: 88916344,
			},
			want: []byte{216, 41, 166, 3, 120, 193, 76, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatTwoUint32ToByteArray(tt.args.secret1, tt.args.secret2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConcatTwoUint32ToByteArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
