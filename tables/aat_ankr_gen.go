package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_ankr_src.go. DO NOT EDIT

func (item *AnrkAnchor) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.X = int16(binary.BigEndian.Uint16(src[0:]))
	item.Y = int16(binary.BigEndian.Uint16(src[2:]))
}

func (item *LookupRecord2) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.LastGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.Value = binary.BigEndian.Uint16(src[4:])
}

func ParseAATLookup(src []byte, valuesCount int) (AATLookup, int, error) {
	var item AATLookup

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AATLookup: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 0:
		item, read, err = ParseAATLoopkup0(src[0:], valuesCount)
	case 10:
		item, read, err = ParseAATLoopkup10(src[0:])
	case 2:
		item, read, err = ParseAATLoopkup2(src[0:])
	case 4:
		item, read, err = ParseAATLoopkup4(src[0:])
	case 6:
		item, read, err = ParseAATLoopkup6(src[0:])
	case 8:
		item, read, err = ParseAATLoopkup8(src[0:])
	default:
		err = fmt.Errorf("unsupported AATLookup format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading AATLookup: %s", err)
	}

	return item, read, nil
}

func ParseAATLookupRecord4(src []byte, parentSrc []byte) (AATLookupRecord4, int, error) {
	var item AATLookupRecord4
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading AATLookupRecord4: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.LastGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	offsetValues := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if offsetValues != 0 { // ignore null offset
			if L := len(parentSrc); L < offsetValues {
				return item, 0, fmt.Errorf("reading AATLookupRecord4: "+"EOF: expected length: %d, got %d", offsetValues, L)
			}

			arrayLength := int(item.nValues())

			if L := len(parentSrc); L < offsetValues+arrayLength*2 {
				return item, 0, fmt.Errorf("reading AATLookupRecord4: "+"EOF: expected length: %d, got %d", offsetValues+arrayLength*2, L)
			}

			item.Values = make([]uint16, arrayLength) // allocation guarded by the previous check
			for i := range item.Values {
				item.Values[i] = binary.BigEndian.Uint16(parentSrc[offsetValues+i*2:])
			}
			offsetValues += arrayLength * 2
		}
	}
	return item, n, nil
}

func ParseAATLoopkup0(src []byte, valuesCount int) (AATLoopkup0, int, error) {
	var item AATLoopkup0
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AATLoopkup0: "+"EOF: expected length: 2, got %d", L)
	}
	item.version = binary.BigEndian.Uint16(src[0:])
	n += 2

	{

		if L := len(src); L < 2+valuesCount*2 {
			return item, 0, fmt.Errorf("reading AATLoopkup0: "+"EOF: expected length: %d, got %d", 2+valuesCount*2, L)
		}

		item.Values = make([]uint16, valuesCount) // allocation guarded by the previous check
		for i := range item.Values {
			item.Values[i] = binary.BigEndian.Uint16(src[2+i*2:])
		}
		n += valuesCount * 2
	}
	return item, n, nil
}

func ParseAATLoopkup10(src []byte) (AATLoopkup10, int, error) {
	var item AATLoopkup10
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading AATLoopkup10: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.unitSize = binary.BigEndian.Uint16(src[2:])
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[4:]))
	arrayLengthValues := int(binary.BigEndian.Uint16(src[6:]))
	n += 8

	{

		if L := len(src); L < 8+arrayLengthValues*2 {
			return item, 0, fmt.Errorf("reading AATLoopkup10: "+"EOF: expected length: %d, got %d", 8+arrayLengthValues*2, L)
		}

		item.Values = make([]uint16, arrayLengthValues) // allocation guarded by the previous check
		for i := range item.Values {
			item.Values[i] = binary.BigEndian.Uint16(src[8+i*2:])
		}
		n += arrayLengthValues * 2
	}
	return item, n, nil
}

func ParseAATLoopkup2(src []byte) (AATLoopkup2, int, error) {
	var item AATLoopkup2
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading AATLoopkup2: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.binSearchHeader.mustParse(src[2:])
	n += 12

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 12+arrayLength*6 {
			return item, 0, fmt.Errorf("reading AATLoopkup2: "+"EOF: expected length: %d, got %d", 12+arrayLength*6, L)
		}

		item.Records = make([]LookupRecord2, arrayLength) // allocation guarded by the previous check
		for i := range item.Records {
			item.Records[i].mustParse(src[12+i*6:])
		}
		n += arrayLength * 6
	}
	return item, n, nil
}

