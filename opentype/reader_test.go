package opentype

import (
	"bytes"
	"math/rand"
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
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

func TestParseCrashers(t *testing.T) {
	font, err := NewLoader(bytes.NewReader([]byte{}))
	if font != nil || err == nil {
		t.Fatal()
	}

	for range [50]int{} {
		L := rand.Intn(100)
		input := make([]byte, L)
		rand.Read(input)

		_, err = NewLoader(bytes.NewReader(input))
		if err == nil {
			t.Error("expected error on random input")
		}

		_, err = NewLoaders(bytes.NewReader(input))
		if font != nil || err == nil {
			t.Error("expected error on random input")
		}
	}
}

func TestCollection(t *testing.T) {
	for _, filename := range filenames(t, "collections") {
		f, err := td.Files.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		fonts, err := NewLoaders(bytes.NewReader(f))
		if err != nil {
			t.Fatal(filename, err)
		}
		for _, font := range fonts {
			if len(font.tables) == 0 {
				t.Fatal()
			}
		}
	}

	// check that it also works for single font files
	for _, filename := range filenames(t, "common") {
		f, err := td.Files.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		fonts, err := NewLoaders(bytes.NewReader(f))
		if err != nil {
			t.Fatal(filename, err)
		}
		if len(fonts) != 1 {
			t.Fatal(filename, "expected only one font")
		}
	}
}
