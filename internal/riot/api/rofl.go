package riot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

type Chunk struct {
	ID              uint32
	Type            uint8
	ID2             uint32
	UncompressedLen uint32
	CompressedLen   uint32
	Payload         []byte
}

func ParseROFL(path string) ([]Chunk, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(raw) < 4 {
		return nil, errors.New("file too short")
	}
	// --- Remove metadata + signature ---
	metaLen := binary.LittleEndian.Uint32(raw[len(raw)-4:])
	if len(raw) < int(metaLen)+4 {
		return nil, errors.New("file is corrupted (metadata length too large)")
	}
	raw = raw[:len(raw)-int(metaLen)-4]
	if len(raw) < 0x100 {
		return nil, errors.New("file too short after signature removal")
	}
	raw = raw[:len(raw)-0x100]

	// --- Remove ROFL header ---
	if len(raw) < 0x10 {
		return nil, errors.New("file too short for header removal")
	}
	raw = raw[0x10:]
	if len(raw) < 0xD {
		return nil, errors.New("file too short for secondary header removal")
	}
	if raw[0xC] == 1 {
		if len(raw) < 0xC {
			return nil, errors.New("header step exceeds buffer")
		}
		raw = raw[0xC:]
	} else {
		if len(raw) < 0xD {
			return nil, errors.New("header step exceeds buffer")
		}
		raw = raw[0xD:]
	}

	// Parse Chunks
	var chunks []Chunk
	buf := bytes.NewReader(raw)

	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init zstd decoder: %w", err)
	}
	defer decoder.Close()

	for {
		// Check if enough data remains for a chunk header (4+1+4+4+4=17 bytes)
		header := make([]byte, 17)
		n, err := buf.Read(header)
		if err == io.EOF || n < 17 {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading chunk header: %w", err)
		}
		chunk := Chunk{}
		chunk.ID = binary.LittleEndian.Uint32(header[0:4])
		chunk.Type = header[4]
		chunk.ID2 = binary.LittleEndian.Uint32(header[5:9])
		chunk.UncompressedLen = binary.LittleEndian.Uint32(header[9:13])
		chunk.CompressedLen = binary.LittleEndian.Uint32(header[13:17])

		// --- Read payload ---
		if chunk.CompressedLen != 0 {
			compressed := make([]byte, chunk.CompressedLen)
			_, err = io.ReadFull(buf, compressed)
			if err != nil {
				return nil, fmt.Errorf("reading compressed payload: %w", err)
			}
			payload, err := decoder.DecodeAll(compressed, nil)
			if err != nil {
				return nil, fmt.Errorf("zstd decompress failed: %w", err)
			}
			chunk.Payload = payload
		} else if chunk.UncompressedLen != 0 {
			// There might be uncompressed payload; usually it's not present
			_, err := buf.Seek(int64(chunk.UncompressedLen), io.SeekCurrent)
			if err != nil {
				return nil, fmt.Errorf("seeking past uncompressed payload: %w", err)
			}
			// chunk.Payload = nil
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func main() {
	chunks, err := ParseROFL("input.rofl")
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	fmt.Printf("Parsed %d chunks\n", len(chunks))
	for i, ch := range chunks {
		fmt.Printf("Chunk %d: ID=%d Type=%d PayloadLen=%d\n",
			i+1, ch.ID, ch.Type, len(ch.Payload))
		// You can process each chunk.Payload as needed
	}
}