func ParseAATLoopkup4(src []byte) (AATLoopkup4, int, error) {
	var item AATLoopkup4
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading AATLoopkup4: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.binSearchHeader.mustParse(src[2:])
	n += 12

	{
		arrayLength := int(item.nUnits - 1)

		offset := 12
		for i := 0; i < arrayLength; i++ {
			elem, read, err := ParseAATLookupRecord4(src[offset:], src)
			if err != nil {
				return item, 0, fmt.Errorf("reading AATLoopkup4: %s", err)
			}
			item.Records = append(item.Records, elem)
			offset += read
		}
		n = offset
	}
	return item, n, nil
}

func ParseAATLoopkup6(src []byte) (AATLoopkup6, int, error) {
	var item AATLoopkup6
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading AATLoopkup6: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.binSearchHeader.mustParse(src[2:])
	n += 12

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 12+arrayLength*4 {
			return item, 0, fmt.Errorf("reading AATLoopkup6: "+"EOF: expected length: %d, got %d", 12+arrayLength*4, L)
		}

		item.Records = make([]loopkupRecord6, arrayLength) // allocation guarded by the previous check
		for i := range item.Records {
			item.Records[i].mustParse(src[12+i*4:])
		}
		n += arrayLength * 4
	}
	return item, n, nil
}

func ParseAATLoopkup8(src []byte) (AATLoopkup8, int, error) {
	var item AATLoopkup8
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AATLoopkup8: "+"EOF: expected length: 2, got %d", L)
	}
	item.version = binary.BigEndian.Uint16(src[0:])
	n += 2

	{
		var (
			err  error
			read int
		)
		item.AATLoopkup8Data, read, err = ParseAATLoopkup8Data(src[2:])
		if err != nil {
			return item, 0, fmt.Errorf("reading AATLoopkup8: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func ParseAATLoopkup8Data(src []byte) (AATLoopkup8Data, int, error) {
	var item AATLoopkup8Data
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading AATLoopkup8Data: "+"EOF: expected length: 4, got %d", L)
	}
	_ = src[3] // early bound checking
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	arrayLengthValues := int(binary.BigEndian.Uint16(src[2:]))
	n += 4

	{

		if L := len(src); L < 4+arrayLengthValues*2 {
			return item, 0, fmt.Errorf("reading AATLoopkup8Data: "+"EOF: expected length: %d, got %d", 4+arrayLengthValues*2, L)
		}

		item.Values = make([]uint16, arrayLengthValues) // allocation guarded by the previous check
		for i := range item.Values {
			item.Values[i] = binary.BigEndian.Uint16(src[4+i*2:])
		}
		n += arrayLengthValues * 2
	}
	return item, n, nil
}

func ParseAnkr(src []byte, valuesCount int) (Ankr, int, error) {
	var item Ankr
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading Ankr: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.flags = binary.BigEndian.Uint16(src[2:])
	offsetLookupTable := int(binary.BigEndian.Uint32(src[4:]))
	offsetGlyphDataTable := int(binary.BigEndian.Uint32(src[8:]))
	n += 12

	{

		if offsetLookupTable != 0 { // ignore null offset
			if L := len(src); L < offsetLookupTable {
				return item, 0, fmt.Errorf("reading Ankr: "+"EOF: expected length: %d, got %d", offsetLookupTable, L)
			}

			var (
				err  error
				read int
			)
			item.LookupTable, read, err = ParseAATLookup(src[offsetLookupTable:], valuesCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading Ankr: %s", err)
			}
			offsetLookupTable += read
		}
	}
	{

		if offsetGlyphDataTable != 0 { // ignore null offset
			if L := len(src); L < offsetGlyphDataTable {
				return item, 0, fmt.Errorf("reading Ankr: "+"EOF: expected length: %d, got %d", offsetGlyphDataTable, L)
			}

			item.GlyphDataTable = src[offsetGlyphDataTable:]
		}
	}
	return item, n, nil
}

func ParseAnrkAnchor(src []byte) (AnrkAnchor, int, error) {
	var item AnrkAnchor
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading AnrkAnchor: "+"EOF: expected length: 4, got %d", L)
	}
	item.mustParse(src)
	n += 4
	return item, n, nil
}

func (item *binSearchHeader) mustParse(src []byte) {
	_ = src[9] // early bound checking
	item.unitSize = binary.BigEndian.Uint16(src[0:])
	item.nUnits = binary.BigEndian.Uint16(src[2:])
	item.searchRange = binary.BigEndian.Uint16(src[4:])
	item.entrySelector = binary.BigEndian.Uint16(src[6:])
	item.rangeShift = binary.BigEndian.Uint16(src[8:])
}

func (item *loopkupRecord6) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.Glyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.Value = binary.BigEndian.Uint16(src[2:])
}
