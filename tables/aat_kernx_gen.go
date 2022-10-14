package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_kernx.go. DO NOT EDIT

func ParseKerx(src []byte) (Kerx, int, error) {
	var item Kerx
	n := 0
	{
		subSlice := src[n:]
		if L := len(subSlice); L < 8 {
			return Kerx{}, 0, fmt.Errorf("EOF: expected length: 8, got %d", L)
		}

		_ = subSlice[7] // early bound checking
		item.version = binary.BigEndian.Uint16(subSlice[0:])
		item.padding = binary.BigEndian.Uint16(subSlice[2:])
		item.nTables = binary.BigEndian.Uint32(subSlice[4:])
		n += 8

	}
	item.rawData = src[n:]
	n = len(src)

	return item, n, nil
}
func parseAatLookupTable8(src []byte) (aatLookupTable8, int, error) {
	var item aatLookupTable8
	n := 0
	{
		subSlice := src[n:]
		if L := len(subSlice); L < 2 {
			return aatLookupTable8{}, 0, fmt.Errorf("EOF: expected length: 2, got %d", L)
		}

		_ = subSlice[1] // early bound checking
		item.firstGlyph = glyphID(binary.BigEndian.Uint16(subSlice[0:]))
		n += 2

	}
	{
		subSlice := src[n:]
		if L := len(subSlice); L < 2 {
			return aatLookupTable8{}, 0, fmt.Errorf("EOF: expected length: %d, got %d", 2, L)
		}
		arrayLength := int(binary.BigEndian.Uint16(subSlice[:]))
		if L := len(subSlice); L < 2+arrayLength*2 {
			return aatLookupTable8{}, 0, fmt.Errorf("EOF: expected length: %d, got %d", 2+arrayLength*2, L)
		}

		item.values = make([]uint16, arrayLength) // allocation guarded by the previous check
		for i := range item.values {
			item.values[i] = binary.BigEndian.Uint16(subSlice[2+i*2:])
		}

		n += 2 + arrayLength*2
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
	item.left = glyphID(binary.BigEndian.Uint16(src[0:]))
	item.right = glyphID(binary.BigEndian.Uint16(src[2:]))
	item.value = int16(binary.BigEndian.Uint16(src[4:]))
}
