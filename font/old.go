package font

// FontExtents exposes font-wide extent values, measured in font units.
// Note that typically ascender is positive and descender negative in coordinate systems that grow up.
type FontExtents struct {
	Ascender  float32 // Typographic ascender.
	Descender float32 // Typographic descender.
	LineGap   float32 // Suggested line spacing gap.
}

// LineMetric identifies one metric about the font.
type LineMetric uint8

const (
	// Distance above the baseline of the top of the underline.
	// Since most fonts have underline positions beneath the baseline, this value is typically negative.
	UnderlinePosition LineMetric = iota

	// Suggested thickness to draw for the underline.
	UnderlineThickness

	// Distance above the baseline of the top of the strikethrough.
	StrikethroughPosition

	// Suggested thickness to draw for the strikethrough.
	StrikethroughThickness

	SuperscriptEmYSize
	SuperscriptEmXOffset

	SubscriptEmYSize
	SubscriptEmYOffset
	SubscriptEmXOffset

	CapHeight
	XHeight
)

// GlyphExtents exposes extent values, measured in font units.
// Note that height is negative in coordinate systems that grow up.
type GlyphExtents struct {
	XBearing float32 // Left side of glyph from origin
	YBearing float32 // Top side of glyph from origin
	Width    float32 // Distance from left to right side
	Height   float32 // Distance from top to bottom side
}

// FaceMetrics exposes details of the font content.
// Implementation must be valid map keys to simplify caching.
type FaceMetrics interface {
	// Upem returns the units per em of the font file.
	// If not found, should return 1000 as fallback value.
	// This value is only relevant for scalable fonts.
	Upem() uint16

	// GlyphName returns the name of the given glyph, or an empty
	// string if the glyph is invalid or has no name.
	GlyphName(gid GID) string

	// LineMetric returns the metric identified by `metric` (in fonts units), or false
	// if the font does not provide such information.
	LineMetric(metric LineMetric) (float32, bool)

	// FontHExtents returns the extents of the font for horizontal text, or false
	// it not available, in font units.
	// `varCoords` (in normalized coordinates) is only useful for variable fonts.
	FontHExtents() (FontExtents, bool)

	// FontVExtents is the same as `FontHExtents`, but for vertical text.
	FontVExtents() (FontExtents, bool)

	// NominalGlyph returns the glyph used to represent the given rune,
	// or false if not found.
	NominalGlyph(ch rune) (GID, bool)

	// HorizontalAdvance returns the horizontal advance in font units.
	// When no data is available but the glyph index is valid, this method
	// should return a default value (the upem number for example).
	// If the glyph is invalid it should return 0.
	// `coords` is used by variable fonts, and is specified in normalized coordinates.
	HorizontalAdvance(gid GID) float32

	// VerticalAdvance is the same as `HorizontalAdvance`, but for vertical advance.
	VerticalAdvance(gid GID) float32

	// GlyphHOrigin fetches the (X,Y) coordinates of the origin (in font units) for a glyph ID,
	// for horizontal text segments.
	// Returns `false` if not available.
	GlyphHOrigin(GID) (x, y int32, found bool)

	// GlyphVOrigin is the same as `GlyphHOrigin`, but for vertical text segments.
	GlyphVOrigin(GID) (x, y int32, found bool)

	// GlyphExtents retrieve the extents for a specified glyph, of false, if not available.
	// `coords` is used by variable fonts, and is specified in normalized coordinates.
	// For bitmap glyphs, the closest resolution to `xPpem` and `yPpem` is selected.
	GlyphExtents(glyph GID, xPpem, yPpem uint16) (GlyphExtents, bool)
}

// FaceRenderer exposes access to glyph contents
type FaceRenderer interface {
	// GlyphData loads the glyph content, or return nil
	// if 'gid' is not supported.
	// For bitmap glyphs, the closest resolution to `xPpem` and `yPpem` is selected.
	GlyphData(gid GID, xPpem, yPpem uint16) GlyphData
}

// GlyphData describe how to graw a glyph.
// It is either an GlyphOutline, GlyphSVG or GlyphBitmap.
type GlyphData interface {
	isGlyphData()
}

func (GlyphOutline) isGlyphData() {}
func (GlyphSVG) isGlyphData()     {}
func (GlyphBitmap) isGlyphData()  {}

// GlyphOutline exposes the path to draw for
// vector glyph.
// Coordinates are expressed in fonts units.
type GlyphOutline struct {
	Segments []Segment
}

type SegmentOp uint8

const (
	SegmentOpMoveTo SegmentOp = iota
	SegmentOpLineTo
	SegmentOpQuadTo
	SegmentOpCubeTo
)

type SegmentPoint struct {
	X, Y float32 // expressed in fonts units
}

// Move translates the point.
func (pt *SegmentPoint) Move(dx, dy float32) {
	pt.X += dx
	pt.Y += dy
}

type Segment struct {
	Op SegmentOp
	// Args is up to three (x, y) coordinates, depending on the
	// operation.
	// The Y axis increases up.
	Args [3]SegmentPoint
}

// ArgsSlice returns the effective slice of points
// used (whose length is between 1 and 3).
func (s *Segment) ArgsSlice() []SegmentPoint {
	switch s.Op {
	case SegmentOpMoveTo, SegmentOpLineTo:
		return s.Args[0:1]
	case SegmentOpQuadTo:
		return s.Args[0:2]
	case SegmentOpCubeTo:
		return s.Args[0:3]
	default:
		panic("unreachable")
	}
}

// GlyphSVG is an SVG description for the glyph,
// as found in Opentype SVG table.
type GlyphSVG struct {
	// The SVG image content, decompressed if needed.
	// The actual glyph description is an SVG element
	// with id="glyph<GID>" (as in id="glyph12"),
	// and several glyphs may share the same Source
	Source []byte

	// According to the specification, a fallback outline
	// should be specified for each SVG glyphs
	Outline GlyphOutline
}

type GlyphBitmap struct {
	// The actual image content, whose interpretation depends
	// on the Format field.
	Data          []byte
	Format        BitmapFormat
	Width, Height int // number of columns and rows
}

// BitmapFormat identifies the format on the glyph
// raw data. Across the various font files, many formats
// may be encountered : black and white bitmaps, PNG, TIFF, JPG.
type BitmapFormat uint8

const (
	_ BitmapFormat = iota
	BlackAndWhite
	PNG
	JPG
	TIFF
)

// BitmapSize expose the size of bitmap glyphs.
// One font may contain several sizes.
type BitmapSize struct {
	Height, Width uint16
	XPpem, YPpem  uint16
}
