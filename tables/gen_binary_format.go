package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from opentype/binary_format.go. DO NOT EDIT

func (item *VarAxis) mustParse(src []byte) {
	_ = src[19] // early bound checking
	item.Tag = Tag(binary.BigEndian.Uint32(src[0:]))
	item.Minimum = Float1616FromUint(binary.BigEndian.Uint32(src[4:]))
	item.Default = Float1616FromUint(binary.BigEndian.Uint32(src[8:]))
	item.Maximum = Float1616FromUint(binary.BigEndian.Uint32(src[12:]))
	item.flags = uint16(binary.BigEndian.Uint16(src[16:]))
	item.strid = NameID(binary.BigEndian.Uint16(src[18:]))
}

func parseVarAxis(src []byte) (VarAxis, int, error) {
	var item VarAxis
	n := 0
	if L := len(src); L < 20 {
		return VarAxis{}, 0, fmt.Errorf("EOF: expected length: 20, got %d", L)
	}

	item.mustParse(src)
	n += 20
	return item, n, nil
}

func (item *fvarHeader) mustParse(src []byte) {
	_ = src[15] // early bound checking
	item.majorVersion = uint16(binary.BigEndian.Uint16(src[0:]))
	item.minorVersion = uint16(binary.BigEndian.Uint16(src[2:]))
	item.axesArrayOffset = uint16(binary.BigEndian.Uint16(src[4:]))
	item.reserved = uint16(binary.BigEndian.Uint16(src[6:]))
	item.axisCount = uint16(binary.BigEndian.Uint16(src[8:]))
	item.axisSize = uint16(binary.BigEndian.Uint16(src[10:]))
	item.instanceCount = uint16(binary.BigEndian.Uint16(src[12:]))
	item.instanceSize = uint16(binary.BigEndian.Uint16(src[14:]))
}

func parseFvarHeader(src []byte) (fvarHeader, int, error) {
	var item fvarHeader
	n := 0
	if L := len(src); L < 16 {
		return fvarHeader{}, 0, fmt.Errorf("EOF: expected length: 16, got %d", L)
	}

	item.mustParse(src)
	n += 16
	return item, n, nil
}

 
func parseMortxSubtable(src []byte, subtableLength int) (mortxSubtable, int, error) {
	var item mortxSubtable
	n := 0
	{
		subSlice := src[n:]
		if L := len(subSlice); L < 12 {
			return mortxSubtable{}, 0, fmt.Errorf("EOF: expected length: 12, got %d", L)
		}

		_ = subSlice[11] // early bound checking
		item.tableLength = uint32(binary.BigEndian.Uint32(subSlice[0:]))
		item.coverage = byte(subSlice[4])
		item.padding = uint16(binary.BigEndian.Uint16(subSlice[5:]))
		item.kind = byte(subSlice[7])
		item.flags = uint32(binary.BigEndian.Uint32(subSlice[8:]))

		n += 12
	}
	{
		subSlice := src[n:]
		if L := len(subSlice); L < +subtableLength {
			return mortxSubtable{}, 0, fmt.Errorf("EOF: expected length: %d, got %d", +subtableLength, L)
		}

		item.subtable = make([]byte, subtableLength) // allocation guarded by the previous check
		for i := range item.subtable {
			item.subtable[i] = byte(subSlice[0+i])
		}

		n += 0 + subtableLength*1
	}
	return item, n, nil
}
