package font

import (
	"github.com/benoitkugler/go-opentype/api"
	"github.com/benoitkugler/go-opentype/api/cff"
	"github.com/benoitkugler/go-opentype/loader"
	"github.com/benoitkugler/go-opentype/tables"
)

type (
	GID = api.GID
	Tag = loader.Tag
)

// Font represents one Opentype font file (or one sub font of a collection).
// It is an educated view of the underlying font file, optimized for quick access
// to information required by text layout engines.
//
// All its methods are read-only and a [*Font] object is thus safe for concurrent use.
type Font struct {
	cmap    api.Cmap
	cmapVar api.UnicodeVariations

	names tables.Name

	hhea *tables.Hhea
	vhea *tables.Vhea
	vorg *tables.VORG // optional
	cff  *cff.Font
	post post       // optional
	svg  tables.SVG // optional

	// Optionnal, only present in variable fonts

	hvar *tables.HVAR // optional
	vvar *tables.VVAR // optional
	avar tables.Avar
	mvar mvar
	gvar gvar
	fvar tables.Fvar

	Glyf   tables.Glyf
	hmtx   tables.Hmtx
	vmtx   tables.Vmtx
	bitmap bitmap
	sbix   tables.Sbix

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

// Face is a font with user-provided settings.
// It is a lightweight wrapper around [*Font], NOT safe for concurrent use.
type Face struct {
	*Font

	// Coords are the current variable coordinates, expressed in normalized units.
	// It is empty for non variable fonts.
	// Use `NormalizeVariations` to convert from design space units.
	Coords []float32

	// Horizontal and vertical pixels-per-em (ppem), used to select bitmap sizes.
	XPpem, YPpem uint16
}
