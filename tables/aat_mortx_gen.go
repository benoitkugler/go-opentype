package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_mortx_src.go. DO NOT EDIT

func ParseMorx(src []byte) (Morx, int, error) {
	var item Morx
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading Morx: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.unused = binary.BigEndian.Uint16(src[2:])
	item.nChains = binary.BigEndian.Uint32(src[4:])
	n += 8

	{
		arrayLength := int(item.nChains)

		offset := 8
		for i := 0; i < arrayLength; i++ {
			elem, read, err := parseMorxChain(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Morx: %s", err)
			}
			item.chains = append(item.chains, elem)
			offset += read
		}
		n = offset
	}
	return item, n, nil
}

func ParseMorxChainSubtable(src []byte, valuesCount int) (MorxChainSubtable, int, error) {
	var item MorxChainSubtable
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading MorxChainSubtable: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.length = binary.BigEndian.Uint32(src[0:])
	item.coverage = src[4]
	item.ignored[0] = src[5]
	item.ignored[1] = src[6]
	item.version = morxSubtableVersion(src[7])
	item.subFeatureFlags = binary.BigEndian.Uint32(src[8:])
	n += 12

	{
		var (
			read int
			err  error
		)
		switch item.version {
		case morxSubtableVersionContextual:
			item.subtableContent, read, err = parseMorxSubtableContextual(src[12:])
		case morxSubtableVersionInsertion:
			item.subtableContent, read, err = parseMorxSubtableInsertion(src[12:])
		case morxSubtableVersionLigature:
			item.subtableContent, read, err = parseMorxSubtableLigature(src[12:])
		case morxSubtableVersionNonContextual:
			item.subtableContent, read, err = parseMorxSubtableNonContextual(src[12:], valuesCount)
		case morxSubtableVersionRearrangement:
			item.subtableContent, read, err = parseMorxSubtableRearrangement(src[12:])
		default:
			err = fmt.Errorf("unsupported morxSubtableVersion %d", item.version)
		}
		if err != nil {
			return item, 0, fmt.Errorf("reading MorxChainSubtable: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func (item *aatFeature) mustParse(src []byte) {
	_ = src[11] // early bound checking
	item.featureType = binary.BigEndian.Uint16(src[0:])
	item.featureSetting = binary.BigEndian.Uint16(src[2:])
	item.enableFlags = binary.BigEndian.Uint32(src[4:])
	item.disableFlags = binary.BigEndian.Uint32(src[8:])
}

func (item *aatSTXHeader) mustParse(src []byte) {
	_ = src[15] // early bound checking
	item.stateSize = binary.BigEndian.Uint32(src[0:])
	item.classTable = binary.BigEndian.Uint32(src[4:])
	item.stateArray = binary.BigEndian.Uint32(src[8:])
	item.entryTable = binary.BigEndian.Uint32(src[12:])
}

func (item *binSearchHeader) mustParse(src []byte) {
	_ = src[9] // early bound checking
	item.unitSize = binary.BigEndian.Uint16(src[0:])
	item.nUnits = binary.BigEndian.Uint16(src[2:])
	item.searchRange = binary.BigEndian.Uint16(src[4:])
	item.entrySelector = binary.BigEndian.Uint16(src[6:])
	item.rangeShift = binary.BigEndian.Uint16(src[8:])
}

func (item *lookupRecord2) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.lastGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.firstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.value = binary.BigEndian.Uint16(src[4:])
}

func (item *loopkupRecord4) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.lastGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.firstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.offsetToValues = binary.BigEndian.Uint16(src[4:])
}

func (item *loopkupRecord6) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.glyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.value = binary.BigEndian.Uint16(src[2:])
}

