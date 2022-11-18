package tables

import (
	"bytes"
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/opentype"
)

func assert(t *testing.T, b bool) {
	t.Helper()
	assertC(t, b, "assertion error")
}

func assertC(t *testing.T, b bool, context string) {
	t.Helper()
	if !b {
		t.Fatal(context)
	}
}

// filenames return the "absolute" file names of the given directory
// excluding directories, and not recursing
func filenames(t *testing.T, dir string) []string {
	t.Helper()

	files, err := td.Files.ReadDir(dir)
	assert(t, err == nil)

	var out []string
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		filename := filepath.Join(dir, entry.Name())
		out = append(out, filename)
	}
	return out
}

func readFontFile(t *testing.T, filepath string) *opentype.Loader {
	t.Helper()

	file, err := td.Files.ReadFile(filepath)
	assert(t, err == nil)

	fp, err := opentype.NewLoader(bytes.NewReader(file))
	assertC(t, err == nil, filepath)

	return fp
}

func readTable(t *testing.T, fl *opentype.Loader, tag string) []byte {
	t.Helper()

	table, err := fl.RawTable(opentype.MustNewTag(tag))
	assert(t, err == nil)

	return table
}

func numGlyphs(t *testing.T, fp *opentype.Loader) int {
	t.Helper()

	table := readTable(t, fp, "maxp")
	maxp, _, err := ParseMaxp(table)
	assert(t, err == nil)

	return int(maxp.numGlyphs)
}

func TestParseBasicTables(t *testing.T) {
	for _, filename := range append(filenames(t, "morx"), filenames(t, "common")...) {
		fp := readFontFile(t, filename)
		_, _, err := ParseOs2(readTable(t, fp, "OS/2"))
		assert(t, err == nil)

		_, _, err = ParseHead(readTable(t, fp, "head"))
		assert(t, err == nil)

		_, _, err = ParseMaxp(readTable(t, fp, "maxp"))
		assert(t, err == nil)

		_, _, err = ParseName(readTable(t, fp, "name"))
		assert(t, err == nil)
	}
}

func TestParseMorx(t *testing.T) {
	files := filenames(t, "morx")
	for _, filename := range files {
		fp := readFontFile(t, filename)
		table, _, err := ParseMorx(readTable(t, fp, "morx"))
		assert(t, err == nil)
		assert(t, int(table.nChains) == len(table.chains))

		ng := numGlyphs(t, fp)
		for _, chain := range table.chains {
			subtables, err := chain.process(ng)
			assert(t, err == nil)
			assert(t, len(subtables) == int(chain.nSubtable))
		}
	}
}

func TestParseCmap(t *testing.T) {
	// general parsing
	for _, filename := range filenames(t, "common") {
		fp := readFontFile(t, filename)
		cmap, _, err := ParseCmap(readTable(t, fp, "cmap"))
		assert(t, err == nil)
		assert(t, len(cmap.records) == len(cmap.subtables))
	}

	// specialized tests for each format
	for _, filename := range filenames(t, "cmap") {
		fp := readFontFile(t, filename)
		cmap, _, err := ParseCmap(readTable(t, fp, "cmap"))
		assert(t, err == nil)
		assert(t, len(cmap.records) == len(cmap.subtables))
	}

	// test format 2 through a single table
	file, err := td.Files.ReadFile("cmap/table/CMAP2.bin")
	assert(t, err == nil)
	cmap, _, err := ParseCmapSubtable2(file)
	assert(t, err == nil)
	assert(t, cmap.format == 2)
}
