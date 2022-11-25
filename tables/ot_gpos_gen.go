package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from ot_gpos_src.go. DO NOT EDIT

func (item *AnchorFormat1) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.anchorFormat = binary.BigEndian.Uint16(src[0:])
	item.XCoordinate = int16(binary.BigEndian.Uint16(src[2:]))
	item.YCoordinate = int16(binary.BigEndian.Uint16(src[4:]))
}

func (item *AnchorFormat2) mustParse(src []byte) {
	_ = src[7] // early bound checking
	item.anchorFormat = binary.BigEndian.Uint16(src[0:])
	item.XCoordinate = int16(binary.BigEndian.Uint16(src[2:]))
	item.YCoordinate = int16(binary.BigEndian.Uint16(src[4:]))
	item.AnchorPoint = binary.BigEndian.Uint16(src[6:])
}

func (item *ClassRangeRecord) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.StartGlyphID = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.EndGlyphID = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.Class = binary.BigEndian.Uint16(src[4:])
}

func (item *DeviceTableHeader) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.first = binary.BigEndian.Uint16(src[0:])
	item.second = binary.BigEndian.Uint16(src[2:])
	item.deltaFormat = binary.BigEndian.Uint16(src[4:])
}

func (item *MarkRecord) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.markClass = binary.BigEndian.Uint16(src[0:])
	item.markAnchorOffset = Offset16(binary.BigEndian.Uint16(src[2:]))
}

func ParseAnchor(src []byte) (Anchor, int, error) {
	var item Anchor

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading Anchor: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 1:
		item, read, err = ParseAnchorFormat1(src[0:])
	case 2:
		item, read, err = ParseAnchorFormat2(src[0:])
	case 3:
		item, read, err = ParseAnchorFormat3(src[0:])
	default:
		err = fmt.Errorf("unsupported Anchor format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading Anchor: %s", err)
	}

	return item, read, nil
}

func ParseAnchorFormat1(src []byte) (AnchorFormat1, int, error) {
	var item AnchorFormat1
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading AnchorFormat1: "+"EOF: expected length: 6, got %d", L)
	}
	item.mustParse(src)
	n += 6
	return item, n, nil
}

func ParseAnchorFormat2(src []byte) (AnchorFormat2, int, error) {
	var item AnchorFormat2
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading AnchorFormat2: "+"EOF: expected length: 8, got %d", L)
	}
	item.mustParse(src)
	n += 8
	return item, n, nil
}

