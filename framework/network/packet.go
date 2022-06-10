package network

import (
	"encoding/binary"
	"fmt"
	"github.com/ferdoran/go-sro/framework/security/blowfish"
	"github.com/ferdoran/go-sro/framework/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math"
	"strings"
)

const (
	BufferSize uint16 = 4096
	HeaderSize uint16 = 6
	DataSize   uint16 = BufferSize - HeaderSize

	HeaderOffset   uint16 = 0
	SizeOffset            = HeaderOffset + 0
	IDOffset              = HeaderOffset + 2
	SequenceOffset        = HeaderOffset + 4
	CRCOffset             = HeaderOffset + 5
	DataOffset            = HeaderOffset + HeaderSize

	EncryptOffset        = HeaderOffset + 2
	EncryptSize          = HeaderSize - EncryptOffset
	EncryptMask   uint16 = 0x8000
)

type Packet struct {
	MessageID             uint16
	Encrypted             bool
	Buffer                []byte
	Data                  []byte
	Size                  int
	DataSize              int
	Sequence              byte
	CRC                   byte
	ServerPacket          bool
	readPosition          int
	writePosition         int
	IsSequenceInitialized bool
}

func NewPacket(buffer []byte) Packet {
	messageID := uint16(0)
	encrypted := false
	crc := byte(0)
	sequence := byte(0)
	data := buffer

	if len(buffer) > int(IDOffset) {
		messageID = binary.LittleEndian.Uint16(buffer[IDOffset : IDOffset+2])
	}
	if len(buffer) > int(SequenceOffset) {
		sequence = buffer[SequenceOffset]
	}
	if len(buffer) > int(CRCOffset) {
		crc = buffer[CRCOffset]
	}
	if len(buffer) > int(DataOffset) {
		data = buffer[DataOffset:]
	}
	if len(buffer) > int(EncryptOffset) && (binary.LittleEndian.Uint16(buffer[:EncryptOffset])&EncryptMask) != 0 {
		encrypted = true
	}
	var size uint16

	if encrypted {
		size = uint16(buffer[0])
	} else {
		size = binary.LittleEndian.Uint16(buffer[:2])
	}

	return Packet{
		MessageID:     messageID,
		Encrypted:     encrypted,
		Buffer:        buffer,
		Data:          data,
		Size:          int(size + DataOffset),
		DataSize:      int(size),
		CRC:           crc,
		Sequence:      sequence,
		ServerPacket:  false,
		readPosition:  0,
		writePosition: 0,
	}
}

func EmptyPacket() Packet {
	buffer := make([]byte, 4096)
	return Packet{
		MessageID:     0,
		Encrypted:     false,
		Buffer:        buffer,
		Data:          buffer[DataOffset:],
		Size:          0,
		DataSize:      0,
		Sequence:      0,
		CRC:           0,
		ServerPacket:  true,
		readPosition:  0,
		writePosition: 0,
	}
}

func EmptyClientPacket() Packet {
	buffer := make([]byte, 4096)
	return Packet{
		MessageID:     0,
		Encrypted:     false,
		Buffer:        buffer,
		Data:          buffer[DataOffset:],
		Size:          0,
		DataSize:      0,
		Sequence:      0,
		CRC:           0,
		ServerPacket:  false,
		readPosition:  0,
		writePosition: 0,
	}
}

func (p *Packet) ReadByte() (byte, error) {
	if len(p.Data)-p.readPosition >= 1 {
		value := p.Data[p.readPosition]
		p.readPosition++
		return value, nil
	}
	return 0, errors.New("Not enough data available to read")
}

func (p *Packet) ReadUInt16() (uint16, error) {
	if len(p.Data)-p.readPosition >= 2 {
		value := p.Data[p.readPosition : p.readPosition+2]
		p.readPosition += 2
		return binary.LittleEndian.Uint16(value), nil
	}
	return 0, errors.New("Not enough data available to read")
}

func (p *Packet) ReadInt16() (int16, error) {
	uVal, err := p.ReadUInt16()
	if err != nil {
		return 0, err
	}
	return int16(uVal), nil
}

