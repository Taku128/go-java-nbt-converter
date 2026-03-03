# go-java-nbt-converter

A high-performance Go library for converting various Minecraft Java Edition schematic formats into the standardized Java Edition Structure NBT format. This library is designed to be used in high-load environments like cloud functions (AWS Lambda).

## Features
- **Accurate NBT Handling**: Uses direct NBT decoding (no JSON round-tripping) to preserve all 64-bit precision and complex NBT types.
- **Support for Multiple Formats**:
  - `.litematic` (Litematica)
  - `.schem` (Sponge/WorldEdit Schematic v2 and v3)
  - `.nbt` (Java Edition Structure NBT)
- **Automatic Block Registry Updates**: Integrated with GitHub Actions to auto-fetch the latest block data from Mojang.
- **Fast and Lightweight**: Minimal dependencies, optimized for speed.

## Installation
```bash
go get github.com/ntaku256/go-java-nbt-converter
```

## Quick Start
```go
import (
	javanbt "github.com/ntaku256/go-java-nbt-converter"
)

func main() {
    // Converts any supported format to a gzip-compressed Java NBT byte slice.
    nbtBytes, err := javanbt.ConvertAny("path/to/my_schematic.litematic")
    if err != nil {
        log.Fatal(err)
    }
    // ... use nbtBytes ...
}
```

## Supported Formats Details
- **Litematica (.litematic)**: Full multi-region support, correct block packing/unpacking, and preservation of Tile Entity data.
- **WorldEdit/Sponge (.schem)**: Support for V2 and V3, VarInt block data, and entities.
- **Java Structure (.nbt)**: Standard output format, also supported as an input for normalization.

## Contributing
We welcome contributions! Please feel free to submit a Pull Request.

## License
MIT
