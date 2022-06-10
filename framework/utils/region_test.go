package utils

import "testing"

func TestUint16ToXAndZ(t *testing.T) {
	type args struct {
		region int16
	}
	tests := []struct {
		name  string
		args  args
		wantX int
		wantZ int
	}{
		{
			name:  "RegionID 25000 is X 168 and Z 97",
			args:  args{25000},
			wantX: 168,
			wantZ: 97,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotZ := Int16ToXAndZ(tt.args.region)
			if gotX != tt.wantX {
				t.Errorf("Int16ToXAndZ() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotZ != tt.wantZ {
				t.Errorf("Int16ToXAndZ() gotZ = %v, want %v", gotZ, tt.wantZ)
			}
		})
	}
}

func TestXAndZToUint16(t *testing.T) {
	type args struct {
		x byte
		z byte
	}
	tests := []struct {
		name         string
		args         args
		wantRegionId int16
	}{
		{
			name:         "X 168 and Z 97 are RegionID 25000",
			args:         args{168, 97},
			wantRegionId: 25000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRegionId := XAndZToInt16(tt.args.x, tt.args.z); gotRegionId != tt.wantRegionId {
				t.Errorf("XAndZToInt16() = %v, want %v", gotRegionId, tt.wantRegionId)
			}
		})
	}
}
