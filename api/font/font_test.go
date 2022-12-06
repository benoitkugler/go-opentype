package font

import (
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

// filenames return the "absolute" file names of the given directory
// excluding directories, and not recursing
func filenames(t testing.TB, dir string) []string {
	t.Helper()

	files, err := td.Files.ReadDir(dir)
	assertNoErr(t, err)

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

func TestCrashes(t *testing.T) {
	for _, filepath := range append(filenames(t, "common"), "toys/chromacheck-svg.ttf") {
		loadFont(t, filepath)
	}
}