func ParseAnchorFormat3(src []byte) (AnchorFormat3, int, error) {
	var item AnchorFormat3
	n := 0
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading AnchorFormat3: "+"EOF: expected length: 10, got %d", L)
		}
		_ = src[9] // early bound checking
		item.anchorFormat = binary.BigEndian.Uint16(src[0:])
		item.XCoordinate = int16(binary.BigEndian.Uint16(src[2:]))
		item.YCoordinate = int16(binary.BigEndian.Uint16(src[4:]))
		item.xDeviceOffset = Offset16(binary.BigEndian.Uint16(src[6:]))
		item.yDeviceOffset = Offset16(binary.BigEndian.Uint16(src[8:]))
		n += 10
	}
	{

		read, err := item.customParseXDevice(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading AnchorFormat3: %s", err)
		}
		n = read
	}
	{

		read, err := item.customParseYDevice(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading AnchorFormat3: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseBaseArray(src []byte, offsetsCount int) (BaseArray, int, error) {
	var item BaseArray
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading BaseArray: "+"EOF: expected length: 2, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[0:]))
		n += 2

		offset := 2
		for i := 0; i < arrayLength; i++ {
			elem, read, err := parseAnchorOffsets(src[offset:], offsetsCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading BaseArray: %s", err)
			}
			item.baseRecords = append(item.baseRecords, elem)
			offset += read
		}
		n = offset
	}
	{

		read, err := item.customParseBaseAnchors(src[:], offsetsCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading BaseArray: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseClassDef(src []byte) (ClassDef, int, error) {
	var item ClassDef

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading ClassDef: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 1:
		item, read, err = ParseClassDef1(src[0:])
	case 2:
		item, read, err = ParseClassDef2(src[0:])
	default:
		err = fmt.Errorf("unsupported ClassDef format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading ClassDef: %s", err)
	}

	return item, read, nil
}

func ParseClassDef1(src []byte) (ClassDef1, int, error) {
	var item ClassDef1
	n := 0
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading ClassDef1: "+"EOF: expected length: 4, got %d", L)
		}
		_ = src[3] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.StartGlyphID = GlyphID(binary.BigEndian.Uint16(src[2:]))
		n += 4
	}
	{
		if L := len(src); L < 6 {
			return item, 0, fmt.Errorf("reading ClassDef1: "+"EOF: expected length: 6, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[4:]))
		n += 2

		if L := len(src); L < 6+arrayLength*2 {
			return item, 0, fmt.Errorf("reading ClassDef1: "+"EOF: expected length: %d, got %d", 6+arrayLength*2, L)
		}

		item.ClassValueArray = make([]uint16, arrayLength) // allocation guarded by the previous check
		for i := range item.ClassValueArray {
			item.ClassValueArray[i] = binary.BigEndian.Uint16(src[6+i*2:])
		}
		n += arrayLength * 2
	}
	return item, n, nil
}

func ParseClassDef2(src []byte) (ClassDef2, int, error) {
	var item ClassDef2
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading ClassDef2: "+"EOF: expected length: 2, got %d", L)
		}
		item.format = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading ClassDef2: "+"EOF: expected length: 4, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[2:]))
		n += 2

		if L := len(src); L < 4+arrayLength*6 {
			return item, 0, fmt.Errorf("reading ClassDef2: "+"EOF: expected length: %d, got %d", 4+arrayLength*6, L)
		}

		item.ClassRangeRecords = make([]ClassRangeRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.ClassRangeRecords {
			item.ClassRangeRecords[i].mustParse(src[4+i*6:])
		}
		n += arrayLength * 6
	}
	return item, n, nil
}

func ParseClassRangeRecord(src []byte) (ClassRangeRecord, int, error) {
	var item ClassRangeRecord
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading ClassRangeRecord: "+"EOF: expected length: 6, got %d", L)
	}
	item.mustParse(src)
	n += 6
	return item, n, nil
}

func ParseCoverage(src []byte) (Coverage, int, error) {
	var item Coverage

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading Coverage: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 1:
		item, read, err = ParseCoverage1(src[0:])
	case 2:
		item, read, err = ParseCoverage2(src[0:])
	default:
		err = fmt.Errorf("unsupported Coverage format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading Coverage: %s", err)
	}

	return item, read, nil
}

func ParseCoverage1(src []byte) (Coverage1, int, error) {
	var item Coverage1
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading Coverage1: "+"EOF: expected length: 2, got %d", L)
		}
		item.format = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading Coverage1: "+"EOF: expected length: 4, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[2:]))
		n += 2

		if L := len(src); L < 4+arrayLength*2 {
			return item, 0, fmt.Errorf("reading Coverage1: "+"EOF: expected length: %d, got %d", 4+arrayLength*2, L)
		}

		item.Glyphs = make([]GlyphID, arrayLength) // allocation guarded by the previous check
		for i := range item.Glyphs {
			item.Glyphs[i] = GlyphID(binary.BigEndian.Uint16(src[4+i*2:]))
		}
		n += arrayLength * 2
	}
	return item, n, nil
}

