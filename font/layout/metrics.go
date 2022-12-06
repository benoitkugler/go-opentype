package layout

import (
	"math"

	"github.com/benoitkugler/go-opentype/font"
	"github.com/benoitkugler/go-opentype/opentype"
	"github.com/benoitkugler/go-opentype/tables"
)

func (f *Font) GetGlyphContourPoint(glyph GID, pointIndex uint16) (x, y int32, ok bool) {
	// harfbuzz seems not to implement this feature
	return 0, 0, false
}

// GlyphName returns the name of the given glyph, or an empty
// string if the glyph is invalid or has no name.
func (f *Font) GlyphName(glyph GID) string {
	if postNames := f.post.names; postNames != nil {
		if name := postNames.glyphName(glyph); name != "" {
			return name
		}
	}
	if f.cff != nil {
		return f.cff.GlyphName(glyph)
	}
	return ""
}

// Upem returns the units per em of the font file.
// This value is only relevant for scalable fonts.
func (f *Font) Upem() uint16 { return f.upem }

var (
	metricsTagHorizontalAscender  = opentype.MustNewTag("hasc")
	metricsTagHorizontalDescender = opentype.MustNewTag("hdsc")
	metricsTagHorizontalLineGap   = opentype.MustNewTag("hlgp")
	metricsTagVerticalAscender    = opentype.MustNewTag("vasc")
	metricsTagVerticalDescender   = opentype.MustNewTag("vdsc")
	metricsTagVerticalLineGap     = opentype.MustNewTag("vlgp")
)

func fixAscenderDescender(value float32, metricsTag Tag) float32 {
	if metricsTag == metricsTagHorizontalAscender || metricsTag == metricsTagVerticalAscender {
		return float32(math.Abs(float64(value)))
	}
	if metricsTag == metricsTagHorizontalDescender || metricsTag == metricsTagVerticalDescender {
		return float32(-math.Abs(float64(value)))
	}
	return value
}

func (f *Font) getPositionCommon(metricTag Tag, varCoords []float32) (float32, bool) {
	deltaVar := f.mvar.getVar(metricTag, varCoords)
	switch metricTag {
	case metricsTagHorizontalAscender:
		if f.os2.useTypoMetrics {
			return fixAscenderDescender(float32(f.os2.sTypoAscender)+deltaVar, metricTag), true
		} else if f.hhea != nil {
			return fixAscenderDescender(float32(f.hhea.Ascender)+deltaVar, metricTag), true
		}

	case metricsTagHorizontalDescender:
		if f.os2.useTypoMetrics {
			return fixAscenderDescender(float32(f.os2.sTypoDescender)+deltaVar, metricTag), true
		} else if f.hhea != nil {
			return fixAscenderDescender(float32(f.hhea.Descender)+deltaVar, metricTag), true
		}
	case metricsTagHorizontalLineGap:
		if f.os2.useTypoMetrics {
			return fixAscenderDescender(float32(f.os2.sTypoLineGap)+deltaVar, metricTag), true
		} else if f.hhea != nil {
			return fixAscenderDescender(float32(f.hhea.LineGap)+deltaVar, metricTag), true
		}
	case metricsTagVerticalAscender:
		if f.vhea != nil {
			return fixAscenderDescender(float32(f.vhea.Ascender)+deltaVar, metricTag), true
		}
	case metricsTagVerticalDescender:
		if f.vhea != nil {
			return fixAscenderDescender(float32(f.vhea.Descender)+deltaVar, metricTag), true
		}
	case metricsTagVerticalLineGap:
		if f.vhea != nil {
			return fixAscenderDescender(float32(f.vhea.LineGap)+deltaVar, metricTag), true
		}
	}
	return 0, false
}

// FontHExtents returns the extents of the font for horizontal text, or false
// it not available, in font units.
func (f *Face) FontHExtents() (font.FontExtents, bool) {
	var (
		out           font.FontExtents
		ok1, ok2, ok3 bool
	)
	out.Ascender, ok1 = f.Font.getPositionCommon(metricsTagHorizontalAscender, f.Coords)
	out.Descender, ok2 = f.Font.getPositionCommon(metricsTagHorizontalDescender, f.Coords)
	out.LineGap, ok3 = f.Font.getPositionCommon(metricsTagHorizontalLineGap, f.Coords)
	return out, ok1 && ok2 && ok3
}

