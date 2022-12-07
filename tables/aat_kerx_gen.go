package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_kerx_src.go. DO NOT EDIT

func (item *KAAnchor) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.Mark = binary.BigEndian.Uint16(src[0:])
	item.Current = binary.BigEndian.Uint16(src[2:])
}

func (item *KAControl) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.Mark = binary.BigEndian.Uint16(src[0:])
	item.Current = binary.BigEndian.Uint16(src[2:])
}

func (item *KACoordinates) mustParse(src []byte) {
	_ = src[7] // early bound checking
	item.MarkX = int16(binary.BigEndian.Uint16(src[0:]))
	item.MarkY = int16(binary.BigEndian.Uint16(src[2:]))
	item.CurrentX = int16(binary.BigEndian.Uint16(src[4:]))
	item.CurrentY = int16(binary.BigEndian.Uint16(src[6:]))
}

func (item *Kernx0Record) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.Left = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.Right = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.Value = int16(binary.BigEndian.Uint16(src[4:]))
}

func ParseAATLookupExt(src []byte, valuesCount int) (AATLookupExt, int, error) {
	var item AATLookupExt

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AATLookupExt: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 0:
		item, read, err = ParseAATLoopkupExt0(src[0:], valuesCount)
	case 10:
		item, read, err = ParseAATLoopkupExt10(src[0:])
	case 2:
		item, read, err = ParseAATLoopkupExt2(src[0:])
	case 4:
		item, read, err = ParseAATLoopkupExt4(src[0:])
	case 6:
		item, read, err = ParseAATLoopkupExt6(src[0:])
	case 8:
		item, read, err = ParseAATLoopkupExt8(src[0:])
	default:
		err = fmt.Errorf("unsupported AATLookupExt format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading AATLookupExt: %s", err)
	}

	return item, read, nil
}

func ParseAATLoopkupExt0(src []byte, valuesCount int) (AATLoopkupExt0, int, error) {
	var item AATLoopkupExt0
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AATLoopkupExt0: "+"EOF: expected length: 2, got %d", L)
	}
	item.version = binary.BigEndian.Uint16(src[0:])
	n += 2

	{

		if L := len(src); L < 2+valuesCount*4 {
			return item, 0, fmt.Errorf("reading AATLoopkupExt0: "+"EOF: expected length: %d, got %d", 2+valuesCount*4, L)
		}

		item.Values = make([]uint32, valuesCount) // allocation guarded by the previous check
		for i := range item.Values {
			item.Values[i] = binary.BigEndian.Uint32(src[2+i*4:])
		}
		n += valuesCount * 4
	}
	return item, n, nil
}

func ParseAATLoopkupExt10(src []byte) (AATLoopkupExt10, int, error) {
	var item AATLoopkupExt10
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading AATLoopkupExt10: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.unitSize = binary.BigEndian.Uint16(src[2:])
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[4:]))
	arrayLengthValues := int(binary.BigEndian.Uint16(src[6:]))
	n += 8

	{

		if L := len(src); L < 8+arrayLengthValues*4 {
			return item, 0, fmt.Errorf("reading AATLoopkupExt10: "+"EOF: expected length: %d, got %d", 8+arrayLengthValues*4, L)
		}

		item.Values = make([]uint32, arrayLengthValues) // allocation guarded by the previous check
		for i := range item.Values {
			item.Values[i] = binary.BigEndian.Uint32(src[8+i*4:])
		}
		n += arrayLengthValues * 4
	}
	return item, n, nil
}

func ParseAATLoopkupExt2(src []byte) (AATLoopkupExt2, int, error) {
	var item AATLoopkupExt2
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading AATLoopkupExt2: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.binSearchHeader.mustParse(src[2:])
	n += 12

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 12+arrayLength*8 {
			return item, 0, fmt.Errorf("reading AATLoopkupExt2: "+"EOF: expected length: %d, got %d", 12+arrayLength*8, L)
		}

		item.Records = make([]lookupRecordExt2, arrayLength) // allocation guarded by the previous check
		for i := range item.Records {
			item.Records[i].mustParse(src[12+i*8:])
		}
		n += arrayLength * 8
	}
	return item, n, nil
}

func ParseAATLoopkupExt4(src []byte) (AATLoopkupExt4, int, error) {
	var item AATLoopkupExt4
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading AATLoopkupExt4: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.binSearchHeader.mustParse(src[2:])
	n += 12

	{
		arrayLength := int(item.nUnits)

		offset := 12
		for i := 0; i < arrayLength; i++ {
			elem, read, err := parseLoopkupRecordExt4(src[offset:], src)
			if err != nil {
				return item, 0, fmt.Errorf("reading AATLoopkupExt4: %s", err)
			}
			item.Records = append(item.Records, elem)
			offset += read
		}
		n = offset
	}
	return item, n, nil
}