func ParseCoverage2(src []byte) (Coverage2, int, error) {
	var item Coverage2
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading Coverage2: "+"EOF: expected length: 2, got %d", L)
		}
		item.format = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading Coverage2: "+"EOF: expected length: 4, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[2:]))
		n += 2

		if L := len(src); L < 4+arrayLength*6 {
			return item, 0, fmt.Errorf("reading Coverage2: "+"EOF: expected length: %d, got %d", 4+arrayLength*6, L)
		}

		item.Ranges = make([]RangeRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.Ranges {
			item.Ranges[i].mustParse(src[4+i*6:])
		}
		n += arrayLength * 6
	}
	return item, n, nil
}

func ParseCursivePos(src []byte) (CursivePos, int, error) {
	var item CursivePos
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading CursivePos: "+"EOF: expected length: 2, got %d", L)
		}
		item.posFormat = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading CursivePos: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading CursivePos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading CursivePos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 6 {
			return item, 0, fmt.Errorf("reading CursivePos: "+"EOF: expected length: 6, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[4:]))
		n += 2

		if L := len(src); L < 6+arrayLength*4 {
			return item, 0, fmt.Errorf("reading CursivePos: "+"EOF: expected length: %d, got %d", 6+arrayLength*4, L)
		}

		item.entryExitRecords = make([]entryExitRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.entryExitRecords {
			item.entryExitRecords[i].mustParse(src[6+i*4:])
		}
		n += arrayLength * 4
	}
	{

		read, err := item.customParseEntryExits(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading CursivePos: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseDeviceTableHeader(src []byte) (DeviceTableHeader, int, error) {
	var item DeviceTableHeader
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading DeviceTableHeader: "+"EOF: expected length: 6, got %d", L)
	}
	item.mustParse(src)
	n += 6
	return item, n, nil
}

func ParseEntryExit(src []byte) (EntryExit, int, error) {
	var item EntryExit
	n := 0
	{
		var (
			err  error
			read int
		)
		item.EntryAnchor, read, err = ParseAnchor(src[0:])
		if err != nil {
			return item, 0, fmt.Errorf("reading EntryExit: %s", err)
		}
		n += read
	}
	{
		var (
			err  error
			read int
		)
		item.ExitAnchor, read, err = ParseAnchor(src[n:])
		if err != nil {
			return item, 0, fmt.Errorf("reading EntryExit: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func ParseLigatureArray(src []byte, offsetsCount int) (LigatureArray, int, error) {
	var item LigatureArray
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading LigatureArray: "+"EOF: expected length: 2, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[0:]))
		n += 2

		if L := len(src); L < 2+arrayLength*2 {
			return item, 0, fmt.Errorf("reading LigatureArray: "+"EOF: expected length: %d, got %d", 2+arrayLength*2, L)
		}

		item.LigatureAttachs = make([]LigatureAttach, arrayLength) // allocation guarded by the previous check
		for i := range item.LigatureAttachs {
			offset := int(binary.BigEndian.Uint16(src[2+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading LigatureArray: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.LigatureAttachs[i], _, err = ParseLigatureAttach(src[offset:], offsetsCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading LigatureArray: %s", err)
			}

		}
		n += arrayLength * 2
	}
	return item, n, nil
}

func ParseLigatureAttach(src []byte, offsetsCount int) (LigatureAttach, int, error) {
	var item LigatureAttach
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading LigatureAttach: "+"EOF: expected length: 2, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[0:]))
		n += 2

		offset := 2
		for i := 0; i < arrayLength; i++ {
			elem, read, err := parseAnchorOffsets(src[offset:], offsetsCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading LigatureAttach: %s", err)
			}
			item.componentRecords = append(item.componentRecords, elem)
			offset += read
		}
		n = offset
	}
	{

		read, err := item.customParseComponentAnchors(src[:], offsetsCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading LigatureAttach: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseMark2Array(src []byte, offsetsCount int) (Mark2Array, int, error) {
	var item Mark2Array
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading Mark2Array: "+"EOF: expected length: 2, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[0:]))
		n += 2

		offset := 2
		for i := 0; i < arrayLength; i++ {
			elem, read, err := parseAnchorOffsets(src[offset:], offsetsCount)
			if err != nil {
				return item, 0, fmt.Errorf("reading Mark2Array: %s", err)
			}
			item.mark2Records = append(item.mark2Records, elem)
			offset += read
		}
		n = offset
	}
	{

		read, err := item.customParseMark2Anchors(src[:], offsetsCount)
		if err != nil {
			return item, 0, fmt.Errorf("reading Mark2Array: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseMarkArray(src []byte) (MarkArray, int, error) {
	var item MarkArray
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading MarkArray: "+"EOF: expected length: 2, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[0:]))
		n += 2

		if L := len(src); L < 2+arrayLength*4 {
			return item, 0, fmt.Errorf("reading MarkArray: "+"EOF: expected length: %d, got %d", 2+arrayLength*4, L)
		}

		item.MarkRecords = make([]MarkRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.MarkRecords {
			item.MarkRecords[i].mustParse(src[2+i*4:])
		}
		n += arrayLength * 4
	}
	{

		read, err := item.customParseMarkAnchors(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkArray: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseMarkBasePos(src []byte) (MarkBasePos, int, error) {
	var item MarkBasePos
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: 2, got %d", L)
		}
		item.posFormat = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.markCoverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkBasePos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 6 {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: 6, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[4:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.baseCoverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkBasePos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 8 {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: 8, got %d", L)
		}
		item.markClassCount = binary.BigEndian.Uint16(src[6:])
		n += 2
	}
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: 10, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[8:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.markArray, read, err = ParseMarkArray(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkBasePos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 12 {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: 12, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[10:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkBasePos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.baseArray, read, err = ParseBaseArray(src[offset:], int(item.markClassCount))
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkBasePos: %s", err)
		}
		offset += read
	}
	return item, n, nil
}

func ParseMarkLigPos(src []byte) (MarkLigPos, int, error) {
	var item MarkLigPos
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: 2, got %d", L)
		}
		item.posFormat = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.MarkCoverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkLigPos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 6 {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: 6, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[4:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.LigatureCoverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkLigPos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 8 {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: 8, got %d", L)
		}
		item.MarkClassCount = binary.BigEndian.Uint16(src[6:])
		n += 2
	}
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: 10, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[8:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.MarkArray, read, err = ParseMarkArray(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkLigPos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 12 {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: 12, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[10:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkLigPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.LigatureArray, read, err = ParseLigatureArray(src[offset:], int(item.MarkClassCount))
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkLigPos: %s", err)
		}
		offset += read
	}
	return item, n, nil
}

func ParseMarkMarkPos(src []byte) (MarkMarkPos, int, error) {
	var item MarkMarkPos
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: 2, got %d", L)
		}
		item.PosFormat = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Mark1Coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkMarkPos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 6 {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: 6, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[4:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Mark2Coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkMarkPos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 8 {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: 8, got %d", L)
		}
		item.MarkClassCount = binary.BigEndian.Uint16(src[6:])
		n += 2
	}
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: 10, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[8:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Mark1Array, read, err = ParseMarkArray(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkMarkPos: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 12 {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: 12, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[10:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading MarkMarkPos: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Mark2Array, read, err = ParseMark2Array(src[offset:], int(item.MarkClassCount))
		if err != nil {
			return item, 0, fmt.Errorf("reading MarkMarkPos: %s", err)
		}
		offset += read
	}
	return item, n, nil
}

func ParseMarkRecord(src []byte) (MarkRecord, int, error) {
	var item MarkRecord
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading MarkRecord: "+"EOF: expected length: 4, got %d", L)
	}
	item.mustParse(src)
	n += 4
	return item, n, nil
}

func ParsePairPos(src []byte, valueFormat1 ValueFormat, valueFormat2 ValueFormat) (PairPos, int, error) {
	var item PairPos
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading PairPos: "+"EOF: expected length: 2, got %d", L)
		}
		item.posFormat = pairPosVersion(binary.BigEndian.Uint16(src[0:]))
		n += 2
	}
	{
		var (
			read int
			err  error
		)
		switch item.posFormat {
		case pairPosVersion1:
			item.Data, read, err = ParsePairPosData1(src[:], valueFormat1, valueFormat2)
		case pairPosVersion2:
			item.Data, read, err = ParsePairPosData2(src[:])
		default:
			err = fmt.Errorf("unsupported PairPosDataVersion %d", item.posFormat)
		}
		if err != nil {
			return item, 0, fmt.Errorf("reading PairPos: %s", err)
		}
		n = read
	}
	return item, n, nil
}

// the actual data starts at src[2:]
func ParsePairPosData1(src []byte, valueFormat1 ValueFormat, valueFormat2 ValueFormat) (PairPosData1, int, error) {
	var item PairPosData1
	n := 2
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading PairPosData1: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading PairPosData1: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading PairPosData1: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 8 {
			return item, 0, fmt.Errorf("reading PairPosData1: "+"EOF: expected length: 8, got %d", L)
		}
		_ = src[7] // early bound checking
		item.valueFormat1 = ValueFormat(binary.BigEndian.Uint16(src[4:]))
		item.valueFormat2 = ValueFormat(binary.BigEndian.Uint16(src[6:]))
		n += 4
	}
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading PairPosData1: "+"EOF: expected length: 10, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint16(src[8:]))
		n += 2

		if L := len(src); L < 10+arrayLength*2 {
			return item, 0, fmt.Errorf("reading PairPosData1: "+"EOF: expected length: %d, got %d", 10+arrayLength*2, L)
		}

		item.PairSetOffset = make([]PairSet, arrayLength) // allocation guarded by the previous check
		for i := range item.PairSetOffset {
			offset := int(binary.BigEndian.Uint16(src[10+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading PairPosData1: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.PairSetOffset[i], _, err = ParsePairSet(src[offset:], ValueFormat(item.valueFormat1), ValueFormat(item.valueFormat2))
			if err != nil {
				return item, 0, fmt.Errorf("reading PairPosData1: %s", err)
			}

		}
		n += arrayLength * 2
	}
	return item, n, nil
}

// the actual data starts at src[2:]
func ParsePairPosData2(src []byte) (PairPosData2, int, error) {
	var item PairPosData2
	n := 2
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading PairPosData2: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 8 {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: 8, got %d", L)
		}
		_ = src[7] // early bound checking
		item.valueFormat1 = ValueFormat(binary.BigEndian.Uint16(src[4:]))
		item.valueFormat2 = ValueFormat(binary.BigEndian.Uint16(src[6:]))
		n += 4
	}
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: 10, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[8:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.classDef1, read, err = ParseClassDef(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading PairPosData2: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 12 {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: 12, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[10:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.classDef2, read, err = ParseClassDef(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading PairPosData2: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 16 {
			return item, 0, fmt.Errorf("reading PairPosData2: "+"EOF: expected length: 16, got %d", L)
		}
		_ = src[15] // early bound checking
		item.class1Count = binary.BigEndian.Uint16(src[12:])
		item.class2Count = binary.BigEndian.Uint16(src[14:])
		n += 4
	}
	{

		read, err := item.customParseClass1Records(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading PairPosData2: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParsePairSet(src []byte, valueFormat1 ValueFormat, valueFormat2 ValueFormat) (PairSet, int, error) {
	var item PairSet
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading PairSet: "+"EOF: expected length: 2, got %d", L)
		}
		item.pairValueCount = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{

		read, err := item.customParsePairValueRecords(src[:], valueFormat1, valueFormat2)
		if err != nil {
			return item, 0, fmt.Errorf("reading PairSet: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseRangeRecord(src []byte) (RangeRecord, int, error) {
	var item RangeRecord
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading RangeRecord: "+"EOF: expected length: 6, got %d", L)
	}
	item.mustParse(src)
	n += 6
	return item, n, nil
}

func ParseSinglePos(src []byte) (SinglePos, int, error) {
	var item SinglePos
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading SinglePos: "+"EOF: expected length: 2, got %d", L)
		}
		item.posFormat = singlePosVersion(binary.BigEndian.Uint16(src[0:]))
		n += 2
	}
	{
		var (
			read int
			err  error
		)
		switch item.posFormat {
		case singlePosVersion1:
			item.Data, read, err = ParseSinglePosData1(src[:])
		case singlePosVersion2:
			item.Data, read, err = ParseSinglePosData2(src[:])
		default:
			err = fmt.Errorf("unsupported SinglePosDataVersion %d", item.posFormat)
		}
		if err != nil {
			return item, 0, fmt.Errorf("reading SinglePos: %s", err)
		}
		n = read
	}
	return item, n, nil
}

// the actual data starts at src[2:]
func ParseSinglePosData1(src []byte) (SinglePosData1, int, error) {
	var item SinglePosData1
	n := 2
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading SinglePosData1: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading SinglePosData1: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading SinglePosData1: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 6 {
			return item, 0, fmt.Errorf("reading SinglePosData1: "+"EOF: expected length: 6, got %d", L)
		}
		item.valueFormat = ValueFormat(binary.BigEndian.Uint16(src[4:]))
		n += 2
	}
	{

		read, err := item.customParseValueRecord(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading SinglePosData1: %s", err)
		}
		n = read
	}
	return item, n, nil
}

// the actual data starts at src[2:]
func ParseSinglePosData2(src []byte) (SinglePosData2, int, error) {
	var item SinglePosData2
	n := 2
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading SinglePosData2: "+"EOF: expected length: 4, got %d", L)
		}
		offset := int(binary.BigEndian.Uint16(src[2:]))
		n += 2
		if L := len(src); L < offset {
			return item, 0, fmt.Errorf("reading SinglePosData2: "+"EOF: expected length: %d, got %d", offset, L)
		}

		var (
			err  error
			read int
		)
		item.Coverage, read, err = ParseCoverage(src[offset:])
		if err != nil {
			return item, 0, fmt.Errorf("reading SinglePosData2: %s", err)
		}
		offset += read
	}
	{
		if L := len(src); L < 8 {
			return item, 0, fmt.Errorf("reading SinglePosData2: "+"EOF: expected length: 8, got %d", L)
		}
		_ = src[7] // early bound checking
		item.valueFormat = ValueFormat(binary.BigEndian.Uint16(src[4:]))
		item.valueCount = binary.BigEndian.Uint16(src[6:])
		n += 4
	}
	{

		read, err := item.customParseValueRecords(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading SinglePosData2: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func (item *RangeRecord) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.StartGlyphID = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.EndGlyphID = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.StartCoverageIndex = binary.BigEndian.Uint16(src[4:])
}

func (item *entryExitRecord) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.entryAnchorOffset = Offset16(binary.BigEndian.Uint16(src[0:]))
	item.exitAnchorOffset = Offset16(binary.BigEndian.Uint16(src[2:]))
}

func parseAnchorOffsets(src []byte, offsetsCount int) (anchorOffsets, int, error) {
	var item anchorOffsets
	n := 0
	{

		if L := len(src); L < offsetsCount*2 {
			return item, 0, fmt.Errorf("reading anchorOffsets: "+"EOF: expected length: %d, got %d", offsetsCount*2, L)
		}

		item.offsets = make([]Offset16, offsetsCount) // allocation guarded by the previous check
		for i := range item.offsets {
			item.offsets[i] = Offset16(binary.BigEndian.Uint16(src[i*2:]))
		}
		n += offsetsCount * 2
	}
	return item, n, nil
}
