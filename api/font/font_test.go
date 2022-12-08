package font

import (
	"testing"

	tu "github.com/benoitkugler/go-opentype/testutils"
)

func TestCrashes(t *testing.T) {
	for _, filepath := range append(tu.Filenames(t, "common"), "toys/chromacheck-svg.ttf") {
		loadFont(t, filepath)
	}
}

func BenchmarkLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, filepath := range tu.Filenames(b, "common") {
			loadFont(b, filepath)
		}
	}
}
