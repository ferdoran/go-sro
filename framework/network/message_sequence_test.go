package network

import (
	"testing"
)

func TestMessageSequence_Next(t *testing.T) {
	type fields struct {
		seed uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{
			name:   "Given a seed of 12 When Next() is called Then it returns 241",
			fields: fields{12},
			want:   224,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMessageSequence(tt.fields.seed)
			if got := ms.Next(); got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestNewMessageSequence(t *testing.T) {
//	type args struct {
//		seed uint
//	}
//	tests := []struct {
//		name string
//		args args
//		want MessageSequence
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewMessageSequence(tt.args.seed); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewMessageSequence() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_generateValue(t *testing.T) {
//	type args struct {
//		value *uint
//	}
//	tests := []struct {
//		name string
//		args args
//		want uint
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := generateValue(tt.args.value); got != tt.want {
//				t.Errorf("generateValue() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
