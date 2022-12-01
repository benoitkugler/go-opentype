package tables

import "sort"

func (c Coverage1) Index(gi GlyphID) (int, bool) {
	num := len(c.Glyphs)
	idx := sort.Search(num, func(i int) bool { return gi <= c.Glyphs[i] })
	if idx < num && c.Glyphs[idx] == gi {
		return idx, true
	}
	return 0, false
}

func (c Coverage2) Index(gi GlyphID) (int, bool) {
	num := len(c.Ranges)
	if num == 0 {
		return 0, false
	}

	idx := sort.Search(num, func(i int) bool { return gi <= c.Ranges[i].StartGlyphID })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if rang := c.Ranges[idx]; gi == rang.StartGlyphID {
			return int(rang.StartCoverageIndex), true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if rang := c.Ranges[idx]; gi >= rang.StartGlyphID && gi <= rang.EndGlyphID {
			return int(rang.StartCoverageIndex) + int(gi-rang.StartGlyphID), true
		}
	}

	return 0, false
}

func (cl ClassDef1) Class(gi GlyphID) (uint16, bool) {
	if gi < cl.StartGlyphID || gi >= cl.StartGlyphID+GlyphID(len(cl.ClassValueArray)) {
		return 0, false
	}
	return cl.ClassValueArray[gi-cl.StartGlyphID], true
}

func (cl ClassDef2) Class(g GlyphID) (uint16, bool) {
	// 'adapted' from golang/x/image/font/sfnt
	c := cl.ClassRangeRecords
	num := len(c)
	if num == 0 {
		return 0, false
	}

	// classRange is an array of startGlyphID, endGlyphID and target class ID.
	// Ranges are non-overlapping.
	// E.g. 130, 135, 1   137, 137, 5   etc

	idx := sort.Search(num, func(i int) bool { return g <= c[i].StartGlyphID })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if class := c[idx]; g == c[idx].StartGlyphID {
			return class.Class, true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if class := c[idx]; g >= class.StartGlyphID && g <= class.EndGlyphID {
			return class.Class, true
		}
	}

	return 0, false
}
