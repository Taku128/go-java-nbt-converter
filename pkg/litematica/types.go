// Package litematica provides types and parsers for the Litematica (.litematic) format.
package litematica

// LitematicaNBT is the top-level structure of a .litematic file.
// Fields are tagged with nbt:"..." for direct NBT decoding (no JSON round-trip).
type LitematicaNBT struct {
	Metadata             MetadataNBT          `nbt:"Metadata"`
	MinecraftDataVersion int32                `nbt:"MinecraftDataVersion"`
	Regions              map[string]RegionNBT `nbt:"Regions"`
	Version              int32                `nbt:"Version"`
}

// MetadataNBT holds the schematic's authoring information.
type MetadataNBT struct {
	Author        string       `nbt:"Author"`
	Description   string       `nbt:"Description"`
	EnclosingSize CoordinateNBT `nbt:"EnclosingSize"`
	Name          string       `nbt:"Name"`
	RegionCount   int32        `nbt:"RegionCount"`
	TimeCreated   int64        `nbt:"TimeCreated"`
	TimeModified  int64        `nbt:"TimeModified"`
	TotalBlocks   int32        `nbt:"TotalBlocks"`
	TotalVolume   int32        `nbt:"TotalVolume"`
}

// CoordinateNBT represents a 3D coordinate as stored in NBT.
type CoordinateNBT struct {
	X int32 `nbt:"x"`
	Y int32 `nbt:"y"`
	Z int32 `nbt:"z"`
}

// RegionNBT represents a single region in a Litematica schematic.
// BlockStates is decoded directly as []int64 (TAG_Long_Array) to avoid
// the precision loss that occurs when round-tripping through JSON.
type RegionNBT struct {
	BlockStatePalette []BlockStatePaletteNBT `nbt:"BlockStatePalette"`
	BlockStates       []int64                `nbt:"BlockStates"`
	Entities          []EntityNBT            `nbt:"Entities"`
	PendingBlockTicks []interface{}          `nbt:"PendingBlockTicks"`
	PendingFluidTicks []interface{}          `nbt:"PendingFluidTicks"`
	Position          CoordinateNBT          `nbt:"Position"`
	Size              CoordinateNBT          `nbt:"Size"`
	TileEntities      []TileEntityNBT        `nbt:"TileEntities"`
}

// BlockStatePaletteNBT represents a block state in the region palette.
type BlockStatePaletteNBT struct {
	Name       string            `nbt:"Name"`
	Properties map[string]string `nbt:"Properties"`
}

// TileEntityNBT represents a tile entity with raw NBT data preserved.
// Only positional fields and id are strongly-typed; the rest is kept as
// a raw map to preserve all NBT type information.
type TileEntityNBT struct {
	Id   string `nbt:"Id"`
	X    int32  `nbt:"x"`
	Y    int32  `nbt:"y"`
	Z    int32  `nbt:"z"`
	// Extra holds all other NBT fields not captured by the typed fields above.
	// This is populated manually after initial decode.
	Extra map[string]interface{} `nbt:"-"`
}

// EntityNBT represents an entity in a Litematica region.
type EntityNBT struct {
	Pos      []float64 `nbt:"Pos"`
	Rotation []float32 `nbt:"Rotation"`
	Motion   []float64 `nbt:"Motion"`
	ID       string    `nbt:"id"`
}
