package pk2

import "time"

const (
	EntrySize = 128
	Key       = "169841"
	BlockSize = 20 * EntrySize
	TypeEmpty = 0
	TypeDir   = 1
	TypeFile  = 2

	HeaderSize            = 256
	HeaderHeaderOffset    = 0
	HeaderVersionOffset   = HeaderHeaderOffset + 30
	HeaderEncryptedOffset = HeaderVersionOffset + 4
	HeaderChecksumOffset  = HeaderEncryptedOffset + 1
	HeaderReservedOffset  = HeaderChecksumOffset + 16

	EntryTypeOffset       = 0
	EntryNameOffset       = EntryTypeOffset + 1
	EntryAccessTimeOffset = EntryNameOffset + 81
	EntryCreateTimeOffset = EntryAccessTimeOffset + 8
	EntryModifyTimeOffset = EntryCreateTimeOffset + 8
	EntryPositionOffset   = EntryModifyTimeOffset + 8
	EntrySizeOffset       = EntryPositionOffset + 8
	EntryNextChainOffset  = EntrySizeOffset + 4
	EntryPaddingOffset    = EntryNextChainOffset + 8
)

var BaseKey = []byte{0x03, 0xF8, 0xE4, 0x44, 0x88, 0x99, 0x3F, 0x64, 0xFE, 0x35}

type PackHeader struct {
	Header    string
	Version   uint32
	Encrypted bool
	Checksum  []byte
	Reserverd []byte
}

type PackFileEntry struct {
	Type       byte
	Name       string
	AccessTime time.Time
	CreateTime time.Time
	ModifyTime time.Time
	Position   uint64
	Size       uint32
	NextChain  uint64
	Padding    []byte
}

type PackFileBlock struct {
	Entries []PackFileEntry
}
