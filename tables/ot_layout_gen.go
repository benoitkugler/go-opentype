package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from ot_layout_src.go. DO NOT EDIT

func (item *FeatureVariationRecord) mustParse(src []byte) {
	_ = src[7] // early bound checking
	item.conditionSetOffset = binary.BigEndian.Uint32(src[0:])
	item.featureTableSubstitutionOffset = binary.BigEndian.Uint32(src[4:])
}

func ParseFeature(src []byte) (Feature, int, error) {
	var item Feature
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading Feature: "+"EOF: expected length: 4, got %d", L)
	}
	_ = src[3] // early bound checking
	item.featureParamsOffset = binary.BigEndian.Uint16(src[0:])
	arrayLengthItemLookupListIndices := int(binary.BigEndian.Uint16(src[2:]))
	n += 4

	{

		if L := len(src); L < 4+arrayLengthItemLookupListIndices*2 {
			return item, 0, fmt.Errorf("reading Feature: "+"EOF: expected length: %d, got %d", 4+arrayLengthItemLookupListIndices*2, L)
		}

		item.LookupListIndices = make([]uint16, arrayLengthItemLookupListIndices) // allocation guarded by the previous check
		for i := range item.LookupListIndices {
			item.LookupListIndices[i] = binary.BigEndian.Uint16(src[4+i*2:])
		}
		n += arrayLengthItemLookupListIndices * 2
	}
	return item, n, nil
}