func (p *Packet) ReadUInt32() (uint32, error) {
	if len(p.Data)-p.readPosition >= 4 {
		value := p.Data[p.readPosition : p.readPosition+4]
		p.readPosition += 4
		return binary.LittleEndian.Uint32(value), nil
	}
	return 0, errors.New("Not enough data available to read")
}

func (p *Packet) ReadInt32() (int32, error) {
	uVal, err := p.ReadUInt32()
	if err != nil {
		return 0, err
	}
	return int32(uVal), nil
}

func (p *Packet) ReadUInt64() (uint64, error) {
	if len(p.Data)-p.readPosition >= 8 {
		value := p.Data[p.readPosition : p.readPosition+8]
		p.readPosition += 8
		return binary.LittleEndian.Uint64(value), nil
	}
	return 0, errors.New("Not enough data available to read")
}

func (p *Packet) ReadInt64() (int64, error) {
	uVal, err := p.ReadUInt64()
	if err != nil {
		return 0, err
	}
	return int64(uVal), nil
}

func (p *Packet) ReadString() (string, error) {
	stringLength, err := p.ReadUInt16()

	if err != nil {
		return "", err
	}

	if len(p.Data)-p.readPosition >= int(stringLength) {
		characters := p.Data[p.readPosition : p.readPosition+int(stringLength)]
		p.readPosition += int(stringLength)
		return string(characters), nil
	}
	return "0", errors.New("Not enough data available to read")
}

func (p *Packet) ReadFloat32() (float32, error) {
	if len(p.Data)-p.readPosition >= 4 {
		value := p.Data[p.readPosition : p.readPosition+4]
		p.readPosition += 4
		return utils.Float32FromByteArray(value), nil
	}
	return 0, errors.New("Not enough data available to read")
}

func (p *Packet) ReadFloat64() (float64, error) {
	if len(p.Data)-p.readPosition >= 8 {
		value := p.Data[p.readPosition : p.readPosition+8]
		p.readPosition += 8
		return utils.Float64FromByteArray(value), nil
	}
	return 0, errors.New("Not enough data available to read")
}

func (p *Packet) ReadBool() (bool, error) {
	b, err := p.ReadByte()
	if err != nil {
		return false, err
	}

	return b != 0, nil
}

func (p *Packet) ReadBytes(n int) ([]byte, error) {
	if len(p.Data)-p.readPosition >= n {
		value := p.Data[p.readPosition : p.readPosition+n]
		p.readPosition += n
		return value, nil
	}
	return nil, errors.New("Not enough data available to read")
}

func (p *Packet) WriteByte(value byte) {
	if cap(p.Data)-p.writePosition < 1 {
		p.extendBuffer()
	}

	p.Data[p.writePosition] = value
	p.writePosition++
	p.DataSize++
}

func (p *Packet) WriteUInt16(value uint16) {
	if cap(p.Data)-p.writePosition < 2 {
		p.extendBuffer()
	}
	p.Data[p.writePosition] = byte(value)
	p.Data[p.writePosition+1] = byte(value >> 8)
	p.writePosition += 2
	p.DataSize += 2
}

func (p *Packet) WriteUInt32(value uint32) {
	if cap(p.Data)-p.writePosition < 4 {
		p.extendBuffer()
	}
	p.Data[p.writePosition] = byte(value)
	p.Data[p.writePosition+1] = byte(value >> 8)
	p.Data[p.writePosition+2] = byte(value >> 16)
	p.Data[p.writePosition+3] = byte(value >> 24)
	p.writePosition += 4
	p.DataSize += 4
}

func (p *Packet) WriteUInt64(value uint64) {
	if cap(p.Data)-p.writePosition < 8 {
		p.extendBuffer()
	}
	p.Data[p.writePosition] = byte(value)
	p.Data[p.writePosition+1] = byte(value >> 8)
	p.Data[p.writePosition+2] = byte(value >> 16)
	p.Data[p.writePosition+3] = byte(value >> 24)
	p.Data[p.writePosition+4] = byte(value >> 32)
	p.Data[p.writePosition+5] = byte(value >> 40)
	p.Data[p.writePosition+6] = byte(value >> 48)
	p.Data[p.writePosition+7] = byte(value >> 56)
	p.writePosition += 8
	p.DataSize += 8
}

