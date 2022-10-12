package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from os2.go. DO NOT EDIT

func ParseOs2(src []byte) (Os2, int, error) {
	var item Os2
	n := 0
	{
		subSlice := src[n:]
		if L := len(subSlice); L < 78 {
			return Os2{}, 0, fmt.Errorf("EOF: expected length: 78, got %d", L)
		}

		_ = subSlice[77] // early bound checking
		item.version = binary.BigEndian.Uint16(subSlice[0:])
		item.xAvgCharWidth = binary.BigEndian.Uint16(subSlice[2:])
		item.uSWeightClass = binary.BigEndian.Uint16(subSlice[4:])
		item.uSWidthClass = binary.BigEndian.Uint16(subSlice[6:])
		item.fSType = binary.BigEndian.Uint16(subSlice[8:])
		item.ySubscriptXSize = int16(binary.BigEndian.Uint16(subSlice[10:]))
		item.ySubscriptYSize = int16(binary.BigEndian.Uint16(subSlice[12:]))
		item.ySubscriptXOffset = int16(binary.BigEndian.Uint16(subSlice[14:]))
		item.ySubscriptYOffset = int16(binary.BigEndian.Uint16(subSlice[16:]))
		item.ySuperscriptXSize = int16(binary.BigEndian.Uint16(subSlice[18:]))
		item.ySuperscriptYSize = int16(binary.BigEndian.Uint16(subSlice[20:]))
		item.ySuperscriptXOffset = int16(binary.BigEndian.Uint16(subSlice[22:]))
		item.ySuperscriptYOffset = int16(binary.BigEndian.Uint16(subSlice[24:]))
		item.yStrikeoutSize = int16(binary.BigEndian.Uint16(subSlice[26:]))
		item.yStrikeoutPosition = int16(binary.BigEndian.Uint16(subSlice[28:]))
		item.sFamilyClass = int16(binary.BigEndian.Uint16(subSlice[30:]))
		for i := range item.panose {
			item.panose[i] = subSlice[32+i]
		}
		for i := range item.ulCharRange {
			item.ulCharRange[i] = binary.BigEndian.Uint32(subSlice[42+i*4:])
		}
		item.achVendID = Tag(binary.BigEndian.Uint32(subSlice[58:]))
		item.fsSelection = binary.BigEndian.Uint16(subSlice[62:])
		item.uSFirstCharIndex = binary.BigEndian.Uint16(subSlice[64:])
		item.uSLastCharIndex = binary.BigEndian.Uint16(subSlice[66:])
		item.sTypoAscender = int16(binary.BigEndian.Uint16(subSlice[68:]))
		item.sTypoDescender = int16(binary.BigEndian.Uint16(subSlice[70:]))
		item.sTypoLineGap = int16(binary.BigEndian.Uint16(subSlice[72:]))
		item.usWinAscent = binary.BigEndian.Uint16(subSlice[74:])
		item.usWinDescent = binary.BigEndian.Uint16(subSlice[76:])
		n += 78

	}
	item.higherVersionData = src[n:]
	n = len(src)

	return item, n, nil
}
