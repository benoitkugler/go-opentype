// Package font provides an high level API to access
// Opentype font properties.
// See [opentype] and [tables] for a lower, more detailled API.
package font

import (
	"math"
)

// GID is used to identify glyphs in a font.
// It is mostly internal to the font and should not be confused with
// Unicode code points.
// Note that, despite Opentype font files using uint16, we choose to use uint32,
// to allow room for future extension.
type GID uint32

// EmptyGlyph represents an invisible glyph, which should not be drawn,
// but whose advance and offsets should still be accounted for when rendering.
const EmptyGlyph GID = math.MaxUint32

// CmapIter is an interator over a Cmap.
type CmapIter interface {
	// Next returns true if the iterator still has data to yield
	Next() bool

	// Char must be called only when `Next` has returned `true`
	Char() (rune, GID)
}

// Cmap stores a compact representation of a cmap,
// offering both on-demand rune lookup and full rune range.
// It is conceptually equivalent to a map[rune]GID, but is often
// implemented more efficiently.
type Cmap interface {
	// Iter returns a new iterator over the cmap
	// Multiple iterators may be used over the same cmap
	// The returned interface is garanted not to be nil.
	Iter() CmapIter

	// Lookup avoid the construction of a map and provides
	// an alternative when only few runes need to be fetched.
	// It returns a default value and false when no glyph is provided.
	Lookup(rune) (GID, bool)
}