func (p *Packet) WriteString(value string) {
	stringLength := len(value)
	if cap(p.Data)-p.writePosition < stringLength+2 {
		p.extendBuffer()
		//p.WriteString(value)
	}
	p.WriteUInt16(uint16(stringLength))
	strBytes := []byte(value)
	for i := 0; i < stringLength; i++ {
		p.Data[p.writePosition+i] = strBytes[i]
	}
	p.writePosition += stringLength
	p.DataSize += stringLength
}

func (p *Packet) WriteFloat32(value float32) {
	if cap(p.Data)-p.writePosition < 4 {
		p.extendBuffer()
	}
	bits := math.Float32bits(value)
	p.Data[p.writePosition] = byte(bits)
	p.Data[p.writePosition+1] = byte(bits >> 8)
	p.Data[p.writePosition+2] = byte(bits >> 16)
	p.Data[p.writePosition+3] = byte(bits >> 24)
	p.writePosition += 4
	p.DataSize += 4
}

func (p *Packet) WriteFloat64(value float64) {
	if cap(p.Data)-p.writePosition < 8 {
		p.extendBuffer()
	}
	bits := math.Float64bits(value)
	p.Data[p.writePosition] = byte(bits)
	p.Data[p.writePosition+1] = byte(bits >> 8)
	p.Data[p.writePosition+2] = byte(bits >> 16)
	p.Data[p.writePosition+3] = byte(bits >> 24)
	p.Data[p.writePosition+4] = byte(bits >> 32)
	p.Data[p.writePosition+5] = byte(bits >> 40)
	p.Data[p.writePosition+6] = byte(bits >> 48)
	p.Data[p.writePosition+7] = byte(bits >> 56)
	p.writePosition += 8
	p.DataSize += 8
}

func (p *Packet) WriteBool(value bool) {
	if value {
		p.WriteByte(1)
	} else {
		p.WriteByte(0)
	}
}

func (p *Packet) WriteBytes(bytes []byte) {
	if cap(p.Data)-p.writePosition < len(bytes) {
		p.extendBuffer()
	}
	copy(p.Data[p.writePosition:p.writePosition+len(bytes)], bytes)
}

func (p *Packet) WriteMessageID() {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, p.MessageID)
	p.Buffer[IDOffset] = bytes[0]
	p.Buffer[IDOffset+1] = bytes[1]
}

func (p *Packet) WriteCRC() {
	if p.ServerPacket {
		p.Buffer[CRCOffset] = 0x00
	} else {
		p.Buffer[CRCOffset] = p.CRC
	}
}

func (p *Packet) WriteSequence() {
	if p.ServerPacket {
		p.Buffer[SequenceOffset] = 0x00
	} else if p.IsSequenceInitialized {
		p.Buffer[SequenceOffset] = p.Sequence
	} else {
		log.Panic("sequence has not been initialized")
	}
}

func (p *Packet) WriteSize() {
	size := uint16(p.DataSize)

	if p.Encrypted {
		size |= uint16(1) << 15
	}

	p.Buffer[0] = byte(size)
	p.Buffer[1] = byte(size >> 8)
}

func (p *Packet) ToBytes() []byte {
	p.WriteSize()
	p.WriteMessageID()
	p.WriteSequence()
	p.WriteCRC()
	return p.Buffer[:p.DataSize+int(HeaderSize)]
}

func (p *Packet) getBufferForBlowfishOperation() []byte {
	numOfBytesForBlowfishOp := p.DataSize + int(EncryptSize)
	numOfBytesAfterBlowfishOp := (numOfBytesForBlowfishOp + 7) & -8
	if uint16(numOfBytesAfterBlowfishOp)+EncryptOffset > BufferSize {
		log.Errorf("Size of packet to de/encrypt is too large: %v\n", numOfBytesForBlowfishOp)
	}
	return p.Buffer[EncryptOffset : int(EncryptOffset)+numOfBytesAfterBlowfishOp]
}

