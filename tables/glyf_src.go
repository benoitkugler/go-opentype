package tables

import (
	"encoding/binary"
	"fmt"
)

// shared with gvar, sbix, eblc
// return an error only if data is not long enough
func parseLoca(src []byte, numGlyphs int, isLong bool) (out []uint32, err error) {
	var size int
	if isLong {
		size = (numGlyphs + 1) * 4
	} else {
		size = (numGlyphs + 1) * 2
	}
	if L := len(src); L < size {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", size, L)
	}
	out = make([]uint32, numGlyphs+1)
	if isLong {
		for i := range out {
			out[i] = binary.BigEndian.Uint32(src[4*i:])
		}
	} else {
		for i := range out {
			out[i] = 2 * uint32(binary.BigEndian.Uint16(src[2*i:])) // The actual local offset divided by 2 is stored.
		}
	}
	return out, nil
}

// Glyph Data
type Glyf []Glyph

// locaOffsets has length numGlyphs + 1
func parseGlyf(src []byte, locaOffsets []uint32) (Glyf, error) {
	out := make(Glyf, len(locaOffsets)-1)
	var err error
	for i := range out {
		start, end := locaOffsets[i], locaOffsets[i+1]
		// If a glyph has no outline, then loca[n] = loca [n+1].
		if start == end {
			continue
		}
		out[i], _, err = ParseGlyph(src[start:end])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

type Glyph struct {
	numberOfContours int16     // If the number of contours is greater than or equal to zero, this is a simple glyph. If negative, this is a composite glyph — the value -1 should be used for composite glyphs.
	XMin             int16     // Minimum x for coordinate data.
	YMin             int16     // Minimum y for coordinate data.
	XMax             int16     // Maximum x for coordinate data.
	YMax             int16     // Maximum y for coordinate data.
	Data             GlyphData `isOpaque:"" subsliceStart:"AtCurrent"`
}

func (gl *Glyph) parseData(src []byte) (read int, err error) {
	if gl.numberOfContours >= 0 { // simple glyph
		gl.Data, read, err = ParseSimpleGlyph(src)
	} else { // composite glyph
		gl.Data, read, err = ParseCompositeGlyph(src)
	}
	return read, err
}

type GlyphData interface {
	isGlyphData()
}

func (SimpleGlyph) isGlyphData()    {}
func (CompositeGlyph) isGlyphData() {}

type SimpleGlyph struct {
	EndPtsOfContours uint16  // [numberOfContours] Array of point indices for the last point of each contour, in increasing numeric order.
	Instructions     []uint8 `arrayCount:"FirstUint16"` // [instructionLength] Array of instruction byte code for the glyph.
	RawData          []byte  `arrayCount:"ToEnd"`
	// flags            uint8   // [variable]	Array of flag elements. See below for details regarding the number of flag array elements.
	// uint8 or int16	xCoordinates[variable]	Contour point x-coordinates. See below for details regarding the number of coordinate array elements. Coordinate for the first point is relative to (0,0); others are relative to previous point.
	// uint8 or int16	yCoordinates[variable]	Contour point y-coordinates. See below for details regarding the number of coordinate array elements. Coordinate for the first point is relative to (0,0); others are relative to previous point.
}

type CompositeGlyph struct {
	RawData []byte `arrayCount:"ToEnd"`
}
