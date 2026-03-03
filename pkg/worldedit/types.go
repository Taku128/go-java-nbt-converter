// Package worldedit provides types and parsers for WorldEdit .schem files
// (Sponge Schematic v2 / v3).
package worldedit

// WorldEditNBT is the top-level structure of a .schem file.
// Fields use nbt tags for direct NBT decoding.
type WorldEditNBT struct {
	BlockData     []int8           `nbt:"BlockData"`
	BlockEntities []interface{}    `nbt:"BlockEntities"`
	DataVersion   int32            `nbt:"DataVersion"`
	Height        int16            `nbt:"Height"`
	Length        int16            `nbt:"Length"`
	Metadata      MetadataNBT      `nbt:"Metadata"`
	Offset        []int32          `nbt:"Offset"`
	Palette       map[string]int32 `nbt:"Palette"`
	PaletteMax    int32            `nbt:"PaletteMax"`
	Version       int32            `nbt:"Version"`
	Width         int16            `nbt:"Width"`
}

// MetadataNBT holds WorldEdit-specific metadata.
type MetadataNBT struct {
	WEOffsetX int32 `nbt:"WEOffsetX"`
	WEOffsetY int32 `nbt:"WEOffsetY"`
	WEOffsetZ int32 `nbt:"WEOffsetZ"`
}
