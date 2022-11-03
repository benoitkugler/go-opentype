package tables

import "github.com/benoitkugler/go-opentype/opentype"

//go:generate ../../binarygen/cmd/generator . _src.go

type glyphID uint16

// NameID is the ID for entries in the font table.
type NameID uint16

type Tag = opentype.Tag

// Float1616 is a float32, represented in
// fixed 16.16 format in font files.
type Float1616 = float32

func Float1616FromUint(v uint32) Float1616 {
	// value are actually signed integers
	return Float1616(int32(v)) / (1 << 16)
}

func Float1616ToUint(f Float1616) uint32 {
	return uint32(int32(f * (1 << 16)))
}

// Number of seconds since 12:00 midnight that started January 1st 1904 in GMT/UTC time zone.
type longdatetime = uint64

// PlatformID represents the platform id for entries in the name table.
type PlatformID uint16

// EncodingID represents the platform specific id for entries in the name table.
// The most common values are provided as constants.
type EncodingID uint16

// LanguageID represents the language used by an entry in the name table,
// the three most common values are provided as constants.
type LanguageID uint16

// same as binary.BigEndian.Uint32, but for 24 bit uint
func parseUint24(b []byte) rune {
	_ = b[2] // BCE
	return rune(b[0])<<16 | rune(b[1])<<8 | rune(b[2])
}
