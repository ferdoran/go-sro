package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ferdoran/go-sro/framework/security/blowfish"
	"io"
)

type PacketReader struct {
	reader io.Reader
}

func NewPacketReader(reader io.Reader) *PacketReader {
	return &PacketReader{
		reader: reader,
	}
}

func (p *PacketReader) ReadPackets() ([]Packet, error) {
	buffer := make([]byte, BufferSize)
	bytesRead, err := p.reader.Read(buffer)
	if err != nil {
		return nil, err
	}
	packets := make([]Packet, 0)
	buffer = buffer[:bytesRead]
	for len(buffer) > 0 {
		packetDataSize := binary.LittleEndian.Uint16(buffer[:2])
		//encrypted := false

		size := int(packetDataSize + HeaderSize)
		if packetDataSize&EncryptMask != 0 {
			// encrypted
			//encrypted = true
			packetDataSize = uint16(buffer[0])
			size = blowfish.GetBufferLength(int(packetDataSize+EncryptSize)) + int(EncryptOffset)
		}
		if size <= len(buffer) {
			// we have multiple packets here
			packets = append(packets, NewPacket(buffer[:size]))
			buffer = buffer[size:]
		} else {
			return nil, errors.New(fmt.Sprintf("packet size is greater than available bytes in buffer. size = %d, avail = %d", size, len(buffer)))
		}
	}
	return packets, nil

}
