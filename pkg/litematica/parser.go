package litematica

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/Tnze/go-mc/nbt"
)

// ParseFile reads and decodes a .litematic file directly into typed structs.
// This avoids the JSON round-trip that was destroying NBT type information
// (int64→float64 precision loss, byte arrays becoming float arrays, etc.).
func ParseFile(filePath string) (*LitematicaNBT, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return Parse(data)
}

// Parse decodes raw (gzip-compressed) NBT bytes directly into a LitematicaNBT.
//
// go-mc/nbt decodes TAG_Long_Array directly into []int64, preserving full
// 64-bit precision. The old JSON approach converted these to []float64,
// losing precision for large values (the root cause of block corruption).
func Parse(data []byte) (*LitematicaNBT, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer r.Close()

	var lit LitematicaNBT
	if _, err = nbt.NewDecoder(r).Decode(&lit); err != nil {
		return nil, fmt.Errorf("failed to decode NBT into LitematicaNBT: %w", err)
	}

	// Post-processing: parse raw tile entity NBT to preserve extra fields.
	// This is done after the initial decode because go-mc doesn't support
	// "catch-all" fields natively. We re-decode the tile entities from raw NBT
	// into a generic map to capture ALL fields including custom mod data.
	if err := parseTileEntitiesRaw(data, &lit); err != nil {
		// Non-fatal: we still have position and id from typed decode
		fmt.Printf("Warning: tile entity raw parse: %v\n", err)
	}

	return &lit, nil
}

// parseTileEntitiesRaw re-reads the raw NBT to extract full tile entity data
// as generic maps, preserving all NBT types.
func parseTileEntitiesRaw(data []byte, lit *LitematicaNBT) error {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer r.Close()

	// Decode as generic interface to get raw map structure
	var raw interface{}
	if _, err = nbt.NewDecoder(r).Decode(&raw); err != nil {
		return err
	}

	rootMap, ok := raw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("root is not a map")
	}

	regionsRaw, ok := rootMap["Regions"].(map[string]interface{})
	if !ok {
		return nil // no regions
	}

	for regionName, regionNBT := range lit.Regions {
		regionRaw, ok := regionsRaw[regionName].(map[string]interface{})
		if !ok {
			continue
		}

		tileEntitiesRaw, ok := regionRaw["TileEntities"].([]interface{})
		if !ok {
			continue
		}

		for i := range regionNBT.TileEntities {
			if i >= len(tileEntitiesRaw) {
				break
			}
			if teMap, ok := tileEntitiesRaw[i].(map[string]interface{}); ok {
				// Copy all fields except the ones we already have typed
				extra := make(map[string]interface{})
				for k, v := range teMap {
					if k == "x" || k == "y" || k == "z" || k == "Id" || k == "id" {
						continue
					}
					extra[k] = v
				}
				regionNBT.TileEntities[i].Extra = extra

				// Some formats use "id" (lowercase) instead of "Id"
				if regionNBT.TileEntities[i].Id == "" {
					if id, ok := teMap["id"].(string); ok {
						regionNBT.TileEntities[i].Id = id
					}
				}
			}
		}

		// Write back modified region
		lit.Regions[regionName] = regionNBT
	}

	return nil
}
