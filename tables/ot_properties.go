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

func (cl Coverage1) Len() int { return len(cl.Glyphs) }

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

func (cr Coverage2) Len() int {
	size := 0
	for _, r := range cr.Ranges {
		size += int(r.EndGlyphID - r.StartGlyphID + 1)
	}
	return size
}

func (cl ClassDef1) Class(gi GlyphID) (uint16, bool) {
	if gi < cl.StartGlyphID || gi >= cl.StartGlyphID+GlyphID(len(cl.ClassValueArray)) {
		return 0, false
	}
	return cl.ClassValueArray[gi-cl.StartGlyphID], true
}

func (cl ClassDef1) Extent() int {
	max := uint16(0)
	for _, cid := range cl.ClassValueArray {
		if cid >= max {
			max = cid
		}
	}
	return int(max) + 1
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

func (cl ClassDef2) Extent() int {
	max := uint16(0)
	for _, r := range cl.ClassRangeRecords {
		if r.Class >= max {
			max = r.Class
		}
	}
	return int(max) + 1
}

// ------------------------------------ layout getters ------------------------------------

// FindLanguage looks for [language] and return its index into the [LangSys] slice,
// or -1 if the tag is not found.
func (sc Script) FindLanguage(language Tag) int {
	// LangSys is sorted: binary search
	low, high := 0, len(sc.langSysRecords)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := sc.langSysRecords[mid].Tag
		if language < p {
			high = mid
		} else if language > p {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

// GetLangSys return the language at [index]. It [index] is out of range (for example with 0xFFFF),
// it returns [DefaultLangSys] (which may be empty)
func (sc Script) GetLangSys(index uint16) LangSys {
	if int(index) >= len(sc.LangSys) {
		if sc.DefaultLangSys != nil {
			return *sc.DefaultLangSys
		}
		return LangSys{}
	}
	return sc.LangSys[index]
}

// FindScript looks for [script] and return its index into the Scripts slice,
// or -1 if the tag is not found.
func (la *Layout) FindScript(script Tag) int {
	// Scripts is sorted: binary search
	low, high := 0, len(la.scriptList.records)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := la.scriptList.records[mid].Tag
		if script < p {
			high = mid
		} else if script > p {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

// FindVariationIndex returns the first feature variation matching
// the specified variation coordinates, as an index in the
// `FeatureVariations` field.
// It returns `-1` if not found.
func (la *Layout) FindVariationIndex(coords []float32) int {
	for i, record := range la.featureVariations.featureVariationRecords {
		if record.evaluate(coords) {
			return i
		}
	}
	return -1
}

// returns `true` if the feature is concerned by the `coords`
func (fv FeatureVariationRecord) evaluate(coords []float32) bool {
	for _, c := range fv.conditionSet.conditions {
		if !c.evaluate(coords) {
			return false
		}
	}
	return true
}

// returns `true` if `coords` match the condition `c`
func (c ConditionFormat1) evaluate(coords []float32) bool {
	var coord float32
	if int(c.axisIndex) < len(coords) {
		coord = coords[c.axisIndex]
	}
	return c.filterRangeMinValue <= coord && coord <= c.filterRangeMaxValue
}

// FindFeatureIndex fetches the index for a given feature tag in the GSUB or GPOS table.
// Returns false if not found
func (la *Layout) FindFeatureIndex(featureTag Tag) (uint16, bool) {
	for i, feat := range la.featureList.records { // i fits in uint16
		if featureTag == feat.Tag {
			return uint16(i), true
		}
	}
	return 0, false
}
