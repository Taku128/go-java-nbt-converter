// Package structure handles Java Edition Structure .nbt files.
package structure

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/Tnze/go-mc/nbt"
	"github.com/ntaku256/go-java-nbt-converter/pkg/format"
)

// StructureNBT is the top-level structure of a Java Edition .nbt file.
type StructureNBT struct {
	DataVersion int32                    `nbt:"DataVersion"`
	Size        []int32                  `nbt:"size"`
	Blocks      []StructureBlock         `nbt:"blocks"`
	Entities    []interface{}            `nbt:"entities"`
	Palette     []StructurePaletteEntry  `nbt:"palette"`
}

// StructureBlock is a placed block.
type StructureBlock struct {
	Pos   []int32     `nbt:"pos"`
	State int32       `nbt:"state"`
	Nbt   interface{} `nbt:"nbt,omitempty"`
}

// StructurePaletteEntry is a palette entry.
type StructurePaletteEntry struct {
	Name       string            `nbt:"Name"`
	Properties map[string]string `nbt:"Properties"`
}

// ParseFile reads and decodes a .nbt structure file directly into typed structs.
func ParseFile(filePath string) (*StructureNBT, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return Parse(data)
}

// Parse decodes raw bytes into a StructureNBT using direct NBT decoding.
func Parse(data []byte) (*StructureNBT, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer r.Close()

	var s StructureNBT
	if _, err = nbt.NewDecoder(r).Decode(&s); err != nil {
		return nil, fmt.Errorf("failed to decode NBT: %w", err)
	}

	return &s, nil
}

// Convert transforms a StructureNBT into StandardFormat.
func Convert(s *StructureNBT) (*format.StandardFormat, error) {
	if s == nil {
		return nil, fmt.Errorf("structure data is nil")
	}

	sf := &format.StandardFormat{
		OriginalFormat: "structure",
		DataVersion:    int(s.DataVersion),
		Blocks:         make([]format.StandardBlock, 0, len(s.Blocks)),
		Palette:        make(map[int]format.StandardPalette, len(s.Palette)),
	}

	if len(s.Size) >= 3 {
		sf.Size.X = int(s.Size[0])
		sf.Size.Y = int(s.Size[1])
		sf.Size.Z = int(s.Size[2])
	}

	for i, p := range s.Palette {
		sf.Palette[i] = format.StandardPalette{
			Name:       p.Name,
			Properties: p.Properties,
		}
	}

	for _, b := range s.Blocks {
		if len(b.Pos) < 3 {
			continue
		}
		sf.Blocks = append(sf.Blocks, format.StandardBlock{
			Position: format.StandardBlockPosition{
				X: float64(b.Pos[0]),
				Y: float64(b.Pos[1]),
				Z: float64(b.Pos[2]),
			},
			State: int(b.State),
			Type:  "block",
			NBT:   b.Nbt,
		})
	}

	return sf, nil
}
