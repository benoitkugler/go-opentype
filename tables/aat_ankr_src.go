package tables

// Ankr is the anchor point table
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6ankr.html
type Ankr struct {
	version uint16 // Version number (set to zero)
	flags   uint16 // Flags (currently unused; set to zero)
	// Offset to the table's lookup table; currently this is always 0x0000000C
	// The lookup table returns uint16 offset from the beginning of the glyph data table, not indices.
	LookupTable AATLookup `offsetSize:"Offset32"`
	// Offset to the glyph data table
	GlyphDataTable []byte `offsetSize:"Offset32" arrayCount:"ToEnd"`
}

// AnrkAnchor is a point within the coordinate space of a given glyph
// independent of the control points used to render the glyph
type AnrkAnchor struct {
	X, Y int16
}
