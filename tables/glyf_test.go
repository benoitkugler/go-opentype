package tables

import (
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

func TestParseGlyf(t *testing.T) {
	for _, file := range td.WithGlyphs {
		fp := readFontFile(t, file.Path)
		head, _, err := ParseHead(readTable(t, fp, "head"))
		assertNoErr(t, err)

		maxp, _, err := ParseMaxp(readTable(t, fp, "maxp"))
		assertNoErr(t, err)

		loca, err := parseLoca(readTable(t, fp, "loca"), int(maxp.numGlyphs), head.IndexToLocFormat == 1)
		assertNoErr(t, err)

		glyphs, err := parseGlyf(readTable(t, fp, "glyf"), loca)
		assertNoErr(t, err)
		assert(t, len(glyphs) == len(loca)-1)
		assert(t, len(glyphs) == file.GlyphNumber)
	}
}
