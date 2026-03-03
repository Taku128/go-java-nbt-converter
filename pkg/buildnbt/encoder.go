package buildnbt

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"

	"github.com/Tnze/go-mc/nbt"
	"github.com/ntaku256/go-java-nbt-converter/pkg/format"
)

// Encode takes a StandardFormat and produces a gzip-compressed Java Edition
// Structure NBT byte slice, ready to be loaded by Minecraft or viewers.
func Encode(sf *format.StandardFormat) ([]byte, error) {
	jNbt := JavaNBT{
		DataVersion: int32(sf.DataVersion),
		Size:        ListTag{int32(sf.Size.X), int32(sf.Size.Y), int32(sf.Size.Z)},
		Entities:    make([]interface{}, 0),
	}

	// Ensure no zero-sized axis
	if jNbt.Size[0] == 0 {
		jNbt.Size[0] = 1
	}
	if jNbt.Size[1] == 0 {
		jNbt.Size[1] = 1
	}
	if jNbt.Size[2] == 0 {
		jNbt.Size[2] = 1
	}

	// Build palette (sparse map → dense slice)
	maxPalette := -1
	for k := range sf.Palette {
		if k > maxPalette {
			maxPalette = k
		}
	}
	if maxPalette >= 0 {
		jNbt.Palette = make([]PaletteEntry, maxPalette+1)
		for k, v := range sf.Palette {
			jNbt.Palette[k] = PaletteEntry{
				Name:       v.Name,
				Properties: v.Properties,
			}
		}
	} else {
		jNbt.Palette = []PaletteEntry{{Name: "minecraft:air"}}
	}

	// Build blocks (relative to Position origin)
	jNbt.Blocks = make([]JavaBlock, 0, len(sf.Blocks))
	for _, b := range sf.Blocks {
		if b.Type == "block" || b.Type == "" {
			jNbt.Blocks = append(jNbt.Blocks, JavaBlock{
				Pos: ListTag{
					int32(b.Position.X) - int32(sf.Position.X),
					int32(b.Position.Y) - int32(sf.Position.Y),
					int32(b.Position.Z) - int32(sf.Position.Z),
				},
				State: int32(b.State),
				Nbt:   b.NBT,
			})
		}
	}

	log.Printf("buildnbt.Encode: DataVersion=%d, Size=%v, Palette=%d, Blocks=%d",
		jNbt.DataVersion, jNbt.Size, len(jNbt.Palette), len(jNbt.Blocks))

	// Encode to gzip-compressed NBT
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if err := nbt.NewEncoder(gw).Encode(jNbt, ""); err != nil {
		gw.Close()
		return nil, fmt.Errorf("failed to encode NBT: %w", err)
	}
	gw.Close()

	return buf.Bytes(), nil
}