// FontVExtents is the same as `FontHExtents`, but for vertical text.
func (f *Face) FontVExtents() (font.FontExtents, bool) {
	var (
		out           font.FontExtents
		ok1, ok2, ok3 bool
	)
	out.Ascender, ok1 = f.Font.getPositionCommon(metricsTagVerticalAscender, f.Coords)
	out.Descender, ok2 = f.Font.getPositionCommon(metricsTagVerticalDescender, f.Coords)
	out.LineGap, ok3 = f.Font.getPositionCommon(metricsTagVerticalLineGap, f.Coords)
	return out, ok1 && ok2 && ok3
}

var (
	tagStrikeoutSize      = opentype.MustNewTag("strs")
	tagStrikeoutOffset    = opentype.MustNewTag("stro")
	tagUnderlineSize      = opentype.MustNewTag("unds")
	tagUnderlineOffset    = opentype.MustNewTag("undo")
	tagSuperscriptYSize   = opentype.MustNewTag("spys")
	tagSuperscriptYOffset = opentype.MustNewTag("spyo")
	tagSuperscriptXSize   = opentype.MustNewTag("spxs")
	tagSuperscriptXOffset = opentype.MustNewTag("spxo")
	tagSubscriptYSize     = opentype.MustNewTag("sbys")
	tagSubscriptYOffset   = opentype.MustNewTag("sbyo")
	tagSubscriptXOffset   = opentype.MustNewTag("sbxo")
	tagXHeight            = opentype.MustNewTag("xhgt")
	tagCapHeight          = opentype.MustNewTag("cpht")
)

// LineMetric returns the metric identified by `metric` (in fonts units).
func (face *Face) LineMetric(metric font.LineMetric) float32 {
	f := face.Font
	switch metric {
	case font.UnderlinePosition:
		return f.post.underlinePosition + f.mvar.getVar(tagUnderlineOffset, face.Coords)
	case font.UnderlineThickness:
		return f.post.underlineThickness + f.mvar.getVar(tagUnderlineSize, face.Coords)
	case font.StrikethroughPosition:
		return float32(f.os2.yStrikeoutPosition) + f.mvar.getVar(tagStrikeoutOffset, face.Coords)
	case font.StrikethroughThickness:
		return float32(f.os2.yStrikeoutSize) + f.mvar.getVar(tagStrikeoutSize, face.Coords)
	case font.SuperscriptEmYSize:
		return float32(f.os2.ySuperscriptYSize) + f.mvar.getVar(tagSuperscriptYSize, face.Coords)
	case font.SuperscriptEmXOffset:
		return float32(f.os2.ySuperscriptXOffset) + f.mvar.getVar(tagSuperscriptXOffset, face.Coords)
	case font.SubscriptEmYSize:
		return float32(f.os2.ySubscriptYSize) + f.mvar.getVar(tagSubscriptYSize, face.Coords)
	case font.SubscriptEmYOffset:
		return float32(f.os2.ySubscriptYOffset) + f.mvar.getVar(tagSubscriptYOffset, face.Coords)
	case font.SubscriptEmXOffset:
		return float32(f.os2.ySubscriptXOffset) + f.mvar.getVar(tagSubscriptXOffset, face.Coords)
	case font.CapHeight:
		return float32(f.os2.sCapHeight) + f.mvar.getVar(tagCapHeight, face.Coords)
	case font.XHeight:
		return float32(f.os2.sxHeigh) + f.mvar.getVar(tagXHeight, face.Coords)
	default:
		return 0
	}
}

// NominalGlyph returns the glyph used to represent the given rune,
// or false if not found.
// Note that it only looks into the cmap, without taking account substitutions
// nor variation selectors.
func (f *Font) NominalGlyph(ch rune) (GID, bool) { return f.cmap.Lookup(ch) }

