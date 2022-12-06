package font

import (
	"fmt"

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

	name tables.Name

	hhea *tables.Hhea
	vhea *tables.Vhea
	vorg *tables.VORG // optional
	cff  *cff.Font
	post post // optional
	svg  svg  // optional

	// Optionnal, only present in variable fonts

	hvar *tables.HVAR // optional
	vvar *tables.VVAR // optional
	avar tables.Avar
	mvar mvar
	gvar gvar
	fvar tables.Fvar

	glyf   tables.Glyf
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
	GSUB GSUB // An absent table has a nil slice of lookups
	GPOS GPOS // An absent table has a nil slice of lookups

	head tables.Head

	upem uint16 // cached value
}

// NewFont load all the font tables, sanitizing them.
// An error is returned for invalid tables, or when mandatory one
// are missing.
func NewFont(ld *loader.Loader) (*Font, error) {
	var (
		out Font
		err error
	)

	raw, err := ld.RawTable(loader.MustNewTag("cmap"))
	if err != nil {
		return nil, err
	}
	tb, _, err := tables.ParseCmap(raw)
	if err != nil {
		return nil, err
	}
	out.cmap, out.cmapVar, err = api.ProcessCmap(tb)
	if err != nil {
		return nil, err
	}

	raw, err = ld.RawTable(loader.MustNewTag("head"))
	if err != nil {
		return nil, err
	}
	out.head, _, err = tables.ParseHead(raw)
	if err != nil {
		return nil, err
	}

	raw, err = ld.RawTable(loader.MustNewTag("maxp"))
	if err != nil {
		return nil, err
	}
	maxp, _, err := tables.ParseMaxp(raw)
	if err != nil {
		return nil, err
	}

	raw, err = ld.RawTable(loader.MustNewTag("fvar"))
	if err == nil { // error only if the table is present and invalid
		out.fvar, _, err = tables.ParseFvar(raw)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("avar"))
	if err == nil { // error only if the table is present and invalid
		out.avar, _, err = tables.ParseAvar(raw)
		if err != nil {
			return nil, err
		}
	}

	out.upem = out.head.Upem()

	raw, err = ld.RawTable(loader.MustNewTag("OS/2"))
	if err == nil { // error only if the table is present and invalid
		os2, _, err := tables.ParseOs2(raw)
		if err != nil {
			return nil, err
		}
		out.os2, err = newOs2(os2)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("loca"))
	if err == nil { // error only if the table is present and invalid
		loca, err := tables.ParseLoca(raw, int(maxp.NumGlyphs), out.head.IndexToLocFormat == 1)
		if err != nil {
			return nil, err
		}

		raw, err = ld.RawTable(loader.MustNewTag("glyf"))
		if err != nil {
			return nil, err
		}
		out.glyf, err = tables.ParseGlyf(raw, loca)
		if err != nil {
			return nil, err
		}
	}

	out.bitmap = selectBitmapTable(ld)

	raw, err = ld.RawTable(loader.MustNewTag("sbix"))
	if err == nil { // error only if the table is present and invalid
		out.sbix, _, err = tables.ParseSbix(raw, int(maxp.NumGlyphs))
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("CFF "))
	if err == nil { // error only if the table is present and invalid
		out.cff, err = cff.ParseFont(raw)
		if err != nil {
			return nil, err
		}

		if N := len(out.cff.Charstrings); N != int(maxp.NumGlyphs) {
			return nil, fmt.Errorf("invalid number of glyphs in CFF table (%d != %d)", N, maxp.NumGlyphs)
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("post"))
	if err == nil { // error only if the table is present and invalid
		post, _, err := tables.ParsePost(raw)
		if err != nil {
			return nil, err
		}

		out.post, err = newPost(post)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("SVG "))
	if err == nil { // error only if the table is present and invalid
		svg, _, err := tables.ParseSVG(raw)
		if err != nil {
			return nil, err
		}

		out.svg, err = newSvg(svg)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("hhea"))
	if err == nil { // error only if the table is present and invalid
		hhea, _, err := tables.ParseHhea(raw)
		if err != nil {
			return nil, err
		}

		out.hhea = &hhea

		raw, err = ld.RawTable(loader.MustNewTag("hmtx"))
		if err == nil { // error only if the table is present and invalid
			out.hmtx, _, err = tables.ParseHmtx(raw, int(hhea.NumOfLongMetrics), int(maxp.NumGlyphs)-int(hhea.NumOfLongMetrics))
			if err != nil {
				return nil, err
			}
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("vhea"))
	if err == nil { // error only if the table is present and invalid
		vhea, _, err := tables.ParseHhea(raw)
		if err != nil {
			return nil, err
		}

		out.vhea = &vhea

		raw, err = ld.RawTable(loader.MustNewTag("vmtx"))
		if err == nil { // error only if the table is present and invalid
			out.vmtx, _, err = tables.ParseHmtx(raw, int(vhea.NumOfLongMetrics), int(maxp.NumGlyphs)-int(vhea.NumOfLongMetrics))
			if err != nil {
				return nil, err
			}
		}
	}

	if len(out.fvar.Axis) != 0 {
		raw, err = ld.RawTable(loader.MustNewTag("MVAR"))
		if err == nil { // error only if the table is present and invalid
			mvar, _, err := tables.ParseMVAR(raw)
			if err != nil {
				return nil, err
			}

			out.mvar = newMvar(mvar)
		}

		raw, err = ld.RawTable(loader.MustNewTag("gvar"))
		if err == nil { // error only if the table is present and invalid
			gvar, _, err := tables.ParseGvar(raw)
			if err != nil {
				return nil, err
			}

			out.gvar, err = newGvar(gvar, out.glyf)
			if err != nil {
				return nil, err
			}
		}

		raw, err = ld.RawTable(loader.MustNewTag("HVAR"))
		if err == nil { // error only if the table is present and invalid
			hvar, _, err := tables.ParseHVAR(raw)
			if err != nil {
				return nil, err
			}

			out.hvar = &hvar
		}
		raw, err = ld.RawTable(loader.MustNewTag("VVAR"))
		if err == nil { // error only if the table is present and invalid
			vvar, _, err := tables.ParseHVAR(raw)
			if err != nil {
				return nil, err
			}

			out.vvar = &vvar
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("VORG"))
	if err == nil { // error only if the table is present and invalid
		vorg, _, err := tables.ParseVORG(raw)
		if err != nil {
			return nil, err
		}

		out.vorg = &vorg
	}

	// layout tables

	raw, err = ld.RawTable(loader.MustNewTag("GDEF"))
	if err == nil { // error only if the table is present and invalid
		out.GDEF, _, err = tables.ParseGDEF(raw)
		if err != nil {
			return nil, err
		}

		err = sanitizeGDEF(out.GDEF, len(out.fvar.Axis))
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("GSUB"))
	if err == nil { // error only if the table is present and invalid
		layout, _, err := tables.ParseLayout(raw)
		if err != nil {
			return nil, err
		}
		out.GSUB, err = newGSUB(layout)
		if err != nil {
			return nil, err
		}
	}
	raw, err = ld.RawTable(loader.MustNewTag("GPOS"))
	if err == nil { // error only if the table is present and invalid
		layout, _, err := tables.ParseLayout(raw)
		if err != nil {
			return nil, err
		}
		out.GPOS, err = newGPOS(layout)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("morx"))
	if err == nil { // error only if the table is present and invalid
		out.Morx, _, err = tables.ParseMorx(raw, int(maxp.NumGlyphs))
		if err != nil {
			return nil, err
		}
	}
	raw, err = ld.RawTable(loader.MustNewTag("kerx"))
	if err == nil { // error only if the table is present and invalid
		out.Kerx, _, err = tables.ParseKerx(raw, int(maxp.NumGlyphs))
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("kern"))
	if err == nil { // error only if the table is present and invalid
		out.Kern, _, err = tables.ParseKern(raw)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ld.RawTable(loader.MustNewTag("ankr"))
	if err == nil { // error only if the table is present and invalid
		out.Ankr, _, err = tables.ParseAnkr(raw, int(maxp.NumGlyphs))
		if err != nil {
			return nil, err
		}
	}
	raw, err = ld.RawTable(loader.MustNewTag("trak"))
	if err == nil { // error only if the table is present and invalid
		out.Trak, _, err = tables.ParseTrak(raw)
		if err != nil {
			return nil, err
		}
	}
	raw, err = ld.RawTable(loader.MustNewTag("feat"))
	if err == nil { // error only if the table is present and invalid
		out.Feat, _, err = tables.ParseFeat(raw)
		if err != nil {
			return nil, err
		}
	}

	return &out, nil
}

// return nil if no table is valid (or present)
func selectBitmapTable(ld *loader.Loader) bitmap {
	color, err := loadBitmap(ld, loader.MustNewTag("CBLC"), loader.MustNewTag("CBDT"))
	if err == nil {
		return color
	}

	gray, err := loadBitmap(ld, loader.MustNewTag("EBLC"), loader.MustNewTag("EBDT"))
	if err == nil {
		return gray
	}

	apple, err := loadBitmap(ld, loader.MustNewTag("bloc"), loader.MustNewTag("bdat"))
	if err == nil {
		return apple
	}

	return nil
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