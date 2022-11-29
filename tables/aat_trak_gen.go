package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_trak_src.go. DO NOT EDIT

func ParseTrackData(src []byte) (TrackData, int, error) {
	var item TrackData
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading TrackData: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.nTracks = binary.BigEndian.Uint16(src[0:])
	item.nSizes = binary.BigEndian.Uint16(src[2:])
	item.sizeTableOffset = binary.BigEndian.Uint32(src[4:])
	n += 8

	{
		arrayLength := int(item.nTracks)

		offset := 8
		for i := 0; i < arrayLength; i++ {
			elem, read, err := ParseTrackTableEntry(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading TrackData: %s", err)
			}
			item.TrackTable = append(item.TrackTable, elem)
			offset += read
		}
		n = offset
	}
	{

		read, err := item.parseSizeTable(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading TrackData: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseTrackTableEntry(src []byte) (TrackTableEntry, int, error) {
	var item TrackTableEntry
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading TrackTableEntry: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.Track = Float1616FromUint(binary.BigEndian.Uint32(src[0:]))
	item.NameIndex = binary.BigEndian.Uint16(src[4:])
	item.offset = binary.BigEndian.Uint16(src[6:])
	n += 8

	{

		read, err := item.parsePerSizeTracking(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading TrackTableEntry: %s", err)
		}
		n = read
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
	offsetItemhoriz := int(binary.BigEndian.Uint16(src[6:]))
	offsetItemvert := int(binary.BigEndian.Uint16(src[8:]))
	item.reserved = binary.BigEndian.Uint16(src[10:])
	n += 12

	{

		if offsetItemhoriz != 0 { // ignore null offset
			if L := len(src); L < offsetItemhoriz {
				return item, 0, fmt.Errorf("reading Trak: "+"EOF: expected length: %d, got %d", offsetItemhoriz, L)
			}

			var (
				err  error
				read int
			)
			item.Horiz, read, err = ParseTrackData(src[offsetItemhoriz:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Trak: %s", err)
			}
			offsetItemhoriz += read

		}
	}
	{

		if offsetItemvert != 0 { // ignore null offset
			if L := len(src); L < offsetItemvert {
				return item, 0, fmt.Errorf("reading Trak: "+"EOF: expected length: %d, got %d", offsetItemvert, L)
			}

			var (
				err  error
				read int
			)
			item.Vert, read, err = ParseTrackData(src[offsetItemvert:])
			if err != nil {
				return item, 0, fmt.Errorf("reading Trak: %s", err)
			}
			offsetItemvert += read

		}
	}
	{

		read, err := item.parsePostProcess(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading Trak: %s", err)
		}
		n = read
	}
	return item, n, nil
}