// VariationGlyph retrieves the glyph ID for a specified Unicode code point
// followed by a specified Variation Selector code point, or false if not found
func (f *Font) VariationGlyph(ch, varSelector rune) (GID, bool) {
	gid, kind := f.cmapVar.GetGlyphVariant(ch, varSelector)
	switch kind {
	case font.VariantNotFound:
		return 0, false
	case font.VariantFound:
		return gid, true
	default: // VariantUseDefault
		return f.NominalGlyph(ch)
	}
}

// do not take into account variations
func (f *Font) getBaseAdvance(gid tables.GlyphID, table tables.Hmtx) int16 {
	LM, LS := len(table.Metrics), len(table.LeftSideBearings)
	index := int(gid)
	if index < LM {
		return table.Metrics[index].AdvanceWidth
	} else if index < LS+LM { // return the last value
		return table.Metrics[len(table.Metrics)-1].AdvanceWidth
	} else {
		/* If `table` is empty, it means we don't have the metrics table
		 * for this direction: return default advance.  Otherwise, it means that the
		 * glyph index is out of bound: return zero. */
		if LM+LS == 0 {
			return int16(f.upem)
		}
		return 0
	}
}

// return the base side bearing, handling invalid glyph index
func getSideBearing(gid tables.GlyphID, table tables.Hmtx) int16 {
	LM, LS := len(table.Metrics), len(table.LeftSideBearings)
	index := int(gid)
	if index < LM {
		return table.Metrics[index].LeftSideBearing
	} else if index < LS+LM {
		return table.Metrics[index-LM].AdvanceWidth
	} else {
		return 0
	}
}

func clamp(v float32) float32 {
	if v < 0 {
		v = 0
	}
	return v
}

func ceil(v float32) int16 {
	return int16(math.Ceil(float64(v)))
}

func (f *Font) getGlyphAdvanceVar(gid tables.GlyphID, varCoords []float32, isVertical bool) float32 {
	_, phantoms := f.getGlyfPoints(gid, varCoords, false)
	if isVertical {
		return clamp(phantoms[phantomTop].Y - phantoms[phantomBottom].Y)
	}
	return clamp(phantoms[phantomRight].X - phantoms[phantomLeft].X)
}

func (face *Face) HorizontalAdvance(gid GID) float32 {
	f := face.Font
	advance := f.getBaseAdvance(tables.GlyphID(gid), f.hmtx)
	if !f.isVar() {
		return float32(advance)
	}
	if f.hvar != nil {
		return float32(advance) + f.hvar.getAdvanceVar(gid, f.varCoords)
	}
	return f.getGlyphAdvanceVar(gid, false)
}

// return `true` is the font is variable and `varCoords` is valid
func (f *Font) isVar() bool {
	return len(f.varCoords) != 0 && len(f.varCoords) == len(f.fvar.Axis)
}

func (f *Font) VerticalAdvance(gid GID) float32 {
	// return the opposite of the advance from the font
	advance := f.getBaseAdvance(gid, f.vmtx)
	if !f.isVar() {
		return -float32(advance)
	}
	if f.vvar != nil {
		return -float32(advance) - f.vvar.getAdvanceVar(gid, f.varCoords)
	}
	return -f.getGlyphAdvanceVar(gid, true)
}

func (f *Font) getGlyphSideBearingVar(gid GID, isVertical bool) int16 {
	extents, phantoms := f.getGlyfPoints(gid, true)
	if isVertical {
		return ceil(phantoms[phantomTop].Y - extents.YBearing)
	}
	return int16(phantoms[phantomLeft].X)
}

// take variations into account
func (f *Font) getHorizontalSideBearing(glyph GID) int16 {
	// base side bearing
	sideBearing := f.Hmtx.getSideBearing(glyph)
	if !f.isVar() {
		return sideBearing
	}
	if f.hvar != nil {
		return sideBearing + int16(f.hvar.getSideBearingVar(glyph, f.varCoords))
	}
	return f.getGlyphSideBearingVar(glyph, false)
}

