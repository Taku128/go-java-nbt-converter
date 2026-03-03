package registry

import (
	"strings"
)

// IsValidBlock checks if a block name is known to this registry.
func IsValidBlock(name string) bool {
	if !strings.Contains(name, ":") {
		name = "minecraft:" + name
	}
	_, ok := BlockProperties[name]
	if ok {
		return true
	}
	// Also check the AllBlocks list for blocks without properties
	for _, b := range AllBlocks {
		if b == name {
			return true
		}
	}
	return false
}

// GetDefaultProperties returns the default properties for a given block.
func GetDefaultProperties(name string) map[string]string {
	if !strings.Contains(name, ":") {
		name = "minecraft:" + name
	}
	props, ok := BlockProperties[name]
	if !ok {
		return nil
	}
	defaults := make(map[string]string)
	for k, v := range props {
		defaults[k] = v.Default
	}
	return defaults
}

// NormalizeTagName ensures a block name has the "minecraft:" prefix if missing.
func NormalizeTagName(name string) string {
	if !strings.Contains(name, ":") {
		return "minecraft:" + name
	}
	return name
}
