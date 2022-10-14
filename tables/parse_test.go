package tables

import (
	"bytes"
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/opentype"
)

// filenames return the "absolute" file names of the given directory
func filenames(t *testing.T, dir string) []string {
	t.Helper()

	files, err := td.Files.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, entry := range files {
		filename := filepath.Join(dir, entry.Name())
		out = append(out, filename)
	}
	return out
}

func readFontFile(t *testing.T, filepath string) *opentype.Loader {
	t.Helper()

	file, err := td.Files.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	fp, err := opentype.NewLoader(bytes.NewReader(file))
	if err != nil {
		t.Fatal(filepath, err)
	}
	return fp
}

func readTable(t *testing.T, fl *opentype.Loader, tag string) []byte {
	t.Helper()

	table, err := fl.RawTable(opentype.MustNewTag(tag))
	if err != nil {
		t.Fatal(err)
	}
	return table
}

func numGlyphs(t *testing.T, fp *opentype.Loader) int {
	t.Helper()

	table := readTable(t, fp, "maxp")
	maxp, _, err := ParseMaxp(table)
	if err != nil {
		t.Fatal(err)
	}
	return int(maxp.numGlyphs)
}

func TestParseBasicTables(t *testing.T) {
	for _, filename := range append(filenames(t, "morx"), filenames(t, "common")...) {
		fp := readFontFile(t, filename)
		_, _, err := ParseOs2(readTable(t, fp, "OS/2"))
		if err != nil {
			t.Fatal(err)
		}
		_, _, err = ParseHead(readTable(t, fp, "head"))
		if err != nil {
			t.Fatal(err)
		}
		_, _, err = ParseMaxp(readTable(t, fp, "maxp"))
		if err != nil {
			t.Fatal(err)
		}
		_, _, err = ParseName(readTable(t, fp, "name"))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseMorx(t *testing.T) {
	files := filenames(t, "morx")
	for _, filename := range files {
		fp := readFontFile(t, filename)
		table, _, err := ParseMorx(readTable(t, fp, "morx"))
		if err != nil {
			t.Fatal(err)
		}
		if int(table.nChains) != len(table.chains) {
			t.Fatal()
		}
		ng := numGlyphs(t, fp)
		for _, chain := range table.chains {
			subtables, err := chain.process(ng)
			if err != nil {
				t.Fatal(err)
			}
			if len(subtables) != int(chain.nSubtable) {
				t.Fatal()
			}
		}
	}
}