func ParseAATLoopkupExt6(src []byte) (AATLoopkupExt6, int, error) {
	var item AATLoopkupExt6
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading AATLoopkupExt6: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.binSearchHeader.mustParse(src[2:])
	n += 12

	{
		arrayLength := int(item.nUnits)

		if L := len(src); L < 12+arrayLength*6 {
			return item, 0, fmt.Errorf("reading AATLoopkupExt6: "+"EOF: expected length: %d, got %d", 12+arrayLength*6, L)
		}

		item.Records = make([]loopkupRecordExt6, arrayLength) // allocation guarded by the previous check
		for i := range item.Records {
			item.Records[i].mustParse(src[12+i*6:])
		}
		n += arrayLength * 6
	}
	return item, n, nil
}

func ParseAATLoopkupExt8(src []byte) (AATLoopkupExt8, int, error) {
	var item AATLoopkupExt8
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AATLoopkupExt8: "+"EOF: expected length: 2, got %d", L)
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
			return item, 0, fmt.Errorf("reading AATLoopkupExt8: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func ParseAATStateTableExt(src []byte, valuesCount int, entryDataSize int) (AATStateTableExt, int, error) {
	var item AATStateTableExt
	n := 0
	if L := len(src); L < 16 {
		return item, 0, fmt.Errorf("reading AATStateTableExt: "+"EOF: expected length: 16, got %d", L)
	}
	_ = src[15] // early bound checking
	item.stateSize = binary.BigEndian.Uint32(src[0:])
	offsetClass := int(binary.BigEndian.Uint32(src[4:]))
	item.stateArray = Offset32(binary.BigEndian.Uint32(src[8:]))
	item.entryTable = Offset32(binary.BigEndian.Uint32(src[12:]))
	n += 16

	{

		if offsetClass != 0 { // ignore null offset
			if L := len(src); L < offsetClass {
				return item, 0, fmt.Errorf("reading AATStateTableExt: "+"EOF: expected length: %d, got %d", offsetClass, L)
			}

			var (
				err  error
				read int
			)
			item.Class, read, err = ParseAATLookup(src[offsetClass:], valuesCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading AATStateTableExt: %s", err)
			}
			offsetClass += read
		}
	}
	{

		err := item.parseStates(src[:], valuesCount, entryDataSize)
		if err != nil {
			return item, 0, fmt.Errorf("reading AATStateTableExt: %s", err)
		}
	}
	{

		read, err := item.parseEntries(src[:], valuesCount, entryDataSize)
		if err != nil {
			return item, 0, fmt.Errorf("reading AATStateTableExt: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseKerx(src []byte, valuesCount int) (Kerx, int, error) {
	var item Kerx
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading Kerx: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.padding = binary.BigEndian.Uint16(src[2:])
	item.nTables = binary.BigEndian.Uint32(src[4:])
	n += 8

	{
		arrayLength := int(item.nTables)

		offset := 8
		for i := 0; i < arrayLength; i++ {
			elem, read, err := ParseKerxSubtable(src[offset:], valuesCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading Kerx: %s", err)
			}
			item.Tables = append(item.Tables, elem)
			offset += read
		}
		n = offset
	}
	return item, n, nil
}

func ParseKerxAnchorAnchors(src []byte, anchorsCount int) (KerxAnchorAnchors, int, error) {
	var item KerxAnchorAnchors
	n := 0
	{

		if L := len(src); L < anchorsCount*4 {
			return item, 0, fmt.Errorf("reading KerxAnchorAnchors: "+"EOF: expected length: %d, got %d", anchorsCount*4, L)
		}

		item.Anchors = make([]KAAnchor, anchorsCount) // allocation guarded by the previous check
		for i := range item.Anchors {
			item.Anchors[i].mustParse(src[i*4:])
		}
		n += anchorsCount * 4
	}
	return item, n, nil
}

func ParseKerxAnchorControls(src []byte, anchorsCount int) (KerxAnchorControls, int, error) {
	var item KerxAnchorControls
	n := 0
	{

		if L := len(src); L < anchorsCount*4 {
			return item, 0, fmt.Errorf("reading KerxAnchorControls: "+"EOF: expected length: %d, got %d", anchorsCount*4, L)
		}

		item.Anchors = make([]KAControl, anchorsCount) // allocation guarded by the previous check
		for i := range item.Anchors {
			item.Anchors[i].mustParse(src[i*4:])
		}
		n += anchorsCount * 4
	}
	return item, n, nil
}

func ParseKerxAnchorCoordinates(src []byte, anchorsCount int) (KerxAnchorCoordinates, int, error) {
	var item KerxAnchorCoordinates
	n := 0
	{

		if L := len(src); L < anchorsCount*8 {
			return item, 0, fmt.Errorf("reading KerxAnchorCoordinates: "+"EOF: expected length: %d, got %d", anchorsCount*8, L)
		}

		item.Anchors = make([]KACoordinates, anchorsCount) // allocation guarded by the previous check
		for i := range item.Anchors {
			item.Anchors[i].mustParse(src[i*8:])
		}
		n += anchorsCount * 8
	}
	return item, n, nil
}

func ParseKerxData0(src []byte, tupleCount int) (KerxData0, int, error) {
	var item KerxData0
	n := 0
	if L := len(src); L < 16 {
		return item, 0, fmt.Errorf("reading KerxData0: "+"EOF: expected length: 16, got %d", L)
	}
	_ = src[15] // early bound checking
	item.nPairs = binary.BigEndian.Uint32(src[0:])
	item.searchRange = binary.BigEndian.Uint32(src[4:])
	item.entrySelector = binary.BigEndian.Uint32(src[8:])
	item.rangeShift = binary.BigEndian.Uint32(src[12:])
	n += 16

	{
		arrayLength := int(item.nPairs)

		if L := len(src); L < 16+arrayLength*6 {
			return item, 0, fmt.Errorf("reading KerxData0: "+"EOF: expected length: %d, got %d", 16+arrayLength*6, L)
		}

		item.Pairs = make([]Kernx0Record, arrayLength) // allocation guarded by the previous check
		for i := range item.Pairs {
			item.Pairs[i].mustParse(src[16+i*6:])
		}
		n += arrayLength * 6
	}
	var err error
	n, err = item.parseEnd(src, tupleCount)
	if err != nil {
		return item, 0, fmt.Errorf("reading KerxData0: %s", err)
	}

	return item, n, nil
}

func ParseKerxData1(src []byte, tupleCount int, valuesCount int) (KerxData1, int, error) {
	var item KerxData1
	n := 0
	{
		var (
			err  error
			read int
		)
		item.AATStateTableExt, read, err = ParseAATStateTableExt(src[0:], int(valuesCount), int(2))
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData1: %s", err)
		}
		n += read
	}
	if L := len(src); L < n+4 {
		return item, 0, fmt.Errorf("reading KerxData1: "+"EOF: expected length: n + 4, got %d", L)
	}
	item.valueTable = Offset32(binary.BigEndian.Uint32(src[n:]))
	n += 4

	{

		err := item.parseValues(src[:], tupleCount, valuesCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData1: %s", err)
		}
	}
	return item, n, nil
}

func ParseKerxData2(src []byte, parentSrc []byte, valuesCount int) (KerxData2, int, error) {
	var item KerxData2
	n := 0
	if L := len(src); L < 16 {
		return item, 0, fmt.Errorf("reading KerxData2: "+"EOF: expected length: 16, got %d", L)
	}
	_ = src[15] // early bound checking
	item.rowWidth = binary.BigEndian.Uint32(src[0:])
	offsetLeft := int(binary.BigEndian.Uint32(src[4:]))
	offsetRight := int(binary.BigEndian.Uint32(src[8:]))
	item.array = Offset32(binary.BigEndian.Uint32(src[12:]))
	n += 16

	{

		if offsetLeft != 0 { // ignore null offset
			if L := len(parentSrc); L < offsetLeft {
				return item, 0, fmt.Errorf("reading KerxData2: "+"EOF: expected length: %d, got %d", offsetLeft, L)
			}

			var (
				err  error
				read int
			)
			item.Left, read, err = ParseAATLookup(parentSrc[offsetLeft:], valuesCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading KerxData2: %s", err)
			}
			offsetLeft += read
		}
	}
	{

		if offsetRight != 0 { // ignore null offset
			if L := len(parentSrc); L < offsetRight {
				return item, 0, fmt.Errorf("reading KerxData2: "+"EOF: expected length: %d, got %d", offsetRight, L)
			}

			var (
				err  error
				read int
			)
			item.Right, read, err = ParseAATLookup(parentSrc[offsetRight:], valuesCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading KerxData2: %s", err)
			}
			offsetRight += read
		}
	}
	{

		item.kerningData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func ParseKerxData4(src []byte, valuesCount int) (KerxData4, int, error) {
	var item KerxData4
	n := 0
	{
		var (
			err  error
			read int
		)
		item.AATStateTableExt, read, err = ParseAATStateTableExt(src[0:], int(valuesCount), int(2))
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData4: %s", err)
		}
		n += read
	}
	if L := len(src); L < n+4 {
		return item, 0, fmt.Errorf("reading KerxData4: "+"EOF: expected length: n + 4, got %d", L)
	}
	item.Flags = binary.BigEndian.Uint32(src[n:])
	n += 4

	{

		err := item.parseAnchors(src[:], valuesCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData4: %s", err)
		}
	}
	return item, n, nil
}

func ParseKerxData6(src []byte, parentSrc []byte, tupleCount int, valuesCount int) (KerxData6, int, error) {
	var item KerxData6
	n := 0
	if L := len(src); L < 24 {
		return item, 0, fmt.Errorf("reading KerxData6: "+"EOF: expected length: 24, got %d", L)
	}
	_ = src[23] // early bound checking
	item.flags = binary.BigEndian.Uint32(src[0:])
	item.rowCount = binary.BigEndian.Uint16(src[4:])
	item.columnCount = binary.BigEndian.Uint16(src[6:])
	item.rowIndexTableOffset = binary.BigEndian.Uint32(src[8:])
	item.columnIndexTableOffset = binary.BigEndian.Uint32(src[12:])
	item.kerningArrayOffset = binary.BigEndian.Uint32(src[16:])
	item.kerningVectorOffset = binary.BigEndian.Uint32(src[20:])
	n += 24

	{

		err := item.parseRow(src[:], parentSrc, tupleCount, valuesCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData6: %s", err)
		}
	}
	{

		err := item.parseColumn(src[:], parentSrc, tupleCount, valuesCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData6: %s", err)
		}
	}
	{

		err := item.parseKernings(src[:], parentSrc, tupleCount, valuesCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxData6: %s", err)
		}
	}
	return item, n, nil
}

func ParseKerxSubtable(src []byte, valuesCount int) (KerxSubtable, int, error) {
	var item KerxSubtable
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading KerxSubtable: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.length = binary.BigEndian.Uint32(src[0:])
	item.Coverage = binary.BigEndian.Uint16(src[4:])
	item.padding = src[6]
	item.version = kerxSTVersion(src[7])
	item.TupleCount = binary.BigEndian.Uint32(src[8:])
	n += 12

	{
		var (
			read int
			err  error
		)
		switch item.version {
		case kerxSTVersion0:
			item.Data, read, err = ParseKerxData0(src[12:], int(item.TupleCount))
		case kerxSTVersion1:
			item.Data, read, err = ParseKerxData1(src[12:], int(item.TupleCount), int(valuesCount))
		case kerxSTVersion2:
			item.Data, read, err = ParseKerxData2(src[12:], src, int(valuesCount))
		case kerxSTVersion4:
			item.Data, read, err = ParseKerxData4(src[12:], int(valuesCount))
		case kerxSTVersion6:
			item.Data, read, err = ParseKerxData6(src[12:], src, int(item.TupleCount), int(valuesCount))
		default:
			err = fmt.Errorf("unsupported KerxDataVersion %d", item.version)
		}
		if err != nil {
			return item, 0, fmt.Errorf("reading KerxSubtable: %s", err)
		}
		n += read
	}
	var err error
	n, err = item.parseEnd(src, valuesCount)
	if err != nil {
		return item, 0, fmt.Errorf("reading KerxSubtable: %s", err)
	}

	return item, n, nil
}

func (item *lookupRecordExt2) mustParse(src []byte) {
	_ = src[7] // early bound checking
	item.LastGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.Value = binary.BigEndian.Uint32(src[4:])
}

func (item *loopkupRecordExt6) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.Glyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.Value = binary.BigEndian.Uint32(src[2:])
}

func parseLoopkupRecordExt4(src []byte, parentSrc []byte) (loopkupRecordExt4, int, error) {
	var item loopkupRecordExt4
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading loopkupRecordExt4: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.LastGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.FirstGlyph = GlyphID(binary.BigEndian.Uint16(src[2:]))
	offsetValues := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if offsetValues != 0 { // ignore null offset
			if L := len(parentSrc); L < offsetValues {
				return item, 0, fmt.Errorf("reading loopkupRecordExt4: "+"EOF: expected length: %d, got %d", offsetValues, L)
			}

			arrayLength := int(item.nValues())

			if L := len(parentSrc); L < offsetValues+arrayLength*4 {
				return item, 0, fmt.Errorf("reading loopkupRecordExt4: "+"EOF: expected length: %d, got %d", offsetValues+arrayLength*4, L)
			}

			item.Values = make([]uint32, arrayLength) // allocation guarded by the previous check
			for i := range item.Values {
				item.Values[i] = binary.BigEndian.Uint32(parentSrc[offsetValues+i*4:])
			}
			offsetValues += arrayLength * 4
		}
	}
	return item, n, nil
}
