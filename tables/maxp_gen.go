package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from maxp_src.go. DO NOT EDIT

func ParseMaxp(src []byte) (Maxp, int, error) {
	var item Maxp
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading Maxp: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.version = maxpVersion(binary.BigEndian.Uint32(src[0:]))
	item.NumGlyphs = binary.BigEndian.Uint16(src[4:])
	n += 6

	{
		var (
			read int
			err  error
		)
		switch item.version {
		case maxpVersion05:
			item.data, read, err = parseMaxpData05(src[6:])
		case maxpVersion1:
			item.data, read, err = parseMaxpData1(src[6:])
		default:
			err = fmt.Errorf("unsupported maxpDataVersion %d", item.version)
		}
		if err != nil {
			return item, 0, fmt.Errorf("reading Maxp: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func (item *maxpData1) mustParse(src []byte) {
	item.rawData[0] = binary.BigEndian.Uint16(src[0:])
	item.rawData[1] = binary.BigEndian.Uint16(src[2:])
	item.rawData[2] = binary.BigEndian.Uint16(src[4:])
	item.rawData[3] = binary.BigEndian.Uint16(src[6:])
	item.rawData[4] = binary.BigEndian.Uint16(src[8:])
	item.rawData[5] = binary.BigEndian.Uint16(src[10:])
	item.rawData[6] = binary.BigEndian.Uint16(src[12:])
	item.rawData[7] = binary.BigEndian.Uint16(src[14:])
	item.rawData[8] = binary.BigEndian.Uint16(src[16:])
	item.rawData[9] = binary.BigEndian.Uint16(src[18:])
	item.rawData[10] = binary.BigEndian.Uint16(src[20:])
	item.rawData[11] = binary.BigEndian.Uint16(src[22:])
	item.rawData[12] = binary.BigEndian.Uint16(src[24:])
}

func parseMaxpData05([]byte) (maxpData05, int, error) {
	var item maxpData05
	n := 0
	return item, n, nil
}

func parseMaxpData1(src []byte) (maxpData1, int, error) {
	var item maxpData1
	n := 0
	if L := len(src); L < 26 {
		return item, 0, fmt.Errorf("reading maxpData1: "+"EOF: expected length: 26, got %d", L)
	}
	item.mustParse(src)
	n += 26
	return item, n, nil
}
