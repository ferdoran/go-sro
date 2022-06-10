package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestParseSilkroadTime(t *testing.T) {
	type args struct {
		timestamp uint32
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Parses Silkroad timestamp as time.Time correctly",
			args: args{timestamp: 0x1D6BA693},
			want: time.Date(2019, 10, 9, 23, 22, 7, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSilkroadTime(tt.args.timestamp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSilkroadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSilkroadTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "Converts time.Time to Silkroad Timestamp correctly",
			args: args{t: time.Date(2019, 10, 9, 23, 22, 7, 0, time.Local)},
			want: 0x1D6BA693,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSilkroadTime(tt.args.t); got != tt.want {
				t.Errorf("ToSilkroadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
