package pk2

import (
	"fmt"
	"github.com/ferdoran/go-sro/framework/security/blowfish"
	"github.com/ferdoran/go-sro/framework/utils"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	checksum = "Joymax Pak File"
)

type Pk2Reader struct {
	file             *os.File
	cipher           *blowfish.Cipher
	Directory        Directory
	finishedIndexing bool
	Files            map[string]PackFileEntry
	FileCache        map[string][]byte
	CachingEnabled   bool
}

func NewPk2Reader(filename string) Pk2Reader {
	f, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	cipher, err := blowfish.NewCipher(GeneratePk2BlowfishKey())
	if err != nil {
		log.Panic(err)
	}
	fileName := filepath.Base(filename)
	fileExt := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(fileExt)]
	return Pk2Reader{file: f, cipher: cipher, Directory: Directory{Name: fileName}, CachingEnabled: true, FileCache: make(map[string][]byte)}
}

func (r *Pk2Reader) IndexArchive() {
	header := r.readHeader()
	r.verifyHeader(header)
	log.Println("Read header successfully")
	r.readEntries(HeaderSize, &r.Directory)
	r.Directory.buildEntryMap()
	r.Directory.buildDirectoryMap()
	r.Files = r.Directory.AllFiles()
	r.finishedIndexing = true
	log.Println("Read entries successfully")
}

func (r *Pk2Reader) ReadFile(filename string) ([]byte, error) {
	filename = strings.ReplaceAll(filename, "\\", string(os.PathSeparator))
	filename = strings.ReplaceAll(filename, "/", string(os.PathSeparator))
	if cacheData, ok := r.FileCache[filename]; r.CachingEnabled && ok {
		return cacheData, nil
	}

	if entry, ok := r.Files[filename]; ok {
		data := r.ReadEntry(&entry)
		if r.CachingEnabled {
			r.FileCache[filename] = data
		}
		return data, nil
	}

	return nil, errors.New(fmt.Sprintf("file not found: %s", filename))
}

func (r *Pk2Reader) ClearFileCache() {
	r.FileCache = make(map[string][]byte)
}

func (r *Pk2Reader) readHeader() PackHeader {
	headerBuffer := make([]byte, HeaderSize)
	bytesRead, err := r.file.Read(headerBuffer)

	if err != nil {
		log.Panic(err)
	} else if bytesRead != HeaderSize {
		log.Panicf("Failed to read pk2 header. Want = %v, got %v\n", HeaderSize, bytesRead)
	}

	enc := false
	if headerBuffer[HeaderEncryptedOffset] == 1 {
		enc = true
	}
	return PackHeader{
		Header:    string(headerBuffer[HeaderHeaderOffset:HeaderVersionOffset]),
		Version:   utils.ByteArrayToUint32(headerBuffer[HeaderVersionOffset:HeaderEncryptedOffset]),
		Encrypted: enc,
		Checksum:  headerBuffer[HeaderChecksumOffset:HeaderReservedOffset],
		Reserverd: headerBuffer[HeaderReservedOffset:HeaderSize],
	}
}

