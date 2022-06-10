package security

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestG_pow_X_mod_P(t *testing.T) {
	type fields struct {
		generator uint32
		private   uint32
		prime     uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "Test with precomputed random values",
			fields: fields{uint32(1891972126), uint32(706780046), uint32(449512514)},
			want:   uint32(61221336),
		},
		{
			name:   "Test with a private of 0, should return 1",
			fields: fields{uint32(1891972126), uint32(0), uint32(449512514)},
			want:   uint32(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := G_pow_X_mod_P(tt.fields.generator, tt.fields.private, tt.fields.prime)
			if result != tt.want {
				t.Errorf("Computation should resulted in = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestKeyTransformValue(t *testing.T) {
	type args struct {
		val1    uint32
		val2    uint32
		key     uint32
		keyByte byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "it works",
			args: args{
				val1:    61221336,
				val2:    88916344,
				key:     106844496,
				keyByte: byte(106844496 & 3),
			},
			want: []byte{240, 83, 162, 10, 176, 211, 230, 14},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			value := make([]byte, 8)
			s1 := make([]byte, 4)
			s2 := make([]byte, 4)
			binary.LittleEndian.PutUint32(s1, tt.args.val1)
			binary.LittleEndian.PutUint32(s2, tt.args.val2)
			value[0] = s1[0]
			value[1] = s1[1]
			value[2] = s1[2]
			value[3] = s1[3]
			value[4] = s2[0]
			value[5] = s2[1]
			value[6] = s2[2]
			value[7] = s2[3]
			if got := KeyTransformValue(value, tt.args.key, tt.args.keyByte); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyTransformValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
