package tables

import (
	"testing"

	tu "github.com/benoitkugler/go-opentype/testutils"
)

func TestParseSVG(t *testing.T) {
	fp := readFontFile(t, "toys/chromacheck-svg.ttf")
	_, _, err := ParseSVG(readTable(t, fp, "SVG "))
	tu.AssertNoErr(t, err)
}

func TestParseCFF(t *testing.T) {
	fp := readFontFile(t, "toys/CFFTest.otf")
	readTable(t, fp, "CFF ")
}
