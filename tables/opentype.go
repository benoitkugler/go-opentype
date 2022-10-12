package tables

type glyphID uint16

// NameID is the ID for entries in the font table.
type NameID uint16

// Tag represents an open-type name.
// These are technically uint32's, but are usually
// displayed in ASCII as they are all acronyms.
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6.html#Overview
type Tag uint32

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
