package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_trak_src.go. DO NOT EDIT

func ParseTrackData(src []byte, parentSrc []byte) (TrackData, int, error) {
	var item TrackData
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading TrackData: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.nTracks = binary.BigEndian.Uint16(src[0:])
	item.nSizes = binary.BigEndian.Uint16(src[2:])
	offsetSizeTable := int(binary.BigEndian.Uint32(src[4:]))
	n += 8

	{

		if offsetSizeTable != 0 { // ignore null offset
			if L := len(parentSrc); L < offsetSizeTable {
				return item, 0, fmt.Errorf("reading TrackData: "+"EOF: expected length: %d, got %d", offsetSizeTable, L)
			}

			arrayLength := int(item.nSizes)

			if L := len(parentSrc); L < offsetSizeTable+arrayLength*4 {
				return item, 0, fmt.Errorf("reading TrackData: "+"EOF: expected length: %d, got %d", offsetSizeTable+arrayLength*4, L)
			}

			item.SizeTable = make([]float32, arrayLength) // allocation guarded by the previous check
			for i := range item.SizeTable {
				item.SizeTable[i] = Float1616FromUint(binary.BigEndian.Uint32(parentSrc[offsetSizeTable+i*4:]))
			}
			offsetSizeTable += arrayLength * 4
		}
	}
	{
		arrayLength := int(item.nTracks)

		offset := 8
		for i := 0; i < arrayLength; i++ {
			elem, read, err := ParseTrackTableEntry(src[offset:], parentSrc, int(item.nSizes))
			if err != nil {
				return item, 0, fmt.Errorf("reading TrackData: %s", err)
			}
			item.TrackTable = append(item.TrackTable, elem)
			offset += read
		}
		n = offset
	}
	return item, n, nil
}

func ParseTrackTableEntry(src []byte, grandParentSrc []byte, perSizeTrackingCount int) (TrackTableEntry, int, error) {
	var item TrackTableEntry
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading TrackTableEntry: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.Track = Float1616FromUint(binary.BigEndian.Uint32(src[0:]))
	item.NameIndex = binary.BigEndian.Uint16(src[4:])
	offsetPerSizeTracking := int(binary.BigEndian.Uint16(src[6:]))
	n += 8

	{

		if offsetPerSizeTracking != 0 { // ignore null offset
			if L := len(grandParentSrc); L < offsetPerSizeTracking {
				return item, 0, fmt.Errorf("reading TrackTableEntry: "+"EOF: expected length: %d, got %d", offsetPerSizeTracking, L)
			}

			if L := len(grandParentSrc); L < offsetPerSizeTracking+perSizeTrackingCount*2 {
				return item, 0, fmt.Errorf("reading TrackTableEntry: "+"EOF: expected length: %d, got %d", offsetPerSizeTracking+perSizeTrackingCount*2, L)
			}

			item.PerSizeTracking = make([]int16, perSizeTrackingCount) // allocation guarded by the previous check
			for i := range item.PerSizeTracking {
				item.PerSizeTracking[i] = int16(binary.BigEndian.Uint16(grandParentSrc[offsetPerSizeTracking+i*2:]))
			}
			offsetPerSizeTracking += perSizeTrackingCount * 2
		}
	}
	return item, n, nil
}

func ParseTrak(src []byte) (Trak, int, error) {
	var item Trak
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading Trak: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint32(src[0:])
	item.format = binary.BigEndian.Uint16(src[4:])
	offsetHoriz := int(binary.BigEndian.Uint16(src[6:]))
	offsetVert := int(binary.BigEndian.Uint16(src[8:]))
	item.reserved = binary.BigEndian.Uint16(src[10:])
	n += 12

	{

		if offsetHoriz != 0 { // ignore null offset
			if L := len(src); L < offsetHoriz {
				return item, 0, fmt.Errorf("reading Trak: "+"EOF: expected length: %d, got %d", offsetHoriz, L)
			}

			var err error
			item.Horiz, _, err = ParseTrackData(src[offsetHoriz:], src)
			if err != nil {
				return item, 0, fmt.Errorf("reading Trak: %s", err)
			}

		}
	}
	{

		if offsetVert != 0 { // ignore null offset
			if L := len(src); L < offsetVert {
				return item, 0, fmt.Errorf("reading Trak: "+"EOF: expected length: %d, got %d", offsetVert, L)
			}

			var err error
			item.Vert, _, err = ParseTrackData(src[offsetVert:], src)
			if err != nil {
				return item, 0, fmt.Errorf("reading Trak: %s", err)
			}

		}
	}
	return item, n, nil
}