func parseAatLookup(src []byte, valuesCount int) (aatLookup, int, error) {
	var item aatLookup
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading aatLookup: "+"EOF: expected length: 2, got %d", L)
	}
	item.version = aatLookupVersion(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{
		var (
			read int
			err  error
		)
		switch item.version {
		case aatLookupTableVersion0:
			item.table, read, err = parseAatLookupTable0(src[2:], valuesCount)
		case aatLookupTableVersion10:
			item.table, read, err = parseAatLookupTable10(src[2:])
		case aatLookupTableVersion2:
			item.table, read, err = parseAatLookupTable2(src[2:])
		case aatLookupTableVersion4:
			item.table, read, err = parseAatLookupTable4(src[2:])
		case aatLookupTableVersion6:
			item.table, read, err = parseAatLookupTable6(src[2:])
		case aatLookupTableVersion8:
			item.table, read, err = parseAatLookupTable8(src[2:])
		default:
			err = fmt.Errorf("unsupported aatLookupTableVersion %d", item.version)
		}
		if err != nil {
			return item, 0, fmt.Errorf("reading aatLookup: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func parseAatLookupTable0(src []byte, valuesCount int) (aatLookupTable0, int, error) {
	var item aatLookupTable0
	n := 0
	{

		if L := len(src); L < valuesCount*2 {
			return item, 0, fmt.Errorf("reading aatLookupTable0: "+"EOF: expected length: %d, got %d", valuesCount*2, L)
		}

		item.values = make([]uint16, valuesCount) // allocation guarded by the previous check
		for i := range item.values {
			item.values[i] = binary.BigEndian.Uint16(src[i*2:])
		}
		n += valuesCount * 2
	}
	return item, n, nil
}

func parseAatLookupTable10(src []byte) (aatLookupTable10, int, error) {
	var item aatLookupTable10
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading aatLookupTable10: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.unitSize = binary.BigEndian.Uint16(src[0:])
	item.firstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	arrayLengthValues := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthValues*2 {
			return item, 0, fmt.Errorf("reading aatLookupTable10: "+"EOF: expected length: %d, got %d", 6+arrayLengthValues*2, L)
		}

		item.values = make([]uint16, arrayLengthValues) // allocation guarded by the previous check
		for i := range item.values {
			item.values[i] = binary.BigEndian.Uint16(src[6+i*2:])
		}
		n += arrayLengthValues * 2
	}
	return item, n, nil
}

func parseAatLookupTable2(src []byte) (aatLookupTable2, int, error) {
	var item aatLookupTable2
	n := 0
	if L := len(src); L < 10 {
		return item, 0, fmt.Errorf("reading aatLookupTable2: "+"EOF: expected length: 10, got %d", L)
	}
	item.binSearchHeader.mustParse(src[0:])
	n += 10

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 10+arrayLength*6 {
			return item, 0, fmt.Errorf("reading aatLookupTable2: "+"EOF: expected length: %d, got %d", 10+arrayLength*6, L)
		}

		item.records = make([]lookupRecord2, arrayLength) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i].mustParse(src[10+i*6:])
		}
		n += arrayLength * 6
	}
	return item, n, nil
}

func parseAatLookupTable4(src []byte) (aatLookupTable4, int, error) {
	var item aatLookupTable4
	n := 0
	if L := len(src); L < 10 {
		return item, 0, fmt.Errorf("reading aatLookupTable4: "+"EOF: expected length: 10, got %d", L)
	}
	item.binSearchHeader.mustParse(src[0:])
	n += 10

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 10+arrayLength*6 {
			return item, 0, fmt.Errorf("reading aatLookupTable4: "+"EOF: expected length: %d, got %d", 10+arrayLength*6, L)
		}

		item.records = make([]loopkupRecord4, arrayLength) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i].mustParse(src[10+i*6:])
		}
		n += arrayLength * 6
	}
	{

		item.rawData = src[n:]
		n = len(src)
	}
	return item, n, nil
}

func parseAatLookupTable6(src []byte) (aatLookupTable6, int, error) {
	var item aatLookupTable6
	n := 0
	if L := len(src); L < 10 {
		return item, 0, fmt.Errorf("reading aatLookupTable6: "+"EOF: expected length: 10, got %d", L)
	}
	item.binSearchHeader.mustParse(src[0:])
	n += 10

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 10+arrayLength*4 {
			return item, 0, fmt.Errorf("reading aatLookupTable6: "+"EOF: expected length: %d, got %d", 10+arrayLength*4, L)
		}

		item.records = make([]loopkupRecord6, arrayLength) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i].mustParse(src[10+i*4:])
		}
		n += arrayLength * 4
	}
	return item, n, nil
}