func ParseFeatureVariation(src []byte) (FeatureVariation, int, error) {
	var item FeatureVariation
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading FeatureVariation: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.majorVersion = binary.BigEndian.Uint16(src[0:])
	item.minorVersion = binary.BigEndian.Uint16(src[2:])
	arrayLengthItemfeatureVariationRecords := int(binary.BigEndian.Uint32(src[4:]))
	n += 8

	{

		if L := len(src); L < 8+arrayLengthItemfeatureVariationRecords*8 {
			return item, 0, fmt.Errorf("reading FeatureVariation: "+"EOF: expected length: %d, got %d", 8+arrayLengthItemfeatureVariationRecords*8, L)
		}

		item.featureVariationRecords = make([]FeatureVariationRecord, arrayLengthItemfeatureVariationRecords) // allocation guarded by the previous check
		for i := range item.featureVariationRecords {
			item.featureVariationRecords[i].mustParse(src[8+i*8:])
		}
		n += arrayLengthItemfeatureVariationRecords * 8
	}
	{

		item.rawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func ParseLangSys(src []byte) (LangSys, int, error) {
	var item LangSys
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading LangSys: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.lookupOrderOffset = binary.BigEndian.Uint16(src[0:])
	item.RequiredFeatureIndex = binary.BigEndian.Uint16(src[2:])
	arrayLengthItemFeatureIndices := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemFeatureIndices*2 {
			return item, 0, fmt.Errorf("reading LangSys: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemFeatureIndices*2, L)
		}

		item.FeatureIndices = make([]uint16, arrayLengthItemFeatureIndices) // allocation guarded by the previous check
		for i := range item.FeatureIndices {
			item.FeatureIndices[i] = binary.BigEndian.Uint16(src[6+i*2:])
		}
		n += arrayLengthItemFeatureIndices * 2
	}
	return item, n, nil
}

func ParseLayout(src []byte) (Layout, int, error) {
	var item Layout
	n := 0
	if L := len(src); L < 10 {
		return item, 0, fmt.Errorf("reading Layout: "+"EOF: expected length: 10, got %d", L)
	}
	_ = src[9] // early bound checking
	item.majorVersion = binary.BigEndian.Uint16(src[0:])
	item.minorVersion = binary.BigEndian.Uint16(src[2:])
	offsetItemscriptList := int(binary.BigEndian.Uint16(src[4:]))
	offsetItemfeatureList := int(binary.BigEndian.Uint16(src[6:]))
	offsetItemlookupList := int(binary.BigEndian.Uint16(src[8:]))
	n += 10

	{

		read, err := item.parseFeatureVariations(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading Layout: %s", err)
		}
		n = read
	}
	{

		if offsetItemscriptList != 0 { // ignore null offset
			if L := len(src); L < offsetItemscriptList {
				return item, 0, fmt.Errorf("reading Layout: "+"EOF: expected length: %d, got %d", offsetItemscriptList, L)
			}

			var (
				err  error
				read int
			)
			item.scriptList, read, err = parseScriptList(src[offsetItemscriptList:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Layout: %s", err)
			}
			offsetItemscriptList += read

		}
	}
	{

		if offsetItemfeatureList != 0 { // ignore null offset
			if L := len(src); L < offsetItemfeatureList {
				return item, 0, fmt.Errorf("reading Layout: "+"EOF: expected length: %d, got %d", offsetItemfeatureList, L)
			}

			var (
				err  error
				read int
			)
			item.featureList, read, err = parseFeatureList(src[offsetItemfeatureList:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Layout: %s", err)
			}
			offsetItemfeatureList += read

		}
	}
	{

		if offsetItemlookupList != 0 { // ignore null offset
			if L := len(src); L < offsetItemlookupList {
				return item, 0, fmt.Errorf("reading Layout: "+"EOF: expected length: %d, got %d", offsetItemlookupList, L)
			}

			var (
				err  error
				read int
			)
			item.lookupList, read, err = parseLookupList(src[offsetItemlookupList:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Layout: %s", err)
			}
			offsetItemlookupList += read

		}
	}
	return item, n, nil
}

func ParseLookup(src []byte) (Lookup, int, error) {
	var item Lookup
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading Lookup: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.lookupType = binary.BigEndian.Uint16(src[0:])
	item.lookupFlag = binary.BigEndian.Uint16(src[2:])
	arrayLengthItemsubtableOffsets := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemsubtableOffsets*2 {
			return item, 0, fmt.Errorf("reading Lookup: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemsubtableOffsets*2, L)
		}

		item.subtableOffsets = make([]Offset16, arrayLengthItemsubtableOffsets) // allocation guarded by the previous check
		for i := range item.subtableOffsets {
			item.subtableOffsets[i] = Offset16(binary.BigEndian.Uint16(src[6+i*2:]))
		}
		n += arrayLengthItemsubtableOffsets * 2
	}
	if L := len(src); L < n+2 {
		return item, 0, fmt.Errorf("reading Lookup: "+"EOF: expected length: n + 2, got %d", L)
	}
	item.markFilteringSet = binary.BigEndian.Uint16(src[n:])
	n += 2

	{

		item.rawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func ParseScript(src []byte) (Script, int, error) {
	var item Script
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading Script: "+"EOF: expected length: 4, got %d", L)
	}
	_ = src[3] // early bound checking
	offsetItemDefaultLangSys := int(binary.BigEndian.Uint16(src[0:]))
	arrayLengthItemlangSysRecords := int(binary.BigEndian.Uint16(src[2:]))
	n += 4

	{

		if L := len(src); L < 4+arrayLengthItemlangSysRecords*6 {
			return item, 0, fmt.Errorf("reading Script: "+"EOF: expected length: %d, got %d", 4+arrayLengthItemlangSysRecords*6, L)
		}

		item.langSysRecords = make([]tagOffsetRecord, arrayLengthItemlangSysRecords) // allocation guarded by the previous check
		for i := range item.langSysRecords {
			item.langSysRecords[i].mustParse(src[4+i*6:])
		}
		n += arrayLengthItemlangSysRecords * 6
	}
	{

		read, err := item.parseLangSys(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading Script: %s", err)
		}
		n = read
	}
	{

		if offsetItemDefaultLangSys != 0 { // ignore null offset
			if L := len(src); L < offsetItemDefaultLangSys {
				return item, 0, fmt.Errorf("reading Script: "+"EOF: expected length: %d, got %d", offsetItemDefaultLangSys, L)
			}

			var (
				err  error
				read int
			)
			item.DefaultLangSys, read, err = ParseLangSys(src[offsetItemDefaultLangSys:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Script: %s", err)
			}
			offsetItemDefaultLangSys += read

		}
	}
	return item, n, nil
}

func parseFeatureList(src []byte) (featureList, int, error) {
	var item featureList
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading featureList: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthItemrecords := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthItemrecords*6 {
			return item, 0, fmt.Errorf("reading featureList: "+"EOF: expected length: %d, got %d", 2+arrayLengthItemrecords*6, L)
		}

		item.records = make([]tagOffsetRecord, arrayLengthItemrecords) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i].mustParse(src[2+i*6:])
		}
		n += arrayLengthItemrecords * 6
	}
	{

		read, err := item.parseFeatures(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading featureList: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func parseLookupList(src []byte) (lookupList, int, error) {
	var item lookupList
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading lookupList: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthItemrecords := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthItemrecords*2 {
			return item, 0, fmt.Errorf("reading lookupList: "+"EOF: expected length: %d, got %d", 2+arrayLengthItemrecords*2, L)
		}

		item.records = make([]Offset16, arrayLengthItemrecords) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i] = Offset16(binary.BigEndian.Uint16(src[2+i*2:]))
		}
		n += arrayLengthItemrecords * 2
	}
	{

		read, err := item.parseLookups(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading lookupList: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func parseScriptList(src []byte) (scriptList, int, error) {
	var item scriptList
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading scriptList: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthItemrecords := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthItemrecords*6 {
			return item, 0, fmt.Errorf("reading scriptList: "+"EOF: expected length: %d, got %d", 2+arrayLengthItemrecords*6, L)
		}

		item.records = make([]tagOffsetRecord, arrayLengthItemrecords) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i].mustParse(src[2+i*6:])
		}
		n += arrayLengthItemrecords * 6
	}
	{

		read, err := item.parseScripts(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading scriptList: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func (item *tagOffsetRecord) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.Tag = Tag(binary.BigEndian.Uint32(src[0:]))
	item.Offset = binary.BigEndian.Uint16(src[4:])
}
