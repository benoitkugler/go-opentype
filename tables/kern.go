package tables

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Kern is the kern table. It has multiple header format, defined in Apple AAT and Microsoft OT
// specs, but the subtable data actually are the same.
//
// Microsoft (OT) format
//
//	version uint16 : Table version number (0)
//	nTables uint16 : Number of subtables in the kerning table.
//
// Apple (AAT) old format
//
//	version uint16 : The version number of the kerning table (0x0001 for the current version).
//	nTables uint16 : The number of subtables included in the kerning table.
//
// Apple (AAT) new format
//
//	version uint32 : The version number of the kerning table (0x00010000 for the current version).
//	nTables uint32 : The number of subtables included in the kerning table.
//
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6kern.html
// and - https://learn.microsoft.com/fr-fr/typography/opentype/spec/kern
type Kern struct {
	version uint16
	Tables  []KernSubtable
}

// We apply the following logic:
//   - read the first uint16 -> it's always the major version
//   - if it's 0, we have a Miscrosoft table
//   - if it's 1, we have an Apple table. We read the next uint16,
//     to differentiate between the old and the new Apple format.
func ParseKern(src []byte) (Kern, int, error) {
	if L := len(src); L < 4 {
		return Kern{}, 0, fmt.Errorf("reading Kern: "+"EOF: expected length: 4, got %d", L)
	}

	var numTables uint32

	major := binary.BigEndian.Uint16(src)
	switch major {
	case 0:
		numTables = uint32(binary.BigEndian.Uint16(src[2:]))
		src = src[4:]
	case 1:
		nextUint16 := binary.BigEndian.Uint16(src[2:])
		if nextUint16 == 0 {
			// either new format or old format with 0 subtables, the later being invalid (or at least useless)
			if len(src) < 8 {
				return Kern{}, 0, errors.New("invalid kern table version 1 (EOF)")
			}
			numTables = binary.BigEndian.Uint32(src[4:])
			src = src[8:]
		} else {
			// old format
			numTables = uint32(nextUint16)
			src = src[4:]
		}

	default:
		return Kern{}, 0, fmt.Errorf("unsupported kern table version: %d", major)
	}

	out := make([]KernSubtable, numTables)
	var (
		err    error
		nbRead int
		isOT   = major == 0
	)
	for i := range out {
		if L := len(src); L < nbRead {
			return Kern{}, 0, fmt.Errorf("reading Kern: "+"EOF: expected length: %d, got %d", nbRead, L)
		}
		src = src[nbRead:]
		if isOT {
			out[i], nbRead, err = ParseOTKernSubtableHeader(src)
		} else {
			out[i], nbRead, err = ParseAATKernSubtableHeader(src)
		}
		if err != nil {
			return Kern{}, 0, err
		}
	}

	return Kern{
		version: major,
		Tables:  out,
	}, 0, nil
}

// kernx coverage flags
const (
	kerxBackwards   = 1 << (12 - 8)
	kerxVariation   = 1 << (13 - 8)
	kerxCrossStream = 1 << (14 - 8)
	kerxVertical    = 1 << (15 - 8)
)

// IsHorizontal returns true if the subtable has horizontal kerning values.
func (k AATKernSubtableHeader) IsHorizontal() bool {
	return k.coverage&kerxVertical == 0
}

// IsBackwards returns true if state-table based should process the glyphs backwards.
func (k AATKernSubtableHeader) IsBackwards() bool {
	return k.coverage&kerxBackwards != 0
}

// IsCrossStream returns true if the subtable has cross-stream kerning values.
func (k AATKernSubtableHeader) IsCrossStream() bool {
	return k.coverage&kerxCrossStream != 0
}

// IsVariation returns true if the subtable has variation kerning values.
func (k AATKernSubtableHeader) IsVariation() bool {
	return k.coverage&kerxVariation != 0
}

func (k AATKernSubtableHeader) Data() KernData { return k.data }

// OT kern flags
const (
	kernHorizontal  = 0x01
	kernCrossStream = 0x04
)

// IsHorizontal returns true if the subtable has horizontal kerning values.
func (k OTKernSubtableHeader) IsHorizontal() bool {
	return k.coverage&kernHorizontal == 1
}

// IsBackwards returns true if state-table based should process the glyphs backwards.
func (k OTKernSubtableHeader) IsBackwards() bool { return false }

// IsCrossStream returns true if the subtable has cross-stream kerning values.
func (k OTKernSubtableHeader) IsCrossStream() bool {
	return k.coverage&kernCrossStream != 0
}

// IsVariation returns true if the subtable has variation kerning values.
func (k OTKernSubtableHeader) IsVariation() bool { return false }

func (k OTKernSubtableHeader) Data() KernData { return k.data }