func (r *Pk2Reader) readEntries(startPosition int64, directory *Directory) {
	entryBuffer := make([]byte, 2560)
	bytesRead, err := r.file.ReadAt(entryBuffer, startPosition)

	if err != nil {
		log.Panic(err)
	}

	if bytesRead%EntrySize != 0 {
		log.Panic("Entries length is not divisible by 128")
	}
	r.decryptBuffer(entryBuffer)

	entries := make([]PackFileEntry, 0)
	for i := 0; i < bytesRead/EntrySize; i++ {
		offset := i * EntrySize
		newEntry := PackFileEntry{
			Type:       entryBuffer[offset+EntryTypeOffset],
			Name:       strings.Trim(string(entryBuffer[offset+EntryNameOffset:offset+EntryAccessTimeOffset]), "\x00"),
			AccessTime: time.Unix(int64(utils.ByteArrayToUint64(entryBuffer[offset+EntryAccessTimeOffset:offset+EntryCreateTimeOffset])), 0),
			CreateTime: time.Unix(int64(utils.ByteArrayToUint64(entryBuffer[offset+EntryCreateTimeOffset:offset+EntryModifyTimeOffset])), 0),
			ModifyTime: time.Unix(int64(utils.ByteArrayToUint64(entryBuffer[offset+EntryModifyTimeOffset:offset+EntryPositionOffset])), 0),
			Position:   utils.ByteArrayToUint64(entryBuffer[offset+EntryPositionOffset : offset+EntrySizeOffset]),
			Size:       utils.ByteArrayToUint32(entryBuffer[offset+EntrySizeOffset : offset+EntryNextChainOffset]),
			NextChain:  utils.ByteArrayToUint64(entryBuffer[offset+EntryNextChainOffset : offset+EntryPaddingOffset]),
			Padding:    entryBuffer[offset+EntryPaddingOffset : offset+EntryPaddingOffset+2],
		}
		newEntry.Name = strings.ReplaceAll(newEntry.Name, "\\", string(os.PathSeparator))
		entries = append(entries, newEntry)
	}

	if nextChain := entries[19].NextChain; nextChain > 0 {
		r.readEntries(int64(nextChain), directory)
	}

	directory.Entries = append(directory.Entries, entries...)

	if directory.Directories == nil {
		directory.Directories = make([]Directory, 0)
	}
	for _, v := range entries {
		if v.Type == TypeDir && !strings.HasPrefix(v.Name, ".") {
			// Expand it
			newDir := Directory{Name: directory.Name + string(os.PathSeparator) + v.Name, Entries: make([]PackFileEntry, 0), Files: make([]PackFileEntry, 0)}
			r.readEntries(int64(v.Position), &newDir)
			directory.Directories = append(directory.Directories, newDir)
		}
	}

}

func GeneratePk2BlowfishKey() []byte {
	keyArr := []byte(Key)
	blowfishKey := make([]byte, len(Key))

	for x := 0; x < len(Key); x++ {
		blowfishKey[x] = keyArr[x] ^ BaseKey[x]
	}
	return blowfishKey
}

func (r *Pk2Reader) verifyHeader(header PackHeader) {
	if header.Encrypted {
		// Decrypt it
		encodedChecksum := []byte(checksum)
		r.cipher.EncryptRev(encodedChecksum, []byte(checksum))

		// Only the first 3 bytes are equal
		for i := 0; i < 3; i++ {
			if encodedChecksum[i] != header.Checksum[i] {
				log.Panicf("Checksum is invalid. Want = %v, got %v\n", encodedChecksum[:3], header.Checksum[:3])
			}
		}
	}
}

func (r *Pk2Reader) decryptBuffer(buffer []byte) {
	for i := 0; i < len(buffer)/8; i++ {
		bufToDecrypt := buffer[i*8 : (i+1)*8]
		r.cipher.DecryptRev(bufToDecrypt, bufToDecrypt)
	}
}

func (r *Pk2Reader) ReadEntry(entry *PackFileEntry) []byte {
	buffer := make([]byte, entry.Size)
	bytesRead, err := r.file.ReadAt(buffer, int64(entry.Position))

	if err != nil {
		log.Panicln(err)
	}

	if bytesRead != int(entry.Size) {
		log.Panicf("Read an incorrect amount of bytes. Want = %v, got %v\n", entry.Size, bytesRead)
	}
	return buffer
}

func (r *Pk2Reader) ExtractFiles(outputDir string) {
	if r.finishedIndexing {
		os.Mkdir(outputDir, os.ModePerm)
		counter := 0
		r.Directory.Extract(outputDir, *r, &counter, r.TotalFiles())
	}
}

func (r *Pk2Reader) TotalFiles() int {
	if r.finishedIndexing {
		return r.Directory.TotalFiles()
	}
	return 0
}

func (r *Pk2Reader) LoadFile(filename string) ([]byte, error) {
	return r.Directory.getFile(r, filename)
}
