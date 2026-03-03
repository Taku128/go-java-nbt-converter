package litematica

import (
	"fmt"
	"log"
	"math/bits"
	"sort"
	"strings"

	"github.com/ntaku256/go-java-nbt-converter/pkg/format"
)

// Convert transforms a parsed LitematicaNBT into a StandardFormat.
// This now works with properly-typed fields (int32, int64) instead of
// the broken float64 values that came from JSON round-tripping.
func Convert(lit *LitematicaNBT) (*format.StandardFormat, error) {
	if lit == nil {
		return nil, fmt.Errorf("litematica data is nil")
	}

	sf := &format.StandardFormat{
		OriginalFormat: "litematica",
		DataVersion:    int(lit.MinecraftDataVersion),
		Version:        int(lit.Version),
		Blocks:         make([]format.StandardBlock, 0),
		Palette:        make(map[int]format.StandardPalette),
	}

	// Metadata
	sf.Metadata.Name = lit.Metadata.Name
	sf.Metadata.Author = lit.Metadata.Author
	sf.Metadata.Description = lit.Metadata.Description
	sf.Metadata.TimeCreated = lit.Metadata.TimeCreated
	sf.Metadata.TimeModified = lit.Metadata.TimeModified
	sf.Metadata.TotalBlocks = int(lit.Metadata.TotalBlocks)
	sf.Metadata.TotalVolume = int(lit.Metadata.TotalVolume)

	// Global palette (merging from multiple regions)
	globalPaletteMap := make(map[string]int)

	minX, minY, minZ := int32(1<<31-1), int32(1<<31-1), int32(1<<31-1)
	maxX, maxY, maxZ := int32(-1<<31), int32(-1<<31), int32(-1<<31)
	foundAny := false

	for regionName, region := range lit.Regions {
		log.Printf("Processing Litematica region: %s (Size: %d,%d,%d, Pos: %d,%d,%d)",
			regionName, region.Size.X, region.Size.Y, region.Size.Z,
			region.Position.X, region.Position.Y, region.Position.Z)

		// 1. Map local palette to global
		localToGlobal := make(map[int]int)
		for i, lp := range region.BlockStatePalette {
			key := canonicalBlockKey(lp.Name, lp.Properties)
			gIdx, exists := globalPaletteMap[key]
			if !exists {
				gIdx = len(sf.Palette)
				globalPaletteMap[key] = gIdx
				sf.Palette[gIdx] = format.StandardPalette{
					Name:       lp.Name,
					Properties: lp.Properties,
				}
			}
			localToGlobal[i] = gIdx
		}

		// 2. Unpack blocks (using proper int32 types, no float64 conversion)
		rSizeX := abs32(region.Size.X)
		rSizeY := abs32(region.Size.Y)
		rSizeZ := abs32(region.Size.Z)

		// Compute region bounds for overall size calculation
		rPosX := region.Position.X
		if region.Size.X < 0 {
			rPosX += region.Size.X + 1
		}
		rPosY := region.Position.Y
		if region.Size.Y < 0 {
			rPosY += region.Size.Y + 1
		}
		rPosZ := region.Position.Z
		if region.Size.Z < 0 {
			rPosZ += region.Size.Z + 1
		}

		if rPosX < minX { minX = rPosX }
		if rPosY < minY { minY = rPosY }
		if rPosZ < minZ { minZ = rPosZ }
		if rPosX+rSizeX > maxX { maxX = rPosX + rSizeX }
		if rPosY+rSizeY > maxY { maxY = rPosY + rSizeY }
		if rPosZ+rSizeZ > maxZ { maxZ = rPosZ + rSizeZ }
		foundAny = true

		numPaletteEntries := len(region.BlockStatePalette)
		if numPaletteEntries == 0 {
			continue
		}

		bitsPerBlock := bits.Len(uint(numPaletteEntries - 1))
		if bitsPerBlock < 2 {
			bitsPerBlock = 2
		}
		mask := uint64((1 << bitsPerBlock) - 1)

		// BlockStates is now []int64 decoded directly from TAG_Long_Array
		// No more float64 precision loss!
		longData := region.BlockStates
		if len(longData) == 0 {
			log.Printf("  Region %s has no BlockStates data", regionName)
			continue
		}

		log.Printf("  Region %s: palette=%d, bitsPerBlock=%d, longs=%d",
			regionName, numPaletteEntries, bitsPerBlock, len(longData))

		for y := int32(0); y < rSizeY; y++ {
			for z := int32(0); z < rSizeZ; z++ {
				for x := int32(0); x < rSizeX; x++ {
					index := int(y)*int(rSizeX)*int(rSizeZ) + int(z)*int(rSizeX) + int(x)

					// Litematica coordinate mapping:
					// index 0 starts at Position and grows in the direction of Size.
					realX := region.Position.X
					if region.Size.X >= 0 { realX += x } else { realX -= x }
					realY := region.Position.Y
					if region.Size.Y >= 0 { realY += y } else { realY -= y }
					realZ := region.Position.Z
					if region.Size.Z >= 0 { realZ += z } else { realZ -= z }

					bitOffset := index * bitsPerBlock
					startLong := bitOffset / 64
					startBit := uint(bitOffset % 64)

					var paletteIndex uint64
					if startLong < len(longData) {
						val := uint64(longData[startLong]) >> startBit
						if startBit+uint(bitsPerBlock) > 64 && startLong+1 < len(longData) {
							val |= uint64(longData[startLong+1]) << (64 - startBit)
						}
						paletteIndex = val & mask
					}

					gIdx := localToGlobal[int(paletteIndex)]
					// Skip air blocks
					if gIdx > 0 || sf.Palette[gIdx].Name != "minecraft:air" {
						sf.Blocks = append(sf.Blocks, format.StandardBlock{
							Position: format.StandardBlockPosition{
								X: float64(realX),
								Y: float64(realY),
								Z: float64(realZ),
							},
							State: gIdx,
							Type:  "block",
						})
					}
				}
			}
		}

		// 3. Tile entities (with preserved raw NBT data)
		for _, te := range region.TileEntities {
			posX := float64(te.X + region.Position.X)
			posY := float64(te.Y + region.Position.Y)
			posZ := float64(te.Z + region.Position.Z)

			// Build NBT data from the preserved Extra fields
			nbtData := make(map[string]interface{})
			if te.Id != "" {
				nbtData["id"] = te.Id
			}
			// Copy all extra fields (Items, CookingTimes, Bees, custom mod data, etc.)
			// These are preserved as their raw NBT types, not JSON-converted types.
			for k, v := range te.Extra {
				nbtData[k] = v
			}

			found := false
			for i := range sf.Blocks {
				if sf.Blocks[i].Position.X == posX && sf.Blocks[i].Position.Y == posY && sf.Blocks[i].Position.Z == posZ {
					sf.Blocks[i].ID = te.Id
					sf.Blocks[i].Type = "block_with_tile_entity"
					sf.Blocks[i].NBT = nbtData
					found = true
					break
				}
			}
			if !found {
				sf.Blocks = append(sf.Blocks, format.StandardBlock{
					Position: format.StandardBlockPosition{X: posX, Y: posY, Z: posZ},
					ID:       te.Id,
					Type:     "tile_entity",
					NBT:      nbtData,
				})
			}
		}

		// 4. Entities
		for _, en := range region.Entities {
			if en.ID == "" || len(en.Pos) < 3 {
				continue
			}
			blk := format.StandardBlock{
				Type:     "entity",
				ID:       en.ID,
				Position: format.StandardBlockPosition{X: en.Pos[0], Y: en.Pos[1], Z: en.Pos[2]},
			}
			if len(en.Rotation) >= 2 {
				blk.Rotation = format.StandardRotation{Yaw: float64(en.Rotation[0]), Pitch: float64(en.Rotation[1])}
			}
			if len(en.Motion) >= 3 {
				blk.Motion = format.StandardMotion{X: en.Motion[0], Y: en.Motion[1], Z: en.Motion[2]}
			}
			sf.Blocks = append(sf.Blocks, blk)
		}
	}

	if foundAny {
		sf.Position.X = int(minX)
		sf.Position.Y = int(minY)
		sf.Position.Z = int(minZ)
		sf.Size.X = int(maxX - minX)
		sf.Size.Y = int(maxY - minY)
		sf.Size.Z = int(maxZ - minZ)
	}

	return sf, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func canonicalBlockKey(name string, props map[string]string) string {
	if len(props) == 0 {
		return name
	}
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, props[k]))
	}
	return name + "[" + strings.Join(parts, ",") + "]"
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
