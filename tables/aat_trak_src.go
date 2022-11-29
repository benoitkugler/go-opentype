package tables

import (
	"encoding/binary"
	"fmt"
)

// Trak is the tracking table.
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6trak.html
type Trak struct {
	version     uint32    // Version number of the tracking table (0x00010000 for the current version).
	format      uint16    // Format of the tracking table (set to 0).
	Horiz       TrackData `offsetSize:"Offset16"` // Offset from start of tracking table to TrackData for horizontal text (or 0 if none).
	Vert        TrackData `offsetSize:"Offset16"` // Offset from start of tracking table to TrackData for vertical text (or 0 if none).
	reserved    uint16    // Reserved. Set to 0.
	postProcess [0]byte   `isOpaque:""`
}

func (tr *Trak) parsePostProcess(src []byte) (int, error) {
	if err := tr.Horiz.postProcess(src); err != nil {
		return 0, err
	}
	if err := tr.Vert.postProcess(src); err != nil {
		return 0, err
	}
	return len(src), nil
}

func parseTrackSizes(src []byte, offset, count int) ([]int16, error) {
	if L := len(src); L < offset+count*2 {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset+count*2, L)
	}
	out := make([]int16, count)
	for i := range out {
		out[i] = int16(binary.BigEndian.Uint16(src[offset+2*i:]))
	}
	return out, nil
}

type TrackData struct {
	nTracks         uint16            // Number of separate tracks included in this table.
	nSizes          uint16            // Number of point sizes included in this table.
	sizeTableOffset uint32            // Offset from start of the tracking table to the start of the size subtable.
	TrackTable      []TrackTableEntry `arrayCount:"ComputedField-nTracks"` // Array[nTracks] of TrackTableEntry records.
	SizeTable       []Float1616       `isOpaque:""`                        // Array[nSizes] of size values.
}

func (tr *TrackData) postProcess(src []byte) error {
	start := int(tr.sizeTableOffset)
	E := start + int(tr.nSizes)*4
	if L := len(src); L < E {
		return fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}
	tr.SizeTable = make([]float32, tr.nSizes)
	for i := range tr.SizeTable {
		tr.SizeTable[i] = Float1616FromUint(binary.BigEndian.Uint32(src[start+i*4:]))
	}

	for i, rec := range tr.TrackTable {
		var err error
		tr.TrackTable[i].PerSizeTracking, err = parseTrackSizes(src, int(rec.offset), len(tr.SizeTable))
		if err != nil {
			return err
		}
	}
	return nil
}

// the actual parsing is performed in Track.parsePostProcess
func (TrackData) parseSizeTable([]byte) (int, error) { return 0, nil }

type TrackTableEntry struct {
	Track           Float1616 // Track value for this record.
	NameIndex       uint16    // The 'name' table index for this track (a short word or phrase like "loose" or "very tight"). NameIndex has a value greater than 255 and less than 32768.
	offset          uint16    // Offset from start of tracking table to per-size tracking values for this track.
	PerSizeTracking []int16   `isOpaque:""` // in font units, with length len(SizeTable)
}

// the actual parsing is performed in Track.parsePostProcess
func (TrackTableEntry) parsePerSizeTracking([]byte) (int, error) { return 0, nil }
