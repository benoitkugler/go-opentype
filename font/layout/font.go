package layout

import (
	"github.com/benoitkugler/go-opentype/font"
	"github.com/benoitkugler/go-opentype/font/cff"
	"github.com/benoitkugler/go-opentype/opentype"
	"github.com/benoitkugler/go-opentype/tables"
)

type (
	GID = font.GID
	Tag = opentype.Tag
)

// Font represents one Opentype font file (or one sub font of a collection).
// It is an educated view of the underlying font file, optimized for quick access
// to information required by text layout engines.
//
// All its methods are read-only and a [*Font] object is thus safe for concurrent use.
type Font struct {
	cmap    font.Cmap
	cmapVar font.UnicodeVariations

	names tables.Name

	hhea *tables.Hhea
	vhea *tables.Vhea
	vorg *tables.VORG // optional
	cff  *cff.Font
	post post       // optional
	svg  tables.SVG // optional

	// Optionnal, only present in variable fonts

	hvar *tables.HVAR // optional
	vvar tables.VVAR  // optional
	avar tables.Avar
	mvar mvar
	gvar gvar
	fvar tables.Fvar

	Glyf tables.Glyf
	hmtx tables.Hmtx
	vmtx tables.Vmtx
	// bitmap bitmapTable // CBDT or EBLC or BLOC
	sbix tables.Sbix

	os2 os2

	// Advanced layout tables.

	GDEF tables.GDEF // An absent table has a nil Class
	Trak tables.Trak
	Ankr tables.Ankr
	Feat tables.Feat
	Morx tables.Morx
	Kern tables.Kern
	Kerx tables.Kerx
	GSUB tables.GSUB // An absent table has a nil slice of lookups
	GPOS tables.GPOS // An absent table has a nil slice of lookups

	head tables.Head

	upem uint16 // cached value
}

// Face is a font with (optional) user-provided variable coordinates.
type Face struct {
	Font *Font

	// Coords are the current variable coordinates, expressed in normalized units.
	// It is empty for non variable fonts.
	// Use `NormalizeVariations` to convert from design space units.
	Coords []float32
}
