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
	scriptList        scriptList        `offsetSize:"Offset16"` // Offset to ScriptList table, from beginning of GPOS table
	featureList       featureList       `offsetSize:"Offset16"` // Offset to FeatureList table, from beginning of GPOS table
	LookupList        lookupList        `offsetSize:"Offset16"` // Offset to LookupList table, from beginning of GPOS table
	featureVariations *FeatureVariation `isOpaque:""`           // Offset to FeatureVariations table, from beginning of GPOS table (may be NULL)
}

func (lt *Layout) parseFeatureVariations(src []byte) (int, error) {
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
	Scripts []Script          `isOpaque:""`
}

func (sl *scriptList) parseScripts(src []byte) (int, error) {
	sl.Scripts = make([]Script, len(sl.records))
	for i, rec := range sl.records {
		var err error
		if L := len(src); L < int(rec.Offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.Offset, L)
		}
		sl.Scripts[i], _, err = ParseScript(src[rec.Offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type Script struct {
	DefaultLangSys LangSys           `offsetSize:"Offset16"`    // Offset to default LangSys table, from beginning of Script table — may be NULL
	langSysRecords []tagOffsetRecord `arrayCount:"FirstUint16"` // [langSysCount]	Array of LangSysRecords, listed alphabetically by LangSys tag
	LangSys        []LangSys         `isOpaque:""`
}

func (sc *Script) parseLangSys(src []byte) (int, error) {
	sc.LangSys = make([]LangSys, len(sc.langSysRecords))
	for i, rec := range sc.langSysRecords {
		var err error
		if L := len(src); L < int(rec.Offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.Offset, L)
		}
		sc.LangSys[i], _, err = ParseLangSys(src[rec.Offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type LangSys struct {
	lookupOrderOffset    uint16   // = NULL (reserved for an offset to a reordering table)
	RequiredFeatureIndex uint16   // Index of a feature required for this language system; if no required features = 0xFFFF
	FeatureIndices       []uint16 `arrayCount:"FirstUint16"` // [featureIndexCount]	Array of indices into the FeatureList, in arbitrary order
}

type featureList struct {
	records  []tagOffsetRecord `arrayCount:"FirstUint16"` // Array of FeatureRecords — zero-based (first feature has FeatureIndex = 0), listed alphabetically by feature tag
	Features []Feature         `isOpaque:""`
}

func (fl *featureList) parseFeatures(src []byte) (int, error) {
	fl.Features = make([]Feature, len(fl.records))
	for i, rec := range fl.records {
		var err error
		if L := len(src); L < int(rec.Offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.Offset, L)
		}
		fl.Features[i], _, err = ParseFeature(src[rec.Offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type Feature struct {
	featureParamsOffset uint16   // Offset from start of Feature table to FeatureParams table, if defined for the feature and present, else NULL
	LookupListIndices   []uint16 `arrayCount:"FirstUint16"` // [lookupIndexCount]	Array of indices into the LookupList — zero-based (first lookup is LookupListIndex = 0)
}

type lookupList struct {
	records []Offset16 `arrayCount:"FirstUint16"` // Array of offsets to Lookup tables, from beginning of LookupList — zero based (first lookup is Lookup index = 0)
	Lookups []Lookup   `isOpaque:""`
}

func (ll *lookupList) parseLookups(src []byte) (int, error) {
	var err error
	ll.Lookups = make([]Lookup, len(ll.records))
	for i, offset := range ll.records {
		if L := len(src); L < int(offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		ll.Lookups[i], _, err = ParseLookup(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

// Lookup is the common format for GSUB and GPOS lookups
type Lookup struct {
	lookupType       uint16     // Different enumerations for GSUB and GPOS
	LookupFlag       uint16     // Lookup qualifiers
	subtableOffsets  []Offset16 `arrayCount:"FirstUint16"` // [subTableCount] Array of offsets to lookup subtables, from beginning of Lookup table
	MarkFilteringSet uint16     // Index (base 0) into GDEF mark glyph sets structure. This field is only present if the USE_MARK_FILTERING_SET lookup flag is set.
	rawData          []byte     `subsliceStart:"AtStart" arrayCount:"ToEnd"`
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
