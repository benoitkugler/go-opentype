package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from maxp.go. DO NOT EDIT

func parseMaxp(src []byte) (maxp, int, error) {
	var item maxp
	n := 0
	{
		subSlice := src[n:]
		if L := len(subSlice); L < 6 {
			return maxp{}, 0, fmt.Errorf("EOF: expected length: 6, got %d", L)
		}

		_ = subSlice[5] // early bound checking
		item.version = maxpVersionKind(binary.BigEndian.Uint32(subSlice[0:]))
		item.numGlyphs = binary.BigEndian.Uint16(subSlice[4:])
		n += 6

	}
	{
		var read int
		var err error
		switch item.version {
		case maxpVersionKind05:
			item.data, read, err = parseMaxpVersion05(src[n:])
		case maxpVersionKind1:
			item.data, read, err = parseMaxpVersion1(src[n:])
		default:
			err = fmt.Errorf("unsupported maxpVersionKind %d", item.version)
		}
		if err != nil {
			return maxp{}, 0, err
		}
		n += read
	}
	return item, n, nil
}

func parseMaxpVersion05([]byte) (maxpVersion05, int, error) {
	var item maxpVersion05
	n := 0
	return item, n, nil
}

func (item *maxpVersion1) mustParse(src []byte) {
	_ = src[25] // early bound checking
	for i := range item.data {
		item.data[i] = binary.BigEndian.Uint16(src[i*2:])
	}
}

func parseMaxpVersion1(src []byte) (maxpVersion1, int, error) {
	var item maxpVersion1
	n := 0
	if L := len(src); L < 26 {
		return maxpVersion1{}, 0, fmt.Errorf("EOF: expected length: 26, got %d", L)
	}

	item.mustParse(src)
	n += 26
	return item, n, nil
}
