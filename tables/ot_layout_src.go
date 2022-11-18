package tables

import (
	"encoding/binary"
	"fmt"
)

// Layout represents the common layout table used by GPOS and GSUB.
// The Features field contains all the features for this layout. However,
// the script and language determines which feature is used.
//
// See https://learn.microsoft.com/typography/opentype/spec/chapter2#organization
// See https://learn.microsoft.com/typography/opentype/spec/gpos
// See https://www.microsoft.com/typography/otspec/GSUB.htm
type Layout struct {
	majorVersion      uint16            // Major version of the GPOS table, = 1
	minorVersion      uint16            // Minor version of the GPOS table, = 0 or 1
	scriptListOffset  scriptList        `offsetSize:"Offset16"`               // Offset to ScriptList table, from beginning of GPOS table
	featureListOffset featureList       `offsetSize:"Offset16"`               // Offset to FeatureList table, from beginning of GPOS table
	lookupListOffset  lookupList        `offsetSize:"Offset16"`               // Offset to LookupList table, from beginning of GPOS table
	featureVariations *FeatureVariation `isOpaque:"" subsliceStart:"AtStart"` // Offset to FeatureVariations table, from beginning of GPOS table (may be NULL)
}

func (lt *Layout) customParseFeatureVariations(src []byte) (int, error) {
	const layoutHeaderSize = 2 + 2 + 2 + 2 + 2
	if lt.minorVersion != 1 {
		return 0, nil
	}
	if L := len(src); L < layoutHeaderSize+4 {
		return 0, fmt.Errorf("reading Layout: EOF: expected length: 4, got %d", L)
	}
	offset := binary.BigEndian.Uint32(src[layoutHeaderSize:])
	if offset == 0 {
		return 4, nil
	}
	if L := len(src); L < int(offset) {
		return 0, fmt.Errorf("reading Layout: EOF: expected length: %d, got %d", offset, L)
	}
	fv, read, err := ParseFeatureVariation(src[offset:])
	if err != nil {
		return 0, err
	}
	lt.featureVariations = &fv

	return read, nil
}

type tagOffsetRecord struct {
	Tag    Tag    // 4-byte script tag identifier
	Offset uint16 // Offset to object from beginning of list
}

