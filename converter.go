// Package javanbt provides a high-performance converter to translate Java Edition
// schematic formats (.litematic, .schem, .nbt) into a standardized Java Edition
// Structure NBT format.
//
// Usage:
//
//	nbtBytes, err := javanbt.ConvertAny("path/to/file.litematic")
//	// nbtBytes is a gzip-compressed Java Edition Structure NBT byte slice.
package javanbt

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ntaku256/go-java-nbt-converter/pkg/buildnbt"
	"github.com/ntaku256/go-java-nbt-converter/pkg/litematica"
	"github.com/ntaku256/go-java-nbt-converter/pkg/structure"
	"github.com/ntaku256/go-java-nbt-converter/pkg/worldedit"
)

// ConvertLitematica converts a .litematic file to Java Structure NBT bytes.
func ConvertLitematica(ctx context.Context, filePath string) ([]byte, error) {
	lit, err := litematica.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("litematica parse: %w", err)
	}
	sf, err := litematica.Convert(lit)
	if err != nil {
		return nil, fmt.Errorf("litematica convert: %w", err)
	}
	return buildnbt.Encode(sf)
}

// ConvertSchem converts a .schem (Sponge Schematic v2/v3) file to Java Structure NBT bytes.
func ConvertSchem(ctx context.Context, filePath string) ([]byte, error) {
	we, err := worldedit.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("worldedit parse: %w", err)
	}
	sf, err := worldedit.Convert(we)
	if err != nil {
		return nil, fmt.Errorf("worldedit convert: %w", err)
	}
	return buildnbt.Encode(sf)
}

// ConvertStructureNBT converts a .nbt (Java Structure) file to standardized
// Java Structure NBT bytes.
func ConvertStructureNBT(ctx context.Context, filePath string) ([]byte, error) {
	s, err := structure.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("structure parse: %w", err)
	}
	sf, err := structure.Convert(s)
	if err != nil {
		return nil, fmt.Errorf("structure convert: %w", err)
	}
	return buildnbt.Encode(sf)
}

// ConvertAny auto-detects the format from the file extension and converts
// accordingly.  Supported extensions: .litematic, .schem, .nbt
func ConvertAny(ctx context.Context, filePath string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".litematic":
		return ConvertLitematica(ctx, filePath)
	case ".schem":
		return ConvertSchem(ctx, filePath)
	case ".nbt":
		return ConvertStructureNBT(ctx, filePath)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s (supported: .litematic, .schem, .nbt)", ext)
	}
}
