package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_kernx_src.go. DO NOT EDIT

func ParseKerx(src []byte) (Kerx, int, error) {
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

		item.rawData = src[8:]
		n = len(src)
	}
	return item, n, nil
}

func (item *aatSTHeader) mustParse(src []byte) {
	_ = src[7] // early bound checking
	item.stateSize = binary.BigEndian.Uint16(src[0:])
	item.classTable = binary.BigEndian.Uint16(src[2:])
	item.stateArray = binary.BigEndian.Uint16(src[4:])
	item.entryTable = binary.BigEndian.Uint16(src[6:])
}

func (item *kernx0Record) mustParse(src []byte) {
	_ = src[5] // early bound checking
	item.left = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.right = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.value = int16(binary.BigEndian.Uint16(src[4:]))
}

func parseKernSubtable0(src []byte) (kernSubtable0, int, error) {
	var item kernSubtable0
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading kernSubtable0: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.nPairs = binary.BigEndian.Uint16(src[0:])
	item.searchRange = binary.BigEndian.Uint16(src[2:])
	item.entrySelector = binary.BigEndian.Uint16(src[4:])
	item.rangeShift = binary.BigEndian.Uint16(src[6:])
	n += 8

	{
		arrayLength := int(item.nPairs)

		if L := len(src); L < 8+arrayLength*6 {
			return item, 0, fmt.Errorf("reading kernSubtable0: "+"EOF: expected length: %d, got %d", 8+arrayLength*6, L)
		}

		item.pairs = make([]kernx0Record, arrayLength) // allocation guarded by the previous check
		for i := range item.pairs {
			item.pairs[i].mustParse(src[8+i*6:])
		}
		n += arrayLength * 6
	}
	{

		item.rawData = src[n:]
		n = len(src)
	}
	return item, n, nil
}

func parseKernSubtable1(src []byte) (kernSubtable1, int, error) {
	var item kernSubtable1
	n := 0
	if L := len(src); L < 10 {
		return item, 0, fmt.Errorf("reading kernSubtable1: "+"EOF: expected length: 10, got %d", L)
	}
	_ = src[9] // early bound checking
	item.aatSTHeader.mustParse(src[0:])
	item.valueTable = binary.BigEndian.Uint16(src[8:])
	n += 10

	{

		item.rawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func parseKernSubtable2(src []byte) (kernSubtable2, int, error) {
	var item kernSubtable2
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading kernSubtable2: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.rowWidth = binary.BigEndian.Uint16(src[0:])
	offsetLeft := int(binary.BigEndian.Uint16(src[2:]))
	offsetRight := int(binary.BigEndian.Uint16(src[4:]))
	item.kerningArrayOffset = binary.BigEndian.Uint16(src[6:])
	n += 8

	{

		if offsetLeft != 0 { // ignore null offset
			if L := len(src); L < offsetLeft {
				return item, 0, fmt.Errorf("reading kernSubtable2: "+"EOF: expected length: %d, got %d", offsetLeft, L)
			}

		}
	}
	{

		if offsetRight != 0 { // ignore null offset
			if L := len(src); L < offsetRight {
				return item, 0, fmt.Errorf("reading kernSubtable2: "+"EOF: expected length: %d, got %d", offsetRight, L)
			}

		}
	}
	{

		item.rawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func parseKerxSubtable0(src []byte) (kerxSubtable0, int, error) {
	var item kerxSubtable0
	n := 0
	if L := len(src); L < 16 {
		return item, 0, fmt.Errorf("reading kerxSubtable0: "+"EOF: expected length: 16, got %d", L)
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
			return item, 0, fmt.Errorf("reading kerxSubtable0: "+"EOF: expected length: %d, got %d", 16+arrayLength*6, L)
		}

		item.pairs = make([]kernx0Record, arrayLength) // allocation guarded by the previous check
		for i := range item.pairs {
			item.pairs[i].mustParse(src[16+i*6:])
		}
		n += arrayLength * 6
	}
	{

		item.rawData = src[n:]
		n = len(src)
	}
	return item, n, nil
}

func parseKerxSubtable1([]byte) (kerxSubtable1, int, error) {
	var item kerxSubtable1
	n := 0
	return item, n, nil
}

func parseKerxSubtable2([]byte) (kerxSubtable2, int, error) {
	var item kerxSubtable2
	n := 0
	return item, n, nil
}

func parseKerxSubtable4([]byte) (kerxSubtable4, int, error) {
	var item kerxSubtable4
	n := 0
	return item, n, nil
}

func parseKerxSubtable6([]byte) (kerxSubtable6, int, error) {
	var item kerxSubtable6
	n := 0
	return item, n, nil
}
