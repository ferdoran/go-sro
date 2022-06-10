package network

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"testing"
)

const (
	Float32Val float32 = 123.45
	Float64Val float64 = 5544332211.1122334455
	UInt16Val  uint16  = 0x05
	UInt32Val  uint32  = 0xABCD
	UInt64Val  uint64  = 0xFFFFDDDD
	StringVal          = "I am a string"
)

func fakeByteArray() []byte {
	arr := make([]byte, 0, 128)
	for i := 0; i < cap(arr); i++ {
		arr = append(arr, byte(i))
	}
	return arr
}

func reversedFakeByteArray() []byte {
	arr := fakeByteArray()

	for i := len(arr)/2 - 1; i >= 0; i-- {
		opp := len(arr) - 1 - i
		arr[i], arr[opp] = arr[opp], arr[i]
	}
	return arr
}

func encryptedFakeByteArray() []byte {
	arr := fakeByteArray()
	arr[EncryptOffset-1] = 0x80
	return arr
}

func bufferIncludingFloat32(value float32) []byte {
	arr := fakeByteArray()
	bits := math.Float32bits(value)
	arr[DataOffset] = byte(bits)
	arr[DataOffset+1] = byte(bits >> 8)
	arr[DataOffset+2] = byte(bits >> 16)
	arr[DataOffset+3] = byte(bits >> 24)
	return arr
}

func bufferIncludingFloat64(value float64) []byte {
	arr := fakeByteArray()
	bits := math.Float64bits(value)
	arr[DataOffset] = byte(bits)
	arr[DataOffset+1] = byte(bits >> 8)
	arr[DataOffset+2] = byte(bits >> 16)
	arr[DataOffset+3] = byte(bits >> 24)
	arr[DataOffset+4] = byte(bits >> 32)
	arr[DataOffset+5] = byte(bits >> 40)
	arr[DataOffset+6] = byte(bits >> 48)
	arr[DataOffset+7] = byte(bits >> 56)
	return arr
}

func bufferIncludingStringValue(value string) []byte {
	arr := fakeByteArray()
	binary.LittleEndian.PutUint16(arr[DataOffset:DataOffset+2], uint16(len(value)))
	strBytes := []byte(value)
	for i := 0; i < len(value); i++ {
		arr[int(DataOffset)+i+2] = strBytes[i]
	}
	return arr
}

func bufferWithMaximumSize() []byte {
	arr := make([]byte, 32767)
	for i := range arr {
		arr[i] = 0xFF
	}
	return arr
}

var buffer = fakeByteArray()
var reversedBuffer = reversedFakeByteArray()
var encryptedBuffer = encryptedFakeByteArray()
var emptyBuffer = make([]byte, 0, 128)
var bufferWithFloat32 = bufferIncludingFloat32(Float32Val)
var bufferWithFloat64 = bufferIncludingFloat64(Float64Val)
var bufferWithString = bufferIncludingStringValue(StringVal)

