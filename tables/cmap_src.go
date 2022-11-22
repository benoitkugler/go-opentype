package tables

import (
	"encoding/binary"
	"fmt"
)

// Cmap is the Character to Glyph Index Mapping table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap
type Cmap struct {
	version   uint16           // Table version number (0).
	numTables uint16           // Number of encoding tables that follow.
	records   []encodingRecord `arrayCount:"ComputedField-numTables"`
	subtables cmapSubtables    `isOpaque:"" subsliceStart:"AtStart"`
}

type cmapSubtables []cmapSubtable

func (cm *Cmap) customParseSubtables(src []byte) (int, error) {
	var err error
	cm.subtables = make(cmapSubtables, len(cm.records))
	for i, rec := range cm.records {
		// add the format length (2)
		if L := len(src); L < int(rec.subtableOffset)+2 {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.subtableOffset, L)
		}
		subtableData := src[rec.subtableOffset:]
		format := binary.BigEndian.Uint16(subtableData)
		switch format {
		case 0:
			cm.subtables[i], _, err = ParseCmapSubtable0(subtableData)
		case 2:
			cm.subtables[i], _, err = ParseCmapSubtable2(subtableData)
		case 4:
			cm.subtables[i], _, err = ParseCmapSubtable4(subtableData)
		case 6:
			cm.subtables[i], _, err = ParseCmapSubtable6(subtableData)
		case 10:
			cm.subtables[i], _, err = ParseCmapSubtable10(subtableData)
		case 12:
			cm.subtables[i], _, err = ParseCmapSubtable12(subtableData)
		case 13:
			cm.subtables[i], _, err = ParseCmapSubtable13(subtableData)
		case 14:
			cm.subtables[i], _, err = ParseCmapSubtable14(subtableData)
		default:
			return 0, fmt.Errorf("unsupported cmap subtable format: %d", format)
		}
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type encodingRecord struct {
	platformID     uint16 // Platform ID.
	encodingID     uint16 // Platform-specific encoding ID.
	subtableOffset uint32 // Byte offset from beginning of table to the subtable for this encoding.
}

type cmapSubtable interface {
	isCmapSubtable()
}

func (CmapSubtable0) isCmapSubtable()  {}
func (CmapSubtable2) isCmapSubtable()  {}
func (CmapSubtable4) isCmapSubtable()  {}
func (CmapSubtable6) isCmapSubtable()  {}
func (CmapSubtable10) isCmapSubtable() {}
func (CmapSubtable12) isCmapSubtable() {}
func (CmapSubtable13) isCmapSubtable() {}
func (CmapSubtable14) isCmapSubtable() {}

type CmapSubtable0 struct {
	format       uint16 // Format number is set to 0.
	length       uint16 // This is the length in bytes of the subtable.
	language     uint16
	glyphIdArray [256]uint8 //	An array that maps character codes to glyph index values.
}

type CmapSubtable2 struct {
	format  uint16 // Format number is set to 2.
	rawData []byte `arrayCount:"ToEnd"`
}

type CmapSubtable4 struct {
	format         uint16 // Format number is set to 0.
	length         uint16 // This is the length in bytes of the subtable.
	language       uint16
	segCountX2     uint16   // 2 × segCount.
	searchRange    uint16   // Maximum power of 2 less than or equal to segCount, times 2 ((2**floor(log2(segCount))) * 2, where “**” is an exponentiation operator)
	entrySelector  uint16   // Log2 of the maximum power of 2 less than or equal to numTables (log2(searchRange/2), which is equal to floor(log2(segCount)))
	rangeShift     uint16   // segCount times 2, minus searchRange ((segCount * 2) - searchRange)
	endCode        []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]uint16 End characterCode for each segment, last=0xFFFF.
	reservedPad    uint16   // Set to 0.
	startCode      []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]uint16 Start character code for each segment.
	idDelta        []int16  `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]int16 Delta for all character codes in segment.
	idRangeOffsets []uint16 `arrayCount:"ComputedField-segCountX2 / 2"` // [segCount]uint16 Offsets into glyphIdArray or 0
	glyphIDArray   []byte   `arrayCount:"ToEnd"`                        // glyphIdArray : uint16[] glyph index array (arbitrary length)
}

type CmapSubtable6 struct {
	format       uint16 // Format number is set to 0.
	length       uint16 // This is the length in bytes of the subtable.
	language     uint16
	firstCode    uint16    // First character code of subrange.
	entryCount   uint16    // Number of character codes in subrange.
	glyphIdArray []GlyphID `arrayCount:"ComputedField-entryCount"` // Array of glyph index values for character codes in the range.
}

type CmapSubtable10 struct {
	format        uint16 //	Subtable format; set to 10.
	reserved      uint16 //	Reserved; set to 0
	length        uint32 //	Byte length of this subtable (including the header)
	language      uint32
	startCharCode uint32    // First character code covered
	numChars      uint32    // Number of character codes covered
	glyphIdArray  []GlyphID `arrayCount:"ComputedField-numChars"` // Array of glyph indices for the character codes covered
}

type CmapSubtable12 struct {
	format    uint16               //	Subtable format; set to 12.
	reserved  uint16               //	Reserved; set to 0
	length    uint32               //	Byte length of this subtable (including the header)
	language  uint32               //	For requirements on use of the language field, see “Use of the language field in 'cmap' subtables” in this document.
	numGroups uint32               //	Number of groupings which follow
	groups    []sequentialMapGroup `arrayCount:"ComputedField-numGroups"` // Array of SequentialMapGroup records.
}

type sequentialMapGroup struct {
	startCharCode uint32 //	First character code in this group
	endCharCode   uint32 //	Last character code in this group
	startGlyphID  uint32 //	Glyph index corresponding to the starting character code
}

type CmapSubtable13 CmapSubtable12

type CmapSubtable14 struct {
	format                uint16              // Subtable format. Set to 14.
	length                uint32              // Byte length of this subtable (including this header)
	numVarSelectorRecords uint32              // Number of variation Selector Records
	varSelectors          []variationSelector `arrayCount:"ComputedField-numVarSelectorRecords"` // [numVarSelectorRecords]	Array of VariationSelector records.
	rawData               []byte              `arrayCount:"ToEnd" subsliceStart:"AtStart"`       // used to interpret [varSelectors]
}

type variationSelector struct {
	varSelector         [3]byte // uint24 Variation selector
	defaultUVSOffset    uint32  // Offset from the start of the format 14 subtable to Default UVS Table. May be 0.
	nonDefaultUVSOffset uint32  // Offset from the start of the format 14 subtable to Non-Default UVS Table. May be 0.
}

// DefaultUVSTable is used in Cmap format 14
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap#default-uvs-table
type DefaultUVSTable struct {
	ranges []unicodeRange `arrayCount:"FirstUint32"`
}

type unicodeRange struct {
	startUnicodeValue [3]byte // uint24 First value in this range
	additionalCount   uint8   // Number of additional values in this range
}

// UVSMappingTable is used in Cmap format 14
// See https://learn.microsoft.com/en-us/typography/opentype/spec/cmap#non-default-uvs-table
type UVSMappingTable struct {
	ranges []uvsMappingRecord `arrayCount:"FirstUint32"`
}

type uvsMappingRecord struct {
	unicodeValue [3]byte // uint24 Base Unicode value of the UVS
	glyphID      uint16  //	Glyph ID of the UVS
}
