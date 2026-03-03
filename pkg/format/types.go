// Package format defines the shared intermediate representation used by all
// format-specific parsers and the final NBT encoder.
//
// Every input format (Litematica, WorldEdit .schem, Classic .schematic, Java Structure .nbt)
// is first converted into StandardFormat.  The encoder then takes a StandardFormat and
// produces a gzip-compressed Java Edition Structure NBT byte slice.
package format

// StandardFormat represents a unified structure that can hold data from
// different Minecraft schematic formats (Litematica, WorldEdit, etc.)
type StandardFormat struct {
	Metadata       StandardMetadata          `json:"metadata"`
	DataVersion    int                       `json:"dataVersion"`
	Version        int                       `json:"version"`
	Size           StandardSize              `json:"size"`
	Position       StandardPosition          `json:"position"`
	Blocks         []StandardBlock           `json:"blocks"`
	Palette        map[int]StandardPalette   `json:"palette"`
	OriginalFormat string                    `json:"originalFormat"`
}

// StandardMetadata contains authoring information about the schematic.
type StandardMetadata struct {
	Name             string `json:"name"`
	Author           string `json:"author"`
	Description      string `json:"description"`
	TimeCreated      int64  `json:"timeCreated"`
	TimeModified     int64  `json:"timeModified"`
	TotalBlocks      int    `json:"totalBlocks"`
	TotalVolume      int    `json:"totalVolume"`
	PreviewImageData []int  `json:"previewImageData,omitempty"`
}

// StandardSize represents the 3D dimensions of the schematic.
type StandardSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// StandardPosition represents a 3D origin position.
type StandardPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// StandardBlock represents a placed block, entity, or tile entity.
type StandardBlock struct {
	Type     string                `json:"type,omitempty"`
	ID       string                `json:"id,omitempty"`
	Position StandardBlockPosition `json:"position"`
	Rotation StandardRotation      `json:"rotation,omitempty"`
	Motion   StandardMotion        `json:"motion,omitempty"`
	State    int                   `json:"state,omitempty"`
	NBT      interface{}           `json:"nbt,omitempty"`
}

// StandardBlockPosition is a float64 position (entities may have sub-block positions).
type StandardBlockPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// StandardRotation represents the rotation of an entity.
type StandardRotation struct {
	Yaw   float64 `json:"yaw,omitempty"`
	Pitch float64 `json:"pitch,omitempty"`
}

// StandardMotion represents the velocity of an entity.
type StandardMotion struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
	Z float64 `json:"z,omitempty"`
}

// StandardPalette represents a block type in the palette.
type StandardPalette struct {
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties,omitempty"`
}
