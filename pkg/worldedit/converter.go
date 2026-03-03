package worldedit

import (
	"fmt"
	"strings"

	"github.com/ntaku256/go-java-nbt-converter/pkg/format"
)

// Convert transforms a parsed WorldEditNBT into a StandardFormat.
func Convert(we *WorldEditNBT) (*format.StandardFormat, error) {
	if we == nil {
		return nil, fmt.Errorf("worldEdit data is nil")
	}

	sf := &format.StandardFormat{
		OriginalFormat: "worldedit",
		DataVersion:    int(we.DataVersion),
		Version:        int(we.Version),
		Blocks:         make([]format.StandardBlock, 0),
	}

	sf.Size.X = int(we.Width)
	sf.Size.Y = int(we.Height)
	sf.Size.Z = int(we.Length)

	if len(we.Offset) >= 3 {
		sf.Position.X = int(we.Offset[0])
		sf.Position.Y = int(we.Offset[1])
		sf.Position.Z = int(we.Offset[2])
	}

	// Convert palette ("minecraft:stone[variant=granite]" → name + properties)
	sf.Palette = make(map[int]format.StandardPalette, len(we.Palette))
	for name, index := range we.Palette {
		blockName, properties := parseBlockString(name)
		sf.Palette[int(index)] = format.StandardPalette{
			Name:       blockName,
			Properties: properties,
		}
	}

	// Unpack BlockData (VarInt-encoded int8 array)
	// BlockData is now []int8 from direct NBT decode
	if len(we.BlockData) > 0 {
		// Convert []int8 to []byte for VarInt reading
		blockDataBytes := make([]byte, len(we.BlockData))
		for i, b := range we.BlockData {
			blockDataBytes[i] = byte(b)
		}

		offset := 0
		totalBlocks := int(we.Width) * int(we.Height) * int(we.Length)
		sf.Blocks = make([]format.StandardBlock, 0, totalBlocks)

		for y := 0; y < int(we.Height); y++ {
			for z := 0; z < int(we.Length); z++ {
				for x := 0; x < int(we.Width); x++ {
					if offset >= len(blockDataBytes) {
						break
					}
					paletteIndex := readVarInt(blockDataBytes, &offset)

					if paletteIndex >= 0 {
						p, exists := sf.Palette[paletteIndex]
						if !exists || p.Name != "minecraft:air" {
							sf.Blocks = append(sf.Blocks, format.StandardBlock{
								Position: format.StandardBlockPosition{
									X: float64(x + sf.Position.X),
									Y: float64(y + sf.Position.Y),
									Z: float64(z + sf.Position.Z),
								},
								State: paletteIndex,
								Type:  "block",
							})
						}
					}
				}
			}
		}
	}

	// Block entities
	for _, beRaw := range we.BlockEntities {
		if beMap, ok := beRaw.(map[string]interface{}); ok {
			id := ""
			if v, ok := beMap["Id"].(string); ok {
				id = v
			} else if v, ok := beMap["id"].(string); ok {
				id = v
			}
			if id != "" {
				x, y, z := extractBlockEntityPosition(beMap)
				sf.Blocks = append(sf.Blocks, format.StandardBlock{
					Position: format.StandardBlockPosition{X: x, Y: y, Z: z},
					ID:       id,
					Type:     "block_entity",
					NBT:      beMap,
				})
			}
		}
	}

	return sf, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseBlockString(s string) (string, map[string]string) {
	parts := strings.SplitN(s, "[", 2)
	name := parts[0]
	props := make(map[string]string)
	if len(parts) > 1 {
		propsStr := strings.TrimSuffix(parts[1], "]")
		for _, p := range strings.Split(propsStr, ",") {
			kv := strings.SplitN(p, "=", 2)
			if len(kv) == 2 {
				props[kv[0]] = kv[1]
			}
		}
	}
	return name, props
}

func readVarInt(data []byte, offset *int) int {
	result := 0
	shift := 0
	for {
		if *offset >= len(data) {
			return -1
		}
		b := data[*offset]
		*offset++
		result |= int(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return result
}

func extractBlockEntityPosition(be map[string]interface{}) (float64, float64, float64) {
	getFloat := func(key string) float64 {
		switch v := be[key].(type) {
		case float64:
			return v
		case int32:
			return float64(v)
		case int64:
			return float64(v)
		case int:
			return float64(v)
		}
		return 0
	}
	return getFloat("x"), getFloat("y"), getFloat("z")
}
