package tables

type BitmapSubtable struct {
	FirstGlyph GlyphID //	First glyph ID of this range.
	LastGlyph  GlyphID //	Last glyph ID of this range (inclusive).
	IndexSubHeader
}

// EBLC is the Embedded Bitmap Location Table
// See - https://learn.microsoft.com/fr-fr/typography/opentype/spec/eblc
type EBLC = CBLC
