package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from name_src.go. DO NOT EDIT

func ParseName(src []byte) (Name, int, error) {
	var item Name
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading Name: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	item.count = binary.BigEndian.Uint16(src[2:])
	offsetStringData := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if offsetStringData != 0 { // ignore null offset
			if L := len(src); L < offsetStringData {
				return item, 0, fmt.Errorf("reading Name: "+"EOF: expected length: %d, got %d", offsetStringData, L)
			}

			item.stringData = src[offsetStringData:]
		}
	}
	{
		arrayLength := int(item.count)

		if L := len(src); L < 6+arrayLength*12 {
			return item, 0, fmt.Errorf("reading Name: "+"EOF: expected length: %d, got %d", 6+arrayLength*12, L)
		}

		item.nameRecords = make([]nameRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.nameRecords {
			item.nameRecords[i].mustParse(src[6+i*12:])
		}
		n += arrayLength * 12
	}
	return item, n, nil
}

func (item *nameRecord) mustParse(src []byte) {
	_ = src[11] // early bound checking
	item.platformID = PlatformID(binary.BigEndian.Uint16(src[0:]))
	item.encodingID = EncodingID(binary.BigEndian.Uint16(src[2:]))
	item.languageID = LanguageID(binary.BigEndian.Uint16(src[4:]))
	item.nameID = NameID(binary.BigEndian.Uint16(src[6:]))
	item.length = binary.BigEndian.Uint16(src[8:])
	item.stringOffset = binary.BigEndian.Uint16(src[10:])
}