func TestNewPacket(t *testing.T) {
	type args struct {
		buffer []byte
	}
	tests := []struct {
		name string
		args args
		want Packet
	}{
		{
			name: "For a given byte array it creates the Packet correctly including all fields",
			args: args{buffer},
			want: Packet{
				Buffer:    buffer,
				Data:      buffer[DataOffset:],
				Size:      int(binary.LittleEndian.Uint16(buffer[:2])) + int(DataOffset),
				DataSize:  int(binary.LittleEndian.Uint16(buffer[:2])),
				Encrypted: false,
				MessageID: 0x0302,
				CRC:       buffer[CRCOffset],
				Sequence:  buffer[SequenceOffset],
			},
		},
		{
			name: "For reversed fake bytes it creates the Packet correctly",
			args: args{reversedBuffer},
			want: Packet{
				Buffer:    reversedBuffer,
				Data:      reversedBuffer[DataOffset:],
				Size:      int(binary.LittleEndian.Uint16(reversedBuffer[:2])) + int(DataOffset),
				DataSize:  int(binary.LittleEndian.Uint16(reversedBuffer[:2])),
				Encrypted: false,
				MessageID: 0x7C7D,
				CRC:       reversedBuffer[CRCOffset],
				Sequence:  reversedBuffer[SequenceOffset],
			},
		},
		{
			name: "For an encrypted buffer Encrypted is true",
			args: args{encryptedBuffer},
			want: Packet{
				Buffer:    encryptedBuffer,
				Data:      encryptedBuffer[DataOffset:],
				Size:      int(encryptedBuffer[0]) + int(DataOffset),
				DataSize:  int(encryptedBuffer[0]),
				Encrypted: true,
				MessageID: 0x0302,
				CRC:       encryptedBuffer[CRCOffset],
				Sequence:  encryptedBuffer[SequenceOffset],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPacket(tt.args.buffer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadByte(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    byte
		wantErr bool
	}{
		{
			name:    "For fake bytes",
			fields:  fields{buffer},
			want:    0x06,
			wantErr: false,
		},
		{
			name:    "For reversed fake bytes",
			fields:  fields{reversedBuffer},
			want:    0x79,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    0x00,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadByte()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadByte() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadFloat32(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    float32
		wantErr bool
	}{
		{
			name:    "For a buffer including a float32 value it reads the float value correctly",
			fields:  fields{bufferWithFloat32},
			want:    Float32Val,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadFloat32()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadFloat32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadFloat64(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{
			name: "For a buffer including a float64 value it reads the float value correctly",
			fields: fields{
				Buffer: bufferWithFloat64,
			},
			want:    Float64Val,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    0x00,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadFloat64()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadFloat64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadString(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "For a buffer including a string value it reads the float value correctly",
			fields:  fields{bufferWithString},
			want:    StringVal,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadString()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadUInt16(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint16
		wantErr bool
	}{
		{
			name:    "For a fake byte array it reads an uint16 correctly",
			fields:  fields{buffer},
			want:    0x0706,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadUInt16()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUInt16() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadUInt32(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint32
		wantErr bool
	}{
		{
			name:    "For a fake byte array it reads an uint32 correctly",
			fields:  fields{buffer},
			want:    0x09080706,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadUInt32()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUInt32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ReadUInt64(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{
			name:    "For a fake byte array it reads an uint64 correctly",
			fields:  fields{buffer},
			want:    0x0D0C0B0A09080706,
			wantErr: false,
		},
		{
			name:    "For an empty buffer it returns an error",
			fields:  fields{emptyBuffer},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			got, err := p.ReadUInt64()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_WriteByte(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a byte is written into it Then it has the byte",
			fields:  fields{make([]byte, 4096)},
			args:    args{0xDE},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a byte is written into it Then it resizes the buffer",
			fields:  fields{make([]byte, 4096)},
			args:    args{0xDE},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 1
			p.WriteByte(tt.args.value)
			v, err := p.ReadByte()
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteByte() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WriteByte() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WriteByte() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_WriteFloat32(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value float32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a float32 is written into it Then it has the float32",
			fields:  fields{make([]byte, 4090)},
			args:    args{Float32Val},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a float32 is written into it Then it resizes the buffer",
			fields:  fields{make([]byte, 4096)},
			args:    args{Float32Val},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 4
			p.WriteFloat32(tt.args.value)
			v, err := p.ReadFloat32()
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFloat32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WriteFloat32() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WriteFloat32() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_WriteFloat64(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a float64 is written into it Then it has the float64",
			fields:  fields{make([]byte, 4090)},
			args:    args{Float64Val},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a float64 is written into it Then it resizes the buffer",
			fields:  fields{make([]byte, 4096)},
			args:    args{Float64Val},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 8
			p.WriteFloat64(tt.args.value)
			v, err := p.ReadFloat64()
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFloat64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WriteFloat64() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WriteFloat64() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_WriteString(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a string is written into it Then it has the string",
			fields:  fields{make([]byte, 4090)},
			args:    args{StringVal},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a string is written into it Then it resizes the buffer",
			fields:  fields{make([]byte, 4090)},
			args:    args{StringVal},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 2 + len(tt.args.value)
			p.WriteString(tt.args.value)
			v, err := p.ReadString()
			if (err != nil) != tt.wantErr {
				t.Errorf("WritrString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WritrString() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WireString() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_WriteUInt16(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a uint16 is written into it Then it has the uint16",
			fields:  fields{make([]byte, 4090)},
			args:    args{UInt16Val},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a uint16 is written into it Then it resizes the buffer",
			fields:  fields{make([]byte, 4096)},
			args:    args{UInt16Val},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 2
			p.WriteUInt16(tt.args.value)
			v, err := p.ReadUInt16()
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteUInt16() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WriteUInt16() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WriteUInt16() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_WriteUInt32(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a uint32 is written into it Then it has the uint32",
			fields:  fields{make([]byte, 4090)},
			args:    args{UInt32Val},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a uint32 is written into it Then it returns an error",
			fields:  fields{make([]byte, 4096)},
			args:    args{UInt32Val},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 4
			p.WriteUInt32(tt.args.value)
			v, err := p.ReadUInt32()
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteUInt32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WriteUInt32() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WriteUInt32() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_WriteUInt64(t *testing.T) {
	type fields struct {
		Buffer []byte
	}
	type args struct {
		value uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Given an empty packet When a uint64 is written into it Then it has the uint64",
			fields:  fields{make([]byte, 4090)},
			args:    args{UInt64Val},
			wantErr: false,
		},
		{
			name:    "Given a full packet When a uint64 is written into it Then it resizes the buffer",
			fields:  fields{make([]byte, 4090)},
			args:    args{UInt64Val},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPacket(tt.fields.Buffer)
			expectedDataSize := p.DataSize + 8
			p.WriteUInt64(tt.args.value)
			v, err := p.ReadUInt64()
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteUInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && v != tt.args.value {
				t.Errorf("WriteUInt64() got = %v, want %v", v, tt.args.value)
			}
			if (err == nil) && expectedDataSize != p.DataSize {
				t.Errorf("WriteUInt64() did not increase the data size correctly. Got = %v, want %v", p.DataSize, expectedDataSize)
			}
		})
	}
}

func TestPacket_ToBytes(t *testing.T) {
	nullBuffer := make([]byte, 128)
	maxSizeBuffer := bufferWithMaximumSize()
	expectedMaxSizeBuffer := maxSizeBuffer
	wantedHeader := []byte{0xFF, 0x8F, 0x00, 0x50, 0x00, 0x00}
	expectedMaxSizeBuffer[0] = wantedHeader[0]
	expectedMaxSizeBuffer[1] = wantedHeader[1]
	expectedMaxSizeBuffer[2] = wantedHeader[2]
	expectedMaxSizeBuffer[3] = wantedHeader[3]
	expectedMaxSizeBuffer[4] = wantedHeader[4]
	expectedMaxSizeBuffer[5] = wantedHeader[5]
	type args struct {
		MessageID    uint16
		Encrypted    bool
		Buffer       []byte
		Data         []byte
		Size         int
		DataSize     int
		CRC          byte
		Sequence     byte
		ServerPacket bool
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Given MessageID, Encrypted, CRC, Sequence and empty Buffer, When ToBytes() is called it writes the values into the buffer correctly",
			args: args{
				0x5000,
				true,
				nullBuffer,
				nullBuffer[DataOffset:],
				6,
				0,
				0xAB,
				0xFE,
				false,
			},
			want: []byte{0x00, 0x80, 0x00, 0x50, 0xFE, 0xAB},
		},
		{
			name: "Given a maximum size buffer When ToBytes() is called Then it writes the values into the buffer correctly",
			args: args{
				0x5000,
				true,
				maxSizeBuffer,
				maxSizeBuffer[DataOffset:],
				len(maxSizeBuffer),
				len(maxSizeBuffer[DataOffset:]),
				0x11,
				0x99,
				true,
			},
			want: expectedMaxSizeBuffer,
		},
		{
			name: "Given an unencrypted buffer When ToBytes() is called Then it writes the values into the buffer correctly",
			args: args{
				0x5000,
				false,
				nullBuffer,
				nullBuffer[DataOffset:],
				6,
				0,
				0x11,
				0x99,
				false,
			},
			want: []byte{0x00, 0x00, 0x00, 0x50, 0x99, 0x11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				MessageID:    tt.args.MessageID,
				Encrypted:    tt.args.Encrypted,
				Buffer:       tt.args.Buffer,
				Data:         tt.args.Data,
				Size:         tt.args.Size,
				DataSize:     tt.args.DataSize,
				Sequence:     tt.args.Sequence,
				CRC:          tt.args.CRC,
				ServerPacket: tt.args.ServerPacket,
			}
			result := p.ToBytes()

			if !bytes.Equal(result, tt.want) {
				t.Errorf("ToBytes() got = %X, want %X", result, tt.want)
			}
		})
	}
}

func TestPacket_WriteMassivePackets(t *testing.T) {
	type args struct {
		bytesToWrite int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Writing 10.000 bytes will write them all and makes the packet massive",
			args: args{
				bytesToWrite: 10_000,
			},
			want: 10_000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := EmptyPacket()

			for i := 0; i < tt.args.bytesToWrite; i++ {
				p.WriteByte(byte(i % 255))
			}
			if p.DataSize < tt.args.bytesToWrite {
				t.Errorf("Did not write all bytes. Got %v, want %v", p.DataSize, tt.args.bytesToWrite)
			}

			if !p.IsMassive() {
				t.Errorf("Packet is not massive! Wanted bytes %v", tt.args.bytesToWrite)
			}

			packets := p.SplitIntoMassivePackets()
			if expectedAmount := int(math.Ceil(float64(p.DataSize)/4096)) + 1; len(packets) != expectedAmount {
				t.Errorf("Invalid amount of packets. Got %v, want %v", len(packets), expectedAmount)
			}
		})
	}
}

// TODO Test String()
