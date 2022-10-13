package tables

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/opentype"
)

// filenames return the "absolute" file names of the given directory
func filenames(t *testing.T, dir string) []string {
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

func TestParseMorx(t *testing.T) {
	files := filenames(t, "morx")
	for _, filename := range files {
		file, err := td.Files.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		fp, err := opentype.NewLoader(bytes.NewReader(file))
		if err != nil {
			t.Fatal(filename, err)
		}
		rawTable, err := fp.RawTable(opentype.MustNewTag("morx"))
		if err != nil {
			t.Fatal(filename, err)
		}
		table, _, err := ParseMorx(rawTable)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(table)
	}
}
