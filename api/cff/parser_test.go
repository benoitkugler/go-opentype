package cff

import (
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/tables"
	tu "github.com/benoitkugler/go-opentype/testutils"
)

func TestParseCFF(t *testing.T) {
	for _, filepath := range tu.Filenames(t, "cff") {
		content, err := td.Files.ReadFile(filepath)
		tu.AssertNoErr(t, err)

		font, err := Parse(content)
		tu.AssertNoErr(t, err)

		tu.Assert(t, len(font.Charstrings) >= 12)

		if font.fdSelect != nil {
			for i := 0; i < len(font.Charstrings); i++ {
				_, err = font.fdSelect.fontDictIndex(tables.GlyphID(i))
				tu.AssertNoErr(t, err)
			}
		}

		for glyphIndex := range font.Charstrings {
			_, _, err := font.LoadGlyph(tables.GlyphID(glyphIndex))
			tu.AssertNoErr(t, err)
		}
	}
}
