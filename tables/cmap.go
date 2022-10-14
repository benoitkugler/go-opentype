package tables

// Cmap is the Character to Glyph Index Mapping table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap
type Cmap struct {
	version   uint16           // Table version number (0).
	numTables uint16           // Number of encoding tables that follow.
	records   []encodingRecord `len:"numTables"`
	rawData   []byte           `len:"__startToEnd"`
}

type encodingRecord struct {
	platformID     uint16 // Platform ID.
	encodingID     uint16 // Platform-specific encoding ID.
	subtableOffset uint32 // Byte offset from beginning of table to the subtable for this encoding.
}

type cmapSubtable0 struct {
	format       uint16 // Format number is set to 0.
	length       uint16 // This is the length in bytes of the subtable.
	language     uint16
	glyphIdArray [256]uint8 //	An array that maps character codes to glyph index values.
}

type cmapSubtable2 struct {
	format  uint16 // Format number is set to 2.
	rawData []byte `len:"__toEnd"`
}

type cmapSubtable4 struct {
	format         uint16 // Format number is set to 0.
	length         uint16 // This is the length in bytes of the subtable.
	language       uint16
	segCountX2     uint16 // 2 × segCount.
	searchRange    uint16 // Maximum power of 2 less than or equal to segCount, times 2 ((2**floor(log2(segCount))) * 2, where “**” is an exponentiation operator)
	entrySelector  uint16 // Log2 of the maximum power of 2 less than or equal to numTables (log2(searchRange/2), which is equal to floor(log2(segCount)))
	rangeShift     uint16 // segCount times 2, minus searchRange ((segCount * 2) - searchRange)
	endCode        []byte `len:"segCountX2"` // actually [segCount]uint16 End characterCode for each segment, last=0xFFFF.
	reservedPad    uint16 // Set to 0.
	startCode      []byte `len:"segCountX2"` // actually [segCount]uint16 Start character code for each segment.
	idDelta        []byte `len:"segCountX2"` // actually [segCount]int16 Delta for all character codes in segment.
	idRangeOffsets []byte `len:"segCountX2"` // actually [segCount]uint16 Offsets into glyphIdArray or 0
	rawData        []byte `len:"__toEnd"`    // glyphIdArray : uint16[] glyph index array (arbitrary length)
}

type cmapSubtable6 struct {
	format       uint16 // Format number is set to 0.
	length       uint16 // This is the length in bytes of the subtable.
	language     uint16
	firstCode    uint16    // First character code of subrange.
	entryCount   uint16    // Number of character codes in subrange.
	glyphIdArray []glyphID `len:"entryCount"` // Array of glyph index values for character codes in the range.
}

type cmapSubtable10 struct {
	format        uint16 //	Subtable format; set to 10.
	reserved      uint16 //	Reserved; set to 0
	length        uint32 //	Byte length of this subtable (including the header)
	language      uint32
	startCharCode uint32    // First character code covered
	numChars      uint32    // Number of character codes covered
	glyphIdArray  []glyphID `len:"numChars"` // Array of glyph indices for the character codes covered
}

type cmapSubtable12 struct {
	format    uint16               //	Subtable format; set to 12.
	reserved  uint16               //	Reserved; set to 0
	length    uint32               //	Byte length of this subtable (including the header)
	language  uint32               //	For requirements on use of the language field, see “Use of the language field in 'cmap' subtables” in this document.
	numGroups uint32               //	Number of groupings which follow
	groups    []sequentialMapGroup `len:"numGroups"` // Array of SequentialMapGroup records.
}

type sequentialMapGroup struct {
	startCharCode uint32 //	First character code in this group
	endCharCode   uint32 //	Last character code in this group
	startGlyphID  uint32 //	Glyph index corresponding to the starting character code
}

type cmapSubtable13 cmapSubtable12

type cmapSubtable14 struct {
	format                uint16              //	Subtable format. Set to 14.
	length                uint32              //	Byte length of this subtable (including this header)
	numVarSelectorRecords uint32              //	Number of variation Selector Records
	varSelector           []variationSelector `len:"numVarSelectorRecords"` // [numVarSelectorRecords]	Array of VariationSelector records.
}

type variationSelector struct {
	varSelector         [3]byte //	uint24 Variation selector
	defaultUVSOffset    uint32  //	Offset from the start of the format 14 subtable to Default UVS Table. May be 0.
	nonDefaultUVSOffset uint32  //	Offset from the start of the format 14 subtable to Non-Default UVS Table. May be 0.
	rawData             []byte  `len:"__startToEnd"`
}