func (p *Packet) Decrypt(c *blowfish.Cipher) Packet {
	if p.Encrypted {
		bufferToDecrypt := p.getBufferForBlowfishOperation()
		for i := 1; i <= len(bufferToDecrypt)/8; i++ {
			chunkToDecrypt := bufferToDecrypt[(i-1)*8 : i*8]
			c.DecryptRev(chunkToDecrypt, chunkToDecrypt)
		}
		return NewPacket(p.Buffer[:p.DataSize+int(HeaderSize)])
	}
	return *p
}

func (p *Packet) Encrypt(c *blowfish.Cipher) Packet {
	if p.Encrypted {
		if p.ServerPacket {
			p.Buffer[SequenceOffset] = 0
			p.Buffer[CRCOffset] = 0
		}
		bufferToEncrypt := p.getBufferForBlowfishOperation()
		for i := 1; i <= len(bufferToEncrypt)/8; i++ {
			chunkToEncrypt := bufferToEncrypt[(i-1)*8 : i*8]
			c.EncryptRev(chunkToEncrypt, chunkToEncrypt)
		}
		p.DataSize = len(bufferToEncrypt) - int(EncryptSize)
	}
	return *p
}

func (p *Packet) IsMassive() bool {
	return p.DataSize > 4090 || p.Size > 4096
}

func (p *Packet) SplitIntoMassivePackets() []*Packet {
	if p.IsMassive() {
		var packets []*Packet
		numOfBodyPackets := uint16(math.Ceil(float64(len(p.Data)) / 4090))
		headerPacket := EmptyPacket()
		headerPacket.MessageID = 0x600D
		headerPacket.WriteBool(true)
		headerPacket.WriteUInt16(numOfBodyPackets)
		headerPacket.WriteUInt16(p.MessageID)
		packets = append(packets, &headerPacket)
		for i := 1; i < int(numOfBodyPackets); i++ {
			packet := EmptyPacket()
			packet.MessageID = 0x600D
			packet.WriteBool(false)
			copy(packet.Data[1:4090], p.Data[i*4089:(i+1)*4089])
			packets = append(packets, &packet)
		}
		return packets

	} else {
		return []*Packet{p}
	}
}

func (p *Packet) String() string {
	sb := strings.Builder{}
	sb2 := strings.Builder{}
	sb.WriteString("\nPacket\n")
	sb.WriteString(fmt.Sprintf("[%X]\n", p.MessageID))
	sb.WriteString(fmt.Sprintf("Encrypted [%v]\n", p.Encrypted))
	sb.WriteString(fmt.Sprintf("Massive [%v]\n", p.IsMassive()))
	sb.WriteString(fmt.Sprintf("Sequence [%v]\n", p.Sequence))
	sb.WriteString(fmt.Sprintf("CRC [%v]\n", p.CRC))
	sb.WriteString(fmt.Sprintf("DataSize [%v]\n", p.DataSize))
	sb.WriteString(fmt.Sprintf("Size [%v]\n", p.Size))
	for i := 0; i < p.DataSize; i++ {
		sb.WriteString(fmt.Sprintf("%02X ", p.Data[i]))
		sb2.WriteString(dumpByte(p.Data[i]))
		if (i+1)%8 == 0 {
			sb.WriteString(fmt.Sprintf("\t\t[%8s]\n", sb2.String()))
			sb2.Reset()
		} else if (i + 1) == p.DataSize {
			m := 8 - ((i + 1) % 8)
			for j := 0; j < m; j++ {
				sb.WriteString("   ")
			}
			sb.WriteString(fmt.Sprintf("\t\t[%-8s]\n", sb2.String()))
			sb2.Reset()
		}
	}
	return sb.String()
}

func dumpByte(value byte) string {
	if value < 32 || value > 126 {
		return "."
	}
	return string(value)
}

func (p *Packet) extendBuffer() {
	newBuffer := make([]byte, cap(p.Data)*2)
	copy(newBuffer, p.Data)
	p.Data = newBuffer
}
