package worldedit

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/Tnze/go-mc/nbt"
)

// ParseFile reads and decodes a .schem file directly into typed structs.
func ParseFile(filePath string) (*WorldEditNBT, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return Parse(data)
}

// Parse decodes raw (gzip-compressed) NBT bytes directly into a WorldEditNBT.
// No JSON round-trip — preserves all NBT type information.
func Parse(data []byte) (*WorldEditNBT, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer r.Close()

	var we WorldEditNBT
	if _, err = nbt.NewDecoder(r).Decode(&we); err != nil {
		return nil, fmt.Errorf("failed to decode NBT: %w", err)
	}

	return &we, nil
}
