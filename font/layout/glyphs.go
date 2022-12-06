package layout

import (
	"github.com/benoitkugler/go-opentype/font"
	"github.com/benoitkugler/go-opentype/tables"
)

type contourPoint struct {
	font.SegmentPoint

	isOnCurve  bool
	isEndPoint bool // this point is the last of the current contour
	isExplicit bool // this point is referenced, i.e., explicit deltas specified */
}

func (c *contourPoint) translate(x, y float32) {
	c.X += x
	c.Y += y
}

func (c *contourPoint) transform(matrix [4]float32) {
	px := c.X*matrix[0] + c.Y*matrix[2]
	c.Y = c.X*matrix[1] + c.Y*matrix[3]
	c.X = px
}

const (
	phantomLeft = iota
	phantomRight
	phantomTop
	phantomBottom
	phantomCount
)

const maxCompositeNesting = 20 // protect against malicious fonts

// use the `glyf` table to fetch the contour points,
// applying variation if needed.
// for composite, recursively calls itself; allPoints includes phantom points and will be at least of length 4
func (f *Font) getPointsForGlyph(gid tables.GlyphID, varCoords []float32, currentDepth int, allPoints *[]contourPoint /* OUT */) {
	// adapted from harfbuzz/src/hb-ot-glyf-table.hh

	if currentDepth > maxCompositeNesting || int(gid) >= len(f.Glyf) {
		return
	}
	g := f.Glyf[gid]

	var points []contourPoint
	if data, ok := g.Data.(tables.SimpleGlyph); ok {
		points = getContourPoints(data) // fetch the "real" points
	} else { // zeros values are enough
		points = make([]contourPoint, pointNumbersCount(g))
	}

	// init phantom point
	points = append(points, make([]contourPoint, phantomCount)...)
	phantoms := points[len(points)-phantomCount:]

	hDelta := float32(g.XMin - getSideBearing(gid, f.hmtx))
	vOrig := float32(g.YMax + getSideBearing(gid, f.vmtx))
	hAdv := float32(f.getBaseAdvance(gid, f.hmtx))
	vAdv := float32(f.getBaseAdvance(gid, f.vmtx))
	phantoms[phantomLeft].X = hDelta
	phantoms[phantomRight].X = hAdv + hDelta
	phantoms[phantomTop].Y = vOrig
	phantoms[phantomBottom].Y = vOrig - vAdv

	if f.isVar() {
		f.gvar.applyDeltasToPoints(gid, varCoords, points)
	}

	switch data := g.Data.(type) {
	case tables.SimpleGlyph:
		*allPoints = append(*allPoints, points...)
	case tables.CompositeGlyph:
		for compIndex, item := range data.Glyphs {
			// recurse on component
			var compPoints []contourPoint

			f.getPointsForGlyph(item.GlyphIndex, varCoords, currentDepth+1, &compPoints)

			LC := len(compPoints)
			if LC < phantomCount { // in case of max depth reached
				return
			}

			/* Copy phantom points from component if USE_MY_METRICS flag set */
			if item.HasUseMyMetrics() {
				copy(phantoms, compPoints[LC-phantomCount:])
			}

			/* Apply component transformation & translation */
			transformPoints(&item, compPoints)

			/* Apply translation from gvar */
			tx, ty := points[compIndex].X, points[compIndex].Y
			for i := range compPoints {
				compPoints[i].translate(tx, ty)
			}

			if item.IsAnchored() {
				p1, p2 := item.ArgsAsIndices()
				if p1 < len(*allPoints) && p2 < LC {
					tx, ty := (*allPoints)[p1].X-compPoints[p2].X, (*allPoints)[p1].Y-compPoints[p2].Y
					for i := range compPoints {
						compPoints[i].translate(tx, ty)
					}
				}
			}

			*allPoints = append(*allPoints, compPoints[0:LC-phantomCount]...)
		}

		*allPoints = append(*allPoints, phantoms...)
	default: // no data for the glyph
		*allPoints = append(*allPoints, phantoms...)
	}

	// apply at top level
	if currentDepth == 0 {
		/* Undocumented rasterizer behavior:
		 * Shift points horizontally by the updated left side bearing */
		tx := -phantoms[phantomLeft].X
		for i := range *allPoints {
			(*allPoints)[i].translate(tx, 0)
		}
	}
}

// does not includes phantom points
func pointNumbersCount(g tables.Glyph) int {
	switch g := g.Data.(type) {
	case tables.SimpleGlyph:
		return len(g.Points)
	case tables.CompositeGlyph:
		/* pseudo component points for each component in composite glyph */
		return len(g.Glyphs)
	}
	return 0
}

// return all the contour points, without phantoms
func getContourPoints(sg tables.SimpleGlyph) []contourPoint {
	const flagOnCurve = 1 << 0 // 0x0001

	points := make([]contourPoint, len(sg.Points))
	for _, end := range sg.EndPtsOfContours {
		points[end].isEndPoint = true
	}
	for i, p := range sg.Points {
		points[i].X, points[i].Y = float32(p.X), float32(p.Y)
		points[i].isOnCurve = p.Flag&flagOnCurve != 0
	}
	return points
}

func extentsFromPoints(allPoints []contourPoint) (ext font.GlyphExtents) {
	truePoints := allPoints[:len(allPoints)-phantomCount]
	if len(truePoints) == 0 {
		// zero extent for the empty glyph
		return ext
	}
	minX, minY := truePoints[0].X, truePoints[0].Y
	maxX, maxY := minX, minY
	for _, p := range truePoints {
		minX = minF(minX, p.X)
		minY = minF(minY, p.Y)
		maxX = maxF(maxX, p.X)
		maxY = maxF(maxY, p.Y)
	}
	ext.XBearing = minX
	ext.YBearing = maxY
	ext.Width = maxX - minX
	ext.Height = minY - maxY
	return ext
}

// walk through the contour points of the given glyph to compute its extends and its phantom points
// As an optimization, if `computeExtents` is false, the extents computation is skipped (a zero value is returned).
func (f *Font) getGlyfPoints(gid tables.GlyphID, varCoords []float32, computeExtents bool) (ext font.GlyphExtents, ph [phantomCount]contourPoint) {
	if int(gid) >= len(f.Glyf) {
		return
	}
	var allPoints []contourPoint
	f.getPointsForGlyph(gid, varCoords, 0, &allPoints)

	copy(ph[:], allPoints[len(allPoints)-phantomCount:])

	if computeExtents {
		ext = extentsFromPoints(allPoints)
	}

	return ext, ph
}

func minF(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxF(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func transformPoints(c *tables.CompositeGlyphPart, points []contourPoint) {
	var transX, transY float32
	if !c.IsAnchored() {
		arg1, arg2 := c.ArgsAsTranslation()
		transX, transY = float32(arg1), float32(arg2)
	}

	scale := c.Scale
	// shortcut identity transform
	if transX == 0 && transY == 0 && scale == [4]float32{1, 0, 0, 1} {
		return
	}

	if c.IsScaledOffsets() {
		for i := range points {
			points[i].translate(transX, transY)
			points[i].transform(scale)
		}
	} else {
		for i := range points {
			points[i].transform(scale)
			points[i].translate(transX, transY)
		}
	}
}
