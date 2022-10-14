package tables

//go:generate ../../binarygen/cmd/generator post.go

// PostScript table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/post
type Post struct {
	version     postVersion
	italicAngle uint32
	// UnderlinePosition is the suggested distance of the top of the
	// underline from the baseline (negative values indicate below baseline).
	underlinePosition int16
	// Suggested values for the underline thickness.
	underlineThickness int16
	// IsFixedPitch indicates that the font is not proportionally spaced
	// (i.e. monospaced).
	isFixedPitch uint32
	memoryUsage  [4]uint32
	names        postNames `version-field:"version"`
}

type postNames interface {
	isPostNames()
}

func (postNames10) isPostNames() {}
func (postNames20) isPostNames() {}
func (postNames30) isPostNames() {}

type postVersion uint32

const (
	postVersion10 postVersion = 0x00010000
	postVersion20 postVersion = 0x00020000
	postVersion30 postVersion = 0x00030000
)

type postNames10 struct{}

type postNames20 struct {
	glyphNameIndexes []uint16 `len:"_first16"` // size numGlyph
	stringData       []byte   `len:"__toEnd"`
}

type postNames30 postNames10
