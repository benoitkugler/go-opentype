package tables

import (
	"testing"
)

func TestParseSVG(t *testing.T) {
	fp := readFontFile(t, "toys/chromacheck-svg.ttf")
	_, _, err := ParseSVG(readTable(t, fp, "SVG "))
	assertNoErr(t, err)
}

func TestParseCFF(t *testing.T) {
	fp := readFontFile(t, "toys/CFFTest.otf")
	readTable(t, fp, "CFF ")
}
