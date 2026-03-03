// Package buildnbt converts a StandardFormat into a gzip-compressed
// Java Edition Structure NBT byte slice.
package buildnbt

import (
	"io"
)

// ---------------------------------------------------------------------------
// ListTag encodes an []int32 as TAG_List(TAG_Int) instead of TAG_Int_Array.
// Java Edition's structure NBT strictly requires TAG_List for pos / size.
// ---------------------------------------------------------------------------

// ListTag is a custom type that serialises as TAG_List of TAG_Int in NBT.
type ListTag []int32

// TagType satisfies go-mc's nbt.Marshaler interface hint.
func (l ListTag) TagType() byte { return 9 } // TAG_List

// MarshalNBT writes the list header (element-type + length) followed by
// each int32 in big-endian order.
func (l ListTag) MarshalNBT(w io.Writer) error {
	var buf [5]byte
	buf[0] = 3 // TAG_Int
	buf[1] = byte(len(l) >> 24)
	buf[2] = byte(len(l) >> 16)
	buf[3] = byte(len(l) >> 8)
	buf[4] = byte(len(l))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	for _, v := range l {
		var vBuf [4]byte
		vBuf[0] = byte(v >> 24)
		vBuf[1] = byte(v >> 16)
		vBuf[2] = byte(v >> 8)
		vBuf[3] = byte(v)
		if _, err := w.Write(vBuf[:]); err != nil {
			return err
		}
	}
	return nil
}

// PaletteEntry is a single block-state in the Java structure palette.
type PaletteEntry struct {
	Name       string            `nbt:"Name"`
	Properties map[string]string `nbt:"Properties,omitempty"`
}

// JavaBlock is a single placed block in Java structure format.
type JavaBlock struct {
	Pos   ListTag     `nbt:"pos"`
	State int32       `nbt:"state"`
	Nbt   interface{} `nbt:"nbt,omitempty"`
}

// JavaNBT is the top-level Java Edition structure NBT compound.
type JavaNBT struct {
	DataVersion int32          `nbt:"DataVersion"`
	Size        ListTag        `nbt:"size"`
	Palette     []PaletteEntry `nbt:"palette"`
	Blocks      []JavaBlock    `nbt:"blocks"`
	Entities    []interface{}  `nbt:"entities"`
}
