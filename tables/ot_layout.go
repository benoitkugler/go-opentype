package tables

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
)

// The following are types shared by GSUB and GPOS tables

// Coverage specifies all the glyphs affected by a substitution or
// positioning operation described in a subtable.
// Conceptually is it a []GlyphIndex, but it may be implemented for efficiently.
// See https://learn.microsoft.com/typography/opentype/spec/chapter2#lookup-table
type Coverage interface {
	isCoverage()
}

func (Coverage1) isCoverage() {}
func (Coverage2) isCoverage() {}

type Coverage1 struct {
	format uint16    `unionTag:"1"`
	Glyphs []GlyphID `arrayCount:"FirstUint16"`
}

type Coverage2 struct {
	format uint16        `unionTag:"2"`
	Ranges []RangeRecord `arrayCount:"FirstUint16"`
}

type RangeRecord struct {
	StartGlyphID       GlyphID // First glyph ID in the range
	EndGlyphID         GlyphID // Last glyph ID in the range
	StartCoverageIndex uint16  // Coverage Index of first glyph ID in range
}

// ClassDef stores a value for a set of GlyphIDs.
// Conceptually it is a map[GlyphID]uint16, but it may
// be implemented more efficiently.
type ClassDef interface {
	isClassDef()
}

func (ClassDef1) isClassDef() {}
func (ClassDef2) isClassDef() {}