func parseMorxChain(src []byte) (morxChain, int, error) {
	var item morxChain
	n := 0
	if L := len(src); L < 16 {
		return item, 0, fmt.Errorf("reading morxChain: "+"EOF: expected length: 16, got %d", L)
	}
	_ = src[15] // early bound checking
	item.flags = binary.BigEndian.Uint32(src[0:])
	item.chainLength = binary.BigEndian.Uint32(src[4:])
	item.nFeatureEntries = binary.BigEndian.Uint32(src[8:])
	item.nSubtable = binary.BigEndian.Uint32(src[12:])
	n += 16

	{
		arrayLength := int(item.nFeatureEntries)

		if L := len(src); L < 16+arrayLength*12 {
			return item, 0, fmt.Errorf("reading morxChain: "+"EOF: expected length: %d, got %d", 16+arrayLength*12, L)
		}

		item.features = make([]aatFeature, arrayLength) // allocation guarded by the previous check
		for i := range item.features {
			item.features[i].mustParse(src[16+i*12:])
		}
		n += arrayLength * 12
	}
	{

		L := int(item.chainLength)
		if len(src) < L {
			return item, 0, fmt.Errorf("reading morxChain: "+"EOF: expected length: %d, got %d", L, len(src))
		}
		item.subtablesData = src[n:L]
		n = L
	}
	return item, n, nil
}

func parseMorxSubtableContextual(src []byte) (morxSubtableContextual, int, error) {
	var item morxSubtableContextual
	n := 0
	if L := len(src); L < 20 {
		return item, 0, fmt.Errorf("reading morxSubtableContextual: "+"EOF: expected length: 20, got %d", L)
	}
	_ = src[19] // early bound checking
	item.aatSTXHeader.mustParse(src[0:])
	item.substitutionTableOffset = binary.BigEndian.Uint32(src[16:])
	n += 20

	{

		item.rawData = src[20:]
		n = len(src)
	}
	return item, n, nil
}

func parseMorxSubtableInsertion(src []byte) (morxSubtableInsertion, int, error) {
	var item morxSubtableInsertion
	n := 0
	if L := len(src); L < 20 {
		return item, 0, fmt.Errorf("reading morxSubtableInsertion: "+"EOF: expected length: 20, got %d", L)
	}
	_ = src[19] // early bound checking
	item.aatSTXHeader.mustParse(src[0:])
	item.insertionActionOffset = binary.BigEndian.Uint32(src[16:])
	n += 20

	{

		item.rawData = src[20:]
		n = len(src)
	}
	return item, n, nil
}

func parseMorxSubtableLigature(src []byte) (morxSubtableLigature, int, error) {
	var item morxSubtableLigature
	n := 0
	if L := len(src); L < 28 {
		return item, 0, fmt.Errorf("reading morxSubtableLigature: "+"EOF: expected length: 28, got %d", L)
	}
	_ = src[27] // early bound checking
	item.aatSTXHeader.mustParse(src[0:])
	item.ligActionOffset = binary.BigEndian.Uint32(src[16:])
	item.componentOffset = binary.BigEndian.Uint32(src[20:])
	item.ligatureOffset = binary.BigEndian.Uint32(src[24:])
	n += 28

	{

		item.rawData = src[28:]
		n = len(src)
	}
	return item, n, nil
}

func parseMorxSubtableNonContextual(src []byte, valuesCount int) (morxSubtableNonContextual, int, error) {
	var item morxSubtableNonContextual
	n := 0
	{
		var (
			err  error
			read int
		)
		item.table, read, err = parseAatLookup(src[0:], valuesCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading morxSubtableNonContextual: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func parseMorxSubtableRearrangement(src []byte) (morxSubtableRearrangement, int, error) {
	var item morxSubtableRearrangement
	n := 0
	if L := len(src); L < 16 {
		return item, 0, fmt.Errorf("reading morxSubtableRearrangement: "+"EOF: expected length: 16, got %d", L)
	}
	item.aatSTXHeader.mustParse(src[0:])
	n += 16

	{

		item.rawData = src[16:]
		n = len(src)
	}
	return item, n, nil
}