type scriptList struct {
	records []tagOffsetRecord `arrayCount:"FirstUint16"` // Array of ScriptRecords, listed alphabetically by script tag
	rawData []byte            `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

type Script struct {
	defaultLangSysOffset uint16            // Offset to default LangSys table, from beginning of Script table — may be NULL
	langSysRecords       []tagOffsetRecord `arrayCount:"FirstUint16"` // [langSysCount]	Array of LangSysRecords, listed alphabetically by LangSys tag
	rawData              []byte            `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

type LangSys struct {
	lookupOrderOffset    uint16   // = NULL (reserved for an offset to a reordering table)
	requiredFeatureIndex uint16   // Index of a feature required for this language system; if no required features = 0xFFFF
	featureIndices       []uint16 `arrayCount:"FirstUint16"` // [featureIndexCount]	Array of indices into the FeatureList, in arbitrary order
}

type featureList struct {
	records []tagOffsetRecord `arrayCount:"FirstUint16"` // Array of FeatureRecords — zero-based (first feature has FeatureIndex = 0), listed alphabetically by feature tag
	rawData []byte            `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

type Feature struct {
	featureParamsOffset uint16   // Offset from start of Feature table to FeatureParams table, if defined for the feature and present, else NULL
	lookupListIndices   []uint16 `arrayCount:"FirstUint16"` // [lookupIndexCount]	Array of indices into the LookupList — zero-based (first lookup is LookupListIndex = 0)
}

type lookupList struct {
	records []uint16 `arrayCount:"FirstUint16"` // Array of offsets to Lookup tables, from beginning of LookupList — zero based (first lookup is Lookup index = 0)
	lookups []Lookup `isOpaque:"" subsliceStart:"AtStart"`
}

func (ll lookupList) customParseLookups(src []byte) (int, error) {
	var err error
	ll.lookups = make([]Lookup, len(ll.records))
	for i, offset := range ll.records {
		if L := len(src); L < int(offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		ll.lookups[i], _, err = ParseLookup(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

// Lookup is the common format for GSUB and GPOS lookups
type Lookup struct {
	lookupType       uint16   // Different enumerations for GSUB and GPOS
	lookupFlag       uint16   // Lookup qualifiers
	subtableOffsets  []uint16 `arrayCount:"FirstUint16"` // [subTableCount] Array of offsets to lookup subtables, from beginning of Lookup table
	markFilteringSet uint16   // Index (base 0) into GDEF mark glyph sets structure. This field is only present if the USE_MARK_FILTERING_SET lookup flag is set.
	rawData          []byte   `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

type FeatureVariation struct {
	majorVersion            uint16                   // Major version of the FeatureVariations table — set to 1.
	minorVersion            uint16                   // Minor version of the FeatureVariations table — set to 0.
	featureVariationRecords []FeatureVariationRecord `arrayCount:"FirstUint32"` //[featureVariationRecordCount]	Array of feature variation records.
	rawData                 []byte                   `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

type FeatureVariationRecord struct {
	conditionSetOffset             uint32 // Offset to a condition set table, from beginning of FeatureVariations table.
	featureTableSubstitutionOffset uint32 // Offset to a feature table substitution table, from beginning of the FeatureVariations table.
}

// The following are types shared by GSUB and GPOS tables

// Coverage specifies all the glyphs affected by a substitution or
// positioning operation described in a subtable.
// Conceptually is it a []GlyphIndex, but it may be implemented for efficiently.
// See https://learn.microsoft.com/typography/opentype/spec/chapter2#lookup-table
type Coverage struct {
	format coverageVersion
	data   coverageData `unionField:"format"`
}

type coverageVersion uint16

const (
	coverageVersion1 coverageVersion = 1
	coverageVersion2 coverageVersion = 2
)

type coverageData interface {
	isCoverage()
}

func (coverageData1) isCoverage() {}
func (coverageData2) isCoverage() {}

type coverageData1 struct {
	glyphs []glyphID `arrayCount:"FirstUint16"`
}

type coverageData2 struct {
	ranges []rangeRecord `arrayCount:"FirstUint16"`
}

type rangeRecord struct {
	startGlyphID       glyphID // First glyph ID in the range
	endGlyphID         glyphID // Last glyph ID in the range
	startCoverageIndex uint16  // Coverage Index of first glyph ID in range
}

type ClassDef struct {
	format classDefVersion
	data   classDefData `unionField:"format"`
}

type classDefVersion uint16

const (
	classDefVersion1 classDefVersion = 1
	classDefVersion2 classDefVersion = 2
)

type classDefData interface {
	isClassDefData()
}

func (classDefData1) isClassDefData() {}
func (classDefData2) isClassDefData() {}

type classDefData1 struct {
	startGlyphID    glyphID  // First glyph ID of the classValueArray
	classValueArray []uint16 `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type classDefData2 struct {
	classRangeRecords []classRangeRecord `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type classRangeRecord struct {
	startGlyphID glyphID // First glyph ID in the range
	endGlyphID   glyphID // Last glyph ID in the range
	class        uint16  // Applied to all glyphs in the range
}

// Lookups

type SequenceLookupRecord struct {
	sequenceIndex   uint16 // Index (zero-based) into the input glyph sequence
	lookupListIndex uint16 // Index (zero-based) into the LookupList
}

// binarygen: startOffset=2
// Format identifier: format = 1
type SequenceContextFormat1 struct {
	Coverage          Coverage `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of SequenceContextFormat1 table
	seqRuleSetOffsets []uint16 `arrayCount:"FirstUint16"` //[seqRuleSetCount]	Array of offsets to SequenceRuleSet tables, from beginning of SequenceContextFormat1 table (offsets may be NULL)
}

type SequenceRuleSet struct {
	seqRuleOffsets []uint16 `arrayCount:"FirstUint16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type SequenceRule struct {
	glyphCount       uint16                 // Number of glyphs in the input glyph sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	inputSequence    []glyphID              `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Array of input glyph IDs—starting with the second glyph
	seqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of Sequence lookup records
}

// binarygen: startOffset=2
// Format identifier: format = 2
type SequenceContextFormat2 struct {
	coverage               Coverage `offsetSize:"Offset16"`    //	Offset to Coverage table, from beginning of SequenceContextFormat2 table
	classDef               ClassDef `offsetSize:"Offset16"`    //	Offset to ClassDef table, from beginning of SequenceContextFormat2 table
	classSeqRuleSetOffsets []uint16 `arrayCount:"FirstUint16"` //[classSeqRuleSetCount]	Array of offsets to ClassSequenceRuleSet tables, from beginning of SequenceContextFormat2 table (may be NULL)
}

type ClassSequenceRuleSet struct {
	classSeqRuleOffsets []uint16 //	Array of offsets to ClassSequenceRule tables, from beginning of ClassSequenceRuleSet table
	rawData             []byte   `sliceStart:"AtStart"`
}

type ClassSequenceRule struct {
	glyphCount       uint16               // Number of glyphs to be matched
	seqLookupCount   uint16               // Number of SequenceLookupRecords
	inputSequence    uint16               `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Sequence of classes to be matched to the input glyph sequence, beginning with the second glyph position
	seqLookupRecords SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of SequenceLookupRecords
}
