package tables

import "sort"

// This file implements routines used to simplify acces to the tables
// data.

func (lk AATLoopkup0) Class(g GlyphID) (uint16, bool) {
	if int(g) >= len(lk.Values) {
		return 0, false
	}
	return lk.Values[g], true
}

func (lk AATLoopkup2) Class(g GlyphID) (uint16, bool) {
	// 'adapted' from golang/x/image/font/sfnt
	c := lk.Records
	num := len(c)
	if num == 0 {
		return 0, false
	}

	// classRange is an array of startGlyphID, endGlyphID and target class ID.
	// Ranges are non-overlapping.
	// E.g. 130, 135, 1   137, 137, 5   etc

	idx := sort.Search(num, func(i int) bool { return g <= c[i].FirstGlyph })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if class := c[idx]; g == c[idx].FirstGlyph {
			return class.Value, true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if class := c[idx]; g >= class.FirstGlyph && g <= class.LastGlyph {
			return class.Value, true
		}
	}

	return 0, false
}

func (lk AATLoopkup4) Class(g GlyphID) (uint16, bool) {
	// binary search
	for i, j := 0, len(lk.Records); i < j; {
		h := i + (j-i)/2
		entry := lk.Records[h]
		if g < entry.FirstGlyph {
			j = h
		} else if entry.LastGlyph < g {
			i = h + 1
		} else {
			return entry.Values[g-entry.FirstGlyph], true
		}
	}
	return 0, false
}

func (lk AATLoopkup6) Class(g GlyphID) (uint16, bool) {
	// binary search
	for i, j := 0, len(lk.Records); i < j; {
		h := i + (j-i)/2
		entry := lk.Records[h]
		if g < entry.Glyph {
			j = h
		} else if entry.Glyph < g {
			i = h + 1
		} else {
			return entry.Value, true
		}
	}
	return 0, false
}

func (lk AATLoopkup8) Class(g GlyphID) (uint16, bool) {
	if g < lk.FirstGlyph || g >= lk.FirstGlyph+GlyphID(len(lk.Values)) {
		return 0, false
	}
	return lk.Values[g-lk.FirstGlyph], true
}

func (lk AATLoopkup10) Class(g GlyphID) (uint16, bool) {
	if g < lk.FirstGlyph || g >= lk.FirstGlyph+GlyphID(len(lk.Values)) {
		return 0, false
	}
	return lk.Values[g-lk.FirstGlyph], true
}