// take variations into account
func (f *Font) getVerticalSideBearing(glyph GID) int16 {
	// base side bearing
	sideBearing := f.vmtx.getSideBearing(glyph)
	if !f.isVar() {
		return sideBearing
	}
	if f.vvar != nil {
		return sideBearing + int16(f.vvar.getSideBearingVar(glyph, f.varCoords))
	}
	return f.getGlyphSideBearingVar(glyph, true)
}

func (f *Font) GlyphHOrigin(GID) (x, y int32, found bool) {
	// zero is the right value here
	return 0, 0, true
}

func (f *Font) GlyphVOrigin(glyph GID) (x, y int32, found bool) {
	x = int32(f.HorizontalAdvance(glyph) / 2)

	if f.vorg != nil {
		y = int32(f.vorg.getYOrigin(glyph))
		return x, y, true
	}

	if extents, ok := f.getExtentsFromGlyf(glyph); ok {
		tsb := f.getVerticalSideBearing(glyph)
		y = int32(extents.YBearing) + int32(tsb)
		return x, y, true
	}

	fontExtents, ok := f.FontHExtents()
	y = int32(fontExtents.Ascender)

	return x, y, ok
}

func (f *Font) getExtentsFromGlyf(glyph GID) (font.GlyphExtents, bool) {
	if int(glyph) >= len(f.Glyf) {
		return font.GlyphExtents{}, false
	}
	g := f.Glyf[glyph]
	if f.isVar() { // we have to compute the outline points and apply variations
		extents, _ := f.getGlyfPoints(glyph, true)
		return extents, true
	}
	return g.getExtents(f.Hmtx, glyph), true
}

func (f *Font) getExtentsFromCBDT(glyph GID, xPpem, yPpem uint16) (font.GlyphExtents, bool) {
	strike := f.bitmap.chooseStrike(xPpem, yPpem)
	if strike == nil || strike.ppemX == 0 || strike.ppemY == 0 {
		return font.GlyphExtents{}, false
	}
	subtable := strike.findTable(glyph)
	if subtable == nil {
		return font.GlyphExtents{}, false
	}
	image := subtable.getImage(glyph)
	if image == nil {
		return font.GlyphExtents{}, false
	}
	extents := image.metrics.glyphExtents()

	/* convert to font units. */
	xScale := float32(f.upem) / float32(strike.ppemX)
	yScale := float32(f.upem) / float32(strike.ppemY)
	extents.XBearing *= xScale
	extents.YBearing *= yScale
	extents.Width *= xScale
	extents.Height *= yScale
	return extents, true
}

func (f *Font) getExtentsFromSbix(glyph GID, xPpem, yPpem uint16) (font.GlyphExtents, bool) {
	strike := f.sbix.chooseStrike(xPpem, yPpem)
	if strike == nil || strike.ppem == 0 {
		return font.GlyphExtents{}, false
	}
	data := strike.getGlyph(glyph, 0)
	if data.isNil() {
		return font.GlyphExtents{}, false
	}
	extents, ok := data.glyphExtents()

	/* convert to font units. */
	scale := float32(f.upem) / float32(strike.ppem)
	extents.XBearing *= scale
	extents.YBearing *= scale
	extents.Width *= scale
	extents.Height *= scale
	return extents, ok
}

func (f *Font) getExtentsFromCff1(glyph GID) (font.GlyphExtents, bool) {
	if f.cff == nil {
		return font.GlyphExtents{}, false
	}
	_, bounds, err := f.cff.LoadGlyph(glyph)
	if err != nil {
		return font.GlyphExtents{}, false
	}
	return bounds.ToExtents(), true
}

// func (f *fontMetrics) getExtentsFromCff2(glyph , coords []float32) (font.GlyphExtents, bool) {
// }

func (f *Font) GlyphExtents(glyph GID, xPpem, yPpem uint16) (font.GlyphExtents, bool) {
	out, ok := f.getExtentsFromSbix(glyph, xPpem, yPpem)
	if ok {
		return out, ok
	}
	out, ok = f.getExtentsFromGlyf(glyph)
	if ok {
		return out, ok
	}
	out, ok = f.getExtentsFromCff1(glyph)
	if ok {
		return out, ok
	}
	out, ok = f.getExtentsFromCBDT(glyph, xPpem, yPpem)
	return out, ok
}