type ClassDef1 struct {
	format          uint16   `unionTag:"1"`
	StartGlyphID    GlyphID  // First glyph ID of the classValueArray
	ClassValueArray []uint16 `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type ClassDef2 struct {
	format            uint16             `unionTag:"2"`
	ClassRangeRecords []ClassRangeRecord `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type ClassRangeRecord struct {
	StartGlyphID GlyphID // First glyph ID in the range
	EndGlyphID   GlyphID // Last glyph ID in the range
	Class        uint16  // Applied to all glyphs in the range
}

// Lookups

type SequenceLookupRecord struct {
	SequenceIndex   uint16 // Index (zero-based) into the input glyph sequence
	LookupListIndex uint16 // Index (zero-based) into the LookupList
}

type SequenceContextFormat1 struct {
	format     uint16            `unionTag:"1"`
	Coverage   Coverage          `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of SequenceContextFormat1 table
	SeqRuleSet []SequenceRuleSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[seqRuleSetCount]	Array of offsets to SequenceRuleSet tables, from beginning of SequenceContextFormat1 table (offsets may be NULL)
}

type SequenceRuleSet struct {
	SeqRule []SequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type SequenceRule struct {
	glyphCount       uint16                 // Number of glyphs in the input glyph sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	InputSequence    []GlyphID              `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Array of input glyph IDs—starting with the second glyph
	SeqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of Sequence lookup records
}

type SequenceContextFormat2 struct {
	format          uint16                 `unionTag:"2"`
	Coverage        Coverage               `offsetSize:"Offset16"`                            //	Offset to Coverage table, from beginning of SequenceContextFormat2 table
	ClassDef        ClassDef               `offsetSize:"Offset16"`                            //	Offset to ClassDef table, from beginning of SequenceContextFormat2 table
	ClassSeqRuleSet []ClassSequenceRuleSet `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[classSeqRuleSetCount]	Array of offsets to ClassSequenceRuleSet tables, from beginning of SequenceContextFormat2 table (may be NULL)
}

type ClassSequenceRuleSet struct {
	ClassSeqRule []ClassSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //	Array of offsets to ClassSequenceRule tables, from beginning of ClassSequenceRuleSet table
}

type ClassSequenceRule struct {
	glyphCount       uint16                 // Number of glyphs to be matched
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	InputSequence    []uint16               `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Sequence of classes to be matched to the input glyph sequence, beginning with the second glyph position
	SeqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of SequenceLookupRecords
}

type SequenceContextFormat3 struct {
	format           uint16                 `unionTag:"3"`
	glyphCount       uint16                 // Number of glyphs in the input sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	CoverageOffsets  []Coverage             `arrayCount:"ComputedField-glyphCount" offsetsArray:"Offset16"` //[glyphCount]	Array of offsets to Coverage tables, from beginning of SequenceContextFormat3 subtable
	SeqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

type ChainedSequenceContextFormat1 struct {
	format            uint16                   `unionTag:"1"`
	Coverage          Coverage                 `offsetSize:"Offset16"`                             //	Offset to Coverage table, from beginning of ChainSequenceContextFormat1 table
	ChainedSeqRuleSet []ChainedSequenceRuleSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[chainedSeqRuleSetCount]	Array of offsets to ChainedSeqRuleSet tables, from beginning of ChainedSequenceContextFormat1 table (may be NULL)
}

type ChainedSequenceRuleSet struct {
	ChainedSeqRules []ChainedSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type ChainedSequenceRule struct {
	BacktrackSequence []GlyphID              `arrayCount:"FirstUint16"` //[backtrackGlyphCount]	Array of backtrack glyph IDs
	inputGlyphCount   uint16                 //	Number of glyphs in the input sequence
	InputSequence     []GlyphID              `arrayCount:"ComputedField-inputGlyphCount-1"` //[inputGlyphCount - 1]	Array of input glyph IDs—start with second glyph
	LookaheadSequence GlyphID                `arrayCount:"FirstUint16"`                     //[lookaheadGlyphCount]	Array of lookahead glyph IDs
	SeqLookupRecords  []SequenceLookupRecord `arrayCount:"FirstUint16"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

type ChainedSequenceContextFormat2 struct {
	format                 uint16                        `unionTag:"2"`
	Coverage               Coverage                      `offsetSize:"Offset16"`                            // Offset to Coverage table, from beginning of ChainedSequenceContextFormat2 table
	BacktrackClassDef      ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing backtrack sequence context, from beginning of ChainedSequenceContextFormat2 table
	InputClassDef          ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing input sequence context, from beginning of ChainedSequenceContextFormat2 table
	LookaheadClassDef      ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing lookahead sequence context, from beginning of ChainedSequenceContextFormat2 table
	ChainedClassSeqRuleSet []ChainedClassSequenceRuleSet `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[chainedClassSeqRuleSetCount]	Array of offsets to ChainedClassSequenceRuleSet tables, from beginning of ChainedSequenceContextFormat2 table (may be NULL)
}

type ChainedClassSequenceRuleSet struct {
	ChainedClassSeqRules []ChainedClassSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to ChainedClassSequenceRule tables, from beginning of ChainedClassSequenceRuleSet
}

type ChainedClassSequenceRule struct {
	BacktrackSequence []uint16               `arrayCount:"FirstUint16"` //[backtrackGlyphCount]	Array of backtrack-sequence classes
	inputGlyphCount   uint16                 //	Total number of glyphs in the input sequence
	InputSequence     []uint16               `arrayCount:"ComputedField-inputGlyphCount-1"` //[inputGlyphCount - 1]	Array of input sequence classes, beginning with the second glyph position
	LookaheadGlyph    []uint16               `arrayCount:"FirstUint16"`                     //[lookaheadGlyphCount]	Array of lookahead-sequence classes
	SeqLookupRecords  []SequenceLookupRecord `arrayCount:"FirstUint16"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

type ChainedSequenceContextFormat3 struct {
	format             uint16                 `unionTag:"3"`
	BacktrackCoverages []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[backtrackGlyphCount]	Array of offsets to coverage tables for the backtrack sequence
	InputCoverages     []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[inputGlyphCount]	Array of offsets to coverage tables for the input sequence
	LookaheadCoverages []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[lookaheadGlyphCount]	Array of offsets to coverage tables for the lookahead sequence
	SeqLookupRecords   []SequenceLookupRecord `arrayCount:"FirstUint16"`                         //[seqLookupCount]	Array of SequenceLookupRecords
}

type Extension struct {
	substFormat         uint16   //	Format identifier. Set to 1.
	ExtensionLookupType uint16   //	Lookup type of subtable referenced by extensionOffset (that is, the extension subtable).
	ExtensionOffset     Offset32 //	Offset to the extension subtable, of lookup type extensionLookupType, relative to the start of the ExtensionSubstFormat1 subtable.
	RawData             []byte   `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

// GSUB is the Glyph Substitution (GSUB) table.
// It provides data for substition of glyphs for appropriate rendering of scripts,
// such as cursively-connecting forms in Arabic script,
// or for advanced typographic effects, such as ligatures.
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/gsub
type GSUB Layout

type GSUBLookup interface {
	isGSUBLookup()
}

func (SingleSubs) isGSUBLookup()             {}
func (MultipleSubs) isGSUBLookup()           {}
func (AlternateSubs) isGSUBLookup()          {}
func (LigatureSubs) isGSUBLookup()           {}
func (ContextualSubs) isGSUBLookup()         {}
func (ChainedContextualSubs) isGSUBLookup()  {}
func (ExtensionSubs) isGSUBLookup()          {}
func (ReverseChainSingleSubs) isGSUBLookup() {}

// AsGSUBLookups returns the GSUB lookup subtables
func (lk Lookup) AsGSUBLookups() ([]GSUBLookup, error) {
	var err error
	out := make([]GSUBLookup, len(lk.subtableOffsets))
	for i, offset := range lk.subtableOffsets {
		if L := len(lk.rawData); L < int(offset) {
			return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		switch lk.lookupType {
		case 1: // Single (format 1.1 1.2)	Replace one glyph with one glyph
			out[i], _, err = ParseSingleSubs(lk.rawData[offset:])
		case 2: // Multiple (format 2.1)	Replace one glyph with more than one glyph
			out[i], _, err = ParseMultipleSubs(lk.rawData[offset:])
		case 3: // Alternate (format 3.1)	Replace one glyph with one of many glyphs
			out[i], _, err = ParseAlternateSubs(lk.rawData[offset:])
		case 4: // Ligature (format 4.1)	Replace multiple glyphs with one glyph
			out[i], _, err = ParseLigatureSubs(lk.rawData[offset:])
		case 5: // Context (format 5.1 5.2 5.3)	Replace one or more glyphs in context
			out[i], _, err = ParseContextualSubs(lk.rawData[offset:])
		case 6: // Chaining Context (format 6.1 6.2 6.3)	Replace one or more glyphs in chained context
			out[i], _, err = ParseChainedContextualSubs(lk.rawData[offset:])
		case 7: // Extension Substitution (format 7.1) Extension mechanism for other substitutions
			out[i], _, err = ParseExtensionSubs(lk.rawData[offset:])
		case 8: // Reverse chaining context single (format 8.1)
			out[i], _, err = ParseReverseChainSingleSubs(lk.rawData[offset:])
		default:
			err = fmt.Errorf("invalid GSUB Loopkup type %d", lk.lookupType)
		}
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ------------------------ GPOS common data structures ------------------------

// GPOS is the Glyph Positioning (GPOS) table.
// It provides precise control over glyph placement
// for sophisticated text layout and rendering in each script
// and language system that a font supports.
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/gpos
type GPOS Layout

type GPOSLookup interface {
	isGPOSLookup()
}

func (SinglePos) isGPOSLookup()            {}
func (PairPos) isGPOSLookup()              {}
func (CursivePos) isGPOSLookup()           {}
func (MarkBasePos) isGPOSLookup()          {}
func (MarkLigPos) isGPOSLookup()           {}
func (MarkMarkPos) isGPOSLookup()          {}
func (ContextualPos) isGPOSLookup()        {}
func (ChainedContextualPos) isGPOSLookup() {}
func (ExtensionPos) isGPOSLookup()         {}

// AsGPOSLookups returns the GPOS lookup subtables
func (lk Lookup) AsGPOSLookups() ([]GPOSLookup, error) {
	var err error
	out := make([]GPOSLookup, len(lk.subtableOffsets))
	for i, offset := range lk.subtableOffsets {
		if L := len(lk.rawData); L < int(offset) {
			return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		switch lk.lookupType {
		case 1: // Single adjustment	Adjust position of a single glyph
			out[i], _, err = ParseSinglePos(lk.rawData[offset:])
		case 2: // Pair adjustment	Adjust position of a pair of glyphs
			out[i], _, err = ParsePairPos(lk.rawData[offset:])
		case 3: // Cursive attachment	Attach cursive glyphs
			out[i], _, err = ParseCursivePos(lk.rawData[offset:])
		case 4: // MarkToBase attachment	Attach a combining mark to a base glyph
			out[i], _, err = ParseMarkBasePos(lk.rawData[offset:])
		case 5: // MarkToLigature attachment	Attach a combining mark to a ligature
			out[i], _, err = ParseMarkLigPos(lk.rawData[offset:])
		case 6: // MarkToMark attachment	Attach a combining mark to another mark
			out[i], _, err = ParseMarkMarkPos(lk.rawData[offset:])
		case 7: // Context positioning	Position one or more glyphs in context
			out[i], _, err = ParseContextualPos(lk.rawData[offset:])
		case 8: // Chained Context positioning	Position one or more glyphs in chained context
			out[i], _, err = ParseChainedContextualPos(lk.rawData[offset:])
		case 9: // Extension positioning	Extension mechanism for other positionings
			out[i], _, err = ParseExtensionPos(lk.rawData[offset:])
		default:
			err = fmt.Errorf("invalid GPOS Loopkup type %d", lk.lookupType)
		}
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ValueFormat is a mask indicating which field
// are set in a [ValueRecord].
// It is often shared between many records.
type ValueFormat uint16

// number of fields present
func (f ValueFormat) size() int { return bits.OnesCount16(uint16(f)) }

const (
	XPlacement ValueFormat = 1 << iota // Includes horizontal adjustment for placement
	YPlacement                         // Includes vertical adjustment for placement
	XAdvance                           // Includes horizontal adjustment for advance
	YAdvance                           // Includes vertical adjustment for advance
	XPlaDevice                         // Includes horizontal Device table for placement
	YPlaDevice                         // Includes vertical Device table for placement
	XAdvDevice                         // Includes horizontal Device table for advance
	YAdvDevice                         // Includes vertical Device table for advance

	//  Mask for having any Device table
	Devices = XPlaDevice | YPlaDevice | XAdvDevice | YAdvDevice
)

// ValueRecord has optional fields
type ValueRecord struct {
	XPlacement int16       // Horizontal adjustment for placement, in design units.
	YPlacement int16       // Vertical adjustment for placement, in design units.
	XAdvance   int16       // Horizontal adjustment for advance, in design units — only used for horizontal layout.
	YAdvance   int16       // Vertical adjustment for advance, in design units — only used for vertical layout.
	XPlaDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for horizontal placement, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup subtable) — may be NULL.
	YPlaDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for vertical placement, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup subtable) — may be NULL.
	XAdvDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for horizontal advance, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup subtable) — may be NULL.
	YAdvDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for vertical advance, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup su
}

// [data] must start at the immediate parent table, [offset] indicating
// the start of the record in it.
// Returns [offset] + the number of bytes read from [offset]
// Note that a [format] with value 0, is supported, resulting in a no-op
func parseValueRecord(format ValueFormat, data []byte, offset int) (out ValueRecord, _ int, err error) {
	if L := len(data); L < offset {
		return out, 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
	}

	size := format.size() // number of fields present
	if size == 0 {        // return early
		return out, offset, nil
	}
	// start by parsing the list of values
	values, err := parseUint16s(data[offset:], size)
	if err != nil {
		return out, 0, fmt.Errorf("invalid value record: %s", err)
	}
	// follow the order
	if format&XPlacement != 0 {
		out.XPlacement = int16(values[0])
		values = values[1:]
	}
	if format&YPlacement != 0 {
		out.YPlacement = int16(values[0])
		values = values[1:]
	}
	if format&XAdvance != 0 {
		out.XAdvance = int16(values[0])
		values = values[1:]
	}
	if format&YAdvance != 0 {
		out.YAdvance = int16(values[0])
		values = values[1:]
	}
	if format&XPlaDevice != 0 {
		if devOffset := values[0]; devOffset != 0 {
			out.XPlaDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		values = values[1:]
	}
	if format&YPlaDevice != 0 {
		if devOffset := values[0]; devOffset != 0 {
			out.YPlaDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		values = values[1:]
	}
	if format&XAdvDevice != 0 {
		if devOffset := values[0]; devOffset != 0 {
			out.XAdvDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		values = values[1:]
	}
	if format&YAdvDevice != 0 {
		if devOffset := values[0]; devOffset != 0 {
			out.YAdvDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		_ = values[1:]
	}
	return out, offset + 2*size, err
}

// DeviceTable is either an DeviceHinting for standard fonts,
// or a DeviceVariation for variable fonts.
type DeviceTable interface {
	isDevice()
}

func (DeviceHinting) isDevice()   {}
func (DeviceVariation) isDevice() {}

type DeviceHinting struct {
	// with length endSize - startSize + 1
	Values []int8
	// correction range, in ppem
	StartSize, EndSize uint16
}

type DeviceVariation VariationStoreIndex

func parseDeviceTable(src []byte, offset uint16) (DeviceTable, error) {
	if L := len(src); L < int(offset)+6 {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset+6, L)
	}
	var header DeviceTableHeader
	header.mustParse(src[offset:])

	switch format := header.deltaFormat; format {
	case 1, 2, 3:
		var out DeviceHinting

		out.StartSize, out.EndSize = header.first, header.second
		if out.EndSize < out.StartSize {
			return nil, errors.New("invalid positionning device subtable")
		}

		nbPerUint16 := 16 / (1 << format) // 8, 4 or 2
		outLength := int(out.EndSize - out.StartSize + 1)
		var count int
		if outLength%nbPerUint16 == 0 {
			count = outLength / nbPerUint16
		} else {
			// add padding
			count = outLength/nbPerUint16 + 1
		}
		uint16s, err := parseUint16s(src[offset+6:], count)
		if err != nil {
			return nil, err
		}
		out.Values = make([]int8, count*nbPerUint16) // handle rounding error by reslicing after
		switch format {
		case 1:
			for i, u := range uint16s {
				uint16As2Bits(out.Values[i*8:], u)
			}
		case 2:
			for i, u := range uint16s {
				uint16As4Bits(out.Values[i*4:], u)
			}
		case 3:
			for i, u := range uint16s {
				uint16As8Bits(out.Values[i*2:], u)
			}
		}
		out.Values = out.Values[:outLength]
		return out, nil
	case 0x8000:
		return DeviceVariation{DeltaSetOuter: header.first, DeltaSetInner: header.second}, nil
	default:
		return nil, fmt.Errorf("unsupported positionning device subtable: %d", format)
	}
}

type PairValueRecord struct {
	SecondGlyph  GlyphID     //	Glyph ID of second glyph in the pair (first glyph is listed in the Coverage table).
	ValueRecord1 ValueRecord //	Positioning data for the first glyph in the pair.
	ValueRecord2 ValueRecord //	Positioning data for the second glyph in the pair.
}

type Class1Record struct {
	Class2Records []Class2Record //[class2Count]	Array of Class2 records, ordered by classes in classDef2.
}

type Class2Record struct {
	ValueRecord1 ValueRecord //	Positioning for first glyph — empty if valueFormat1 = 0.
	ValueRecord2 ValueRecord //	Positioning for second glyph — empty if valueFormat2 = 0.
}

func (AnchorFormat1) isAnchor() {}
func (AnchorFormat2) isAnchor() {}
func (AnchorFormat3) isAnchor() {}

type AnchorFormat1 struct {
	anchorFormat uint16 `unionTag:"1"`
	XCoordinate  int16  // Horizontal value, in design units
	YCoordinate  int16  // Vertical value, in design units
}

type AnchorFormat2 struct {
	anchorFormat uint16 `unionTag:"2"`
	XCoordinate  int16  // Horizontal value, in design units
	YCoordinate  int16  // Vertical value, in design units
	AnchorPoint  uint16 // Index to glyph contour point
}

type AnchorFormat3 struct {
	anchorFormat  uint16      `unionTag:"3"`
	XCoordinate   int16       // Horizontal value, in design units
	YCoordinate   int16       // Vertical value, in design units
	xDeviceOffset Offset16    // Offset to Device table (non-variable font) / VariationIndex table (variable font) for X coordinate, from beginning of Anchor table (may be NULL)
	yDeviceOffset Offset16    // Offset to Device table (non-variable font) / VariationIndex table (variable font) for Y coordinate, from beginning of Anchor table (may be NULL)
	XDevice       DeviceTable `isOpaque:""` // Offset to Device table (non-variable font) / VariationIndex table (variable font) for X coordinate, from beginning of Anchor table (may be NULL)
	YDevice       DeviceTable `isOpaque:""` // Offset to Device table (non-variable font) / VariationIndex table (variable font) for Y coordinate, from beginning of Anchor table (may be NULL)
}

func (af AnchorFormat3) customParseXDevice(src []byte) (int, error) {
	if af.xDeviceOffset == 0 {
		return 0, nil
	}
	var err error
	af.XDevice, err = parseDeviceTable(src, uint16(af.xDeviceOffset))
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (af AnchorFormat3) customParseYDevice(src []byte) (int, error) {
	if af.yDeviceOffset == 0 {
		return 0, nil
	}
	var err error
	af.YDevice, err = parseDeviceTable(src, uint16(af.yDeviceOffset))
	if err != nil {
		return 0, err
	}
	return 0, nil
}

type MarkArray struct {
	MarkRecords []MarkRecord `arrayCount:"FirstUint16"` //[markCount]	Array of MarkRecords, ordered by corresponding glyphs in the associated mark Coverage table.
	MarkAnchors []Anchor     `isOpaque:""`              // with same length as MarkRecords
}

func (ma *MarkArray) customParseMarkAnchors(src []byte) (int, error) {
	ma.MarkAnchors = make([]Anchor, len(ma.MarkRecords))
	var err error
	for i, rec := range ma.MarkRecords {
		if L := len(src); L < int(rec.markAnchorOffset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.markAnchorOffset, L)
		}
		ma.MarkAnchors[i], _, err = ParseAnchor(src[rec.markAnchorOffset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type MarkRecord struct {
	markClass        uint16   // Class defined for the associated mark.
	markAnchorOffset Offset16 // Offset to Anchor table, from beginning of MarkArray table.
}

// ------------------------------ parsing helpers ------------------------------

// write 8 elements
func uint16As2Bits(dst []int8, u uint16) {
	const mask = 0xFE // 11111110
	dst[0] = int8((0-uint8(u>>15&1))&mask | uint8(u>>14&1))
	dst[1] = int8((0-uint8(u>>13&1))&mask | uint8(u>>12&1))
	dst[2] = int8((0-uint8(u>>11&1))&mask | uint8(u>>10&1))
	dst[3] = int8((0-uint8(u>>9&1))&mask | uint8(u>>8&1))
	dst[4] = int8((0-uint8(u>>7&1))&mask | uint8(u>>6&1))
	dst[5] = int8((0-uint8(u>>5&1))&mask | uint8(u>>4&1))
	dst[6] = int8((0-uint8(u>>3&1))&mask | uint8(u>>2&1))
	dst[7] = int8((0-uint8(u>>1&1))&mask | uint8(u>>0&1))
}

// write 4 elements
func uint16As4Bits(dst []int8, u uint16) {
	const mask = 0xF8 // 11111000

	dst[0] = int8((0-uint8(u>>15&1))&mask | uint8(u>>12&0x07))
	dst[1] = int8((0-uint8(u>>11&1))&mask | uint8(u>>8&0x07))
	dst[2] = int8((0-uint8(u>>7&1))&mask | uint8(u>>4&0x07))
	dst[3] = int8((0-uint8(u>>3&1))&mask | uint8(u>>0&0x07))
}

// write 2 elements
func uint16As8Bits(dst []int8, u uint16) {
	dst[0] = int8(u >> 8)
	dst[1] = int8(u)
}

// parseUint16s interprets data as a (big endian) uint16 slice.
// It returns an error if [data] is not long enough for the given [count].
func parseUint16s(src []byte, count int) ([]uint16, error) {
	if L := len(src); L < 2*count {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", 2*count, L)
	}
	out := make([]uint16, count)
	for i := range out {
		out[i] = binary.BigEndian.Uint16(src[2*i:])
	}
	return out, nil
}
