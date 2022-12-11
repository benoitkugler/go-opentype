package harfbuzz

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/harfbuzz"
	"github.com/benoitkugler/go-opentype/api"
	"github.com/benoitkugler/go-opentype/api/font"
	"github.com/benoitkugler/go-opentype/language"
	"github.com/benoitkugler/go-opentype/loader"
	"github.com/benoitkugler/go-opentype/testutils"
)

// ported from harfbuzz/util/hb-shape.cc, main-font-text.hh Copyright © 2010, 2011,2012  Google, Inc. Behdad Esfahbod

// parse and run the test cases directly copied from harfbuzz/test/shaping

type formatOptions struct {
	hideGlyphNames bool
	hidePositions  bool
	hideAdvances   bool
	hideClusters   bool
	showExtents    bool
	showFlags      bool
}

// Generates gidDDD if glyph has no name.
func (f *Font) glyphToString(glyph GID) string {
	if name := f.face.GlyphName(glyph); name != "" {
		return name
	}

	return fmt.Sprintf("gid%d", glyph)
}

// return a compact representation of the buffer contents
func (b *Buffer) serialize(font *Font, opt formatOptions) string {
	if len(b.Info) == 0 {
		return "" //  the reference does not return []
	}
	gs := new(strings.Builder)
	gs.WriteByte('[')
	var x, y Position
	for i, glyph := range b.Info {
		if opt.hideGlyphNames {
			fmt.Fprintf(gs, "%d", glyph.Glyph)
		} else {
			gs.WriteString(font.glyphToString(glyph.Glyph))
		}

		if !opt.hideClusters {
			fmt.Fprintf(gs, "=%d", glyph.Cluster)
		}
		pos := b.Pos[i]

		if !opt.hidePositions {
			if x+pos.XOffset != 0 || y+pos.YOffset != 0 {
				fmt.Fprintf(gs, "@%d,%d", x+pos.XOffset, y+pos.YOffset)
			}
			if !opt.hideAdvances {
				fmt.Fprintf(gs, "+%d", pos.XAdvance)
				if pos.YAdvance != 0 {
					fmt.Fprintf(gs, ",%d", pos.YAdvance)
				}
			}
		}

		if opt.showExtents {
			extents, _ := font.GlyphExtents(glyph.Glyph)
			fmt.Fprintf(gs, "<%d,%d,%d,%d>", extents.XBearing, extents.YBearing, extents.Width, extents.Height)
		}

		if i != len(b.Info)-1 {
			gs.WriteByte('|')
		}

		if opt.hideAdvances {
			x += pos.XAdvance
			y += pos.YAdvance
		}
	}
	gs.WriteByte(']')
	return gs.String()
}

type fontOptions struct {
	font *Font // cached value of getFont()

	fontRef    api.FontID
	variations []font.Variation

	subpixelBits         int
	fontSizeX, fontSizeY int
	ptem                 float64
	yPpem, xPpem         uint16
}

const fontSizeUpem = 0x7FFFFFFF

func newFontOptions() fontOptions {
	return fontOptions{
		subpixelBits: 0,
		fontSizeX:    fontSizeUpem,
		fontSizeY:    fontSizeUpem,
	}
}

func (fo *fontOptions) getFont() *Font {
	if fo.font != nil {
		return fo.font
	}

	/* Create the blob */
	if fo.fontRef.File == "" {
		check(errors.New("no font file specified"))
	}

	f, err := td.Files.ReadFile(fo.fontRef.File)
	check(err)

	fonts, err := loader.NewLoaders(bytes.NewReader(f))
	check(err)

	if int(fo.fontRef.Index) >= len(fonts) {
		check(fmt.Errorf("invalid font Index %d for length %d", fo.fontRef.Index, len(fonts)))
	}
	ft, err := font.NewFont(fonts[fo.fontRef.Index])
	check(err)

	/* Create the face */
	face := font.Face{Font: ft, XPpem: fo.xPpem, YPpem: fo.yPpem}
	face.SetVariations(fo.variations)

	fo.font = NewFont(&face)

	if fo.fontSizeX == fontSizeUpem {
		fo.fontSizeX = int(fo.font.faceUpem)
	}
	if fo.fontSizeY == fontSizeUpem {
		fo.fontSizeY = int(fo.font.faceUpem)
	}

	fo.font.Ptem = float32(fo.ptem)

	scaleX := scalbnf(float64(fo.fontSizeX), fo.subpixelBits)
	scaleY := scalbnf(float64(fo.fontSizeY), fo.subpixelBits)
	fo.font.XScale, fo.font.YScale = scaleX, scaleY

	return fo.font
}

func scalbnf(x float64, exp int) int32 {
	return int32(x * (math.Pow(2, float64(exp))))
}

// see variationsUsage
func (opts *fontOptions) parseVariations(s string) error {
	// remove possible quote
	s = strings.Trim(s, `"`)

	variations := strings.Split(s, ",")
	opts.variations = make([]font.Variation, len(variations))

	var err error
	for i, feature := range variations {
		opts.variations[i], err = ParseVariation(feature)
		if err != nil {
			return err
		}
	}
	return nil
}

type textInput struct {
	textBefore, textAfter []rune
	text                  []rune
}

func parseUnicodes(s string) ([]rune, error) {
	runes := strings.Split(s, ",")
	text := make([]rune, len(runes))
	for i, r := range runes {
		if _, err := fmt.Sscanf(r, "U+%x", &text[i]); err == nil {
			continue
		}
		if _, err := fmt.Sscanf(r, "0x%x", &text[i]); err == nil {
			continue
		}
		if _, err := fmt.Sscanf(r, "%x", &text[i]); err == nil {
			continue
		}
		return text, fmt.Errorf("invalid unicode rune : %s", r)
	}
	return text, nil
}

type testOptions struct {
	input    textInput
	shaper   shapeOptions
	fontOpts fontOptions
	format   formatOptions
}

type shapeOptions struct {
	shaper                    string
	features                  string
	props                     SegmentProperties
	invisibleGlyph            GID
	clusterLevel              ClusterLevel
	bot                       bool
	eot                       bool
	preserveDefaultIgnorables bool
	removeDefaultIgnorables   bool
}

func (so *shapeOptions) setupBuffer(buffer *Buffer) {
	buffer.Props = so.props
	var flags ShappingOptions
	if so.bot {
		flags |= Bot
	}
	if so.eot {
		flags |= Eot
	}
	if so.preserveDefaultIgnorables {
		flags |= PreserveDefaultIgnorables
	}
	if so.removeDefaultIgnorables {
		flags |= RemoveDefaultIgnorables
	}
	buffer.Flags = flags
	buffer.Invisible = so.invisibleGlyph
	buffer.ClusterLevel = so.clusterLevel
	buffer.GuessSegmentProperties()
}

func copyBufferProperties(dst, src *Buffer) {
	dst.Props = src.Props
	dst.Flags = src.Flags
	dst.ClusterLevel = src.ClusterLevel
}

func appendBuffer(dst, src *Buffer, start, end int) {
	origLen := len(dst.Info)

	dst.Info = append(dst.Info, src.Info[start:end]...)
	dst.Pos = append(dst.Pos, src.Pos[start:end]...)

	/* pre-context */
	if origLen == 0 && start+len(src.context[0]) > 0 {
		dst.clearContext(0)
		for start > 0 && len(dst.context[0]) < contextLength {
			start--
			dst.context[0] = append(dst.context[0], src.Info[start].codepoint)
		}

		for i := 0; i < len(src.context[0]) && len(dst.context[0]) < contextLength; i++ {
			dst.context[0] = append(dst.context[0], src.context[0][i])
		}
	}

	/* post-context */
	dst.clearContext(1)
	for end < len(src.Info) && len(dst.context[1]) < contextLength {
		dst.context[1] = append(dst.context[1], src.Info[end].codepoint)
		end++
	}
	for i := 0; i < len(src.context[1]) && len(dst.context[1]) < contextLength; i++ {
		dst.context[1] = append(dst.context[1], src.context[1][i])
	}
}

func (so *shapeOptions) populateBuffer(input textInput) *Buffer {
	buffer := NewBuffer()

	if input.textBefore != nil {
		buffer.AddRunes(input.textBefore, len(input.textBefore), 0)
	}

	buffer.AddRunes(input.text, 0, len(input.text))

	if input.textAfter != nil {
		buffer.AddRunes(input.textAfter, 0, 0)
	}

	so.setupBuffer(buffer)

	return buffer
}

func (so *shapeOptions) shape(font *Font, buffer *Buffer, verify bool) error {
	var textBuffer *Buffer

	if verify {
		textBuffer = NewBuffer()
		appendBuffer(textBuffer, buffer, 0, len(buffer.Info))
	}

	features, err := so.parseFeatures()
	if err != nil {
		return err
	}
	buffer.Shape(font, features)

	if verify {
		if err := so.verifyBuffer(buffer, textBuffer, font); err != nil {
			return err
		}
	}

	return nil
}

func (so *shapeOptions) verifyBuffer(buffer, textBuffer *Buffer, font *Font) error {
	if err := so.verifyBufferMonotone(buffer); err != nil {
		return err
	}
	if err := so.verifyBufferSafeToBreak(buffer, textBuffer, font); err != nil {
		return err
	}
	if err := so.verifyValidGID(buffer, font); err != nil {
		log.Println(err)
	}
	return nil
}

func (so *shapeOptions) verifyValidGID(buffer *Buffer, font *Font) error {
	for _, glyph := range buffer.Info {
		_, ok := font.GlyphExtents(glyph.Glyph)
		if !ok {
			return fmt.Errorf("Unknow glyph %d in font", glyph.Glyph)
		}
	}
	return nil
}

/* Check that clusters are monotone. */
func (so *shapeOptions) verifyBufferMonotone(buffer *Buffer) error {
	if so.clusterLevel == MonotoneGraphemes || so.clusterLevel == MonotoneCharacters {
		isForward := buffer.Props.Direction.isForward()

		info := buffer.Info

		for i := 1; i < len(info); i++ {
			if info[i-1].Cluster != info[i].Cluster && (info[i-1].Cluster < info[i].Cluster) != isForward {
				return fmt.Errorf("cluster at index %d is not monotone", i)
			}
		}
	}

	return nil
}

func (so *shapeOptions) verifyBufferSafeToBreak(buffer, textBuffer *Buffer, font *Font) error {
	if so.clusterLevel != MonotoneGraphemes && so.clusterLevel != MonotoneCharacters {
		/* Cannot perform this check without monotone clusters.
		 * Then again, unsafe-to-break flag is much harder to use without
		 * monotone clusters. */
		return nil
	}

	/* Check that breaking up shaping at safe-to-break is indeed safe. */

	fragment, reconstruction := NewBuffer(), NewBuffer()
	copyBufferProperties(reconstruction, buffer)

	info := buffer.Info
	text := textBuffer.Info

	/* Chop text and shape fragments. */
	forward := buffer.Props.Direction.isForward()
	start := 0
	textStart := len(textBuffer.Info)
	if forward {
		textStart = 0
	}
	textEnd := textStart
	for end := 1; end < len(buffer.Info)+1; end++ {
		offset := 1
		if forward {
			offset = 0
		}
		if end < len(buffer.Info) && (info[end].Cluster == info[end-1].Cluster ||
			info[end-offset].Mask&GlyphUnsafeToBreak != 0) {
			continue
		}

		/* Shape segment corresponding to glyphs start..end. */
		if end == len(buffer.Info) {
			if forward {
				textEnd = len(textBuffer.Info)
			} else {
				textStart = 0
			}
		} else {
			if forward {
				cluster := info[end].Cluster
				for textEnd < len(textBuffer.Info) && text[textEnd].Cluster < cluster {
					textEnd++
				}
			} else {
				cluster := info[end-1].Cluster
				for textStart != 0 && text[textStart-1].Cluster >= cluster {
					textStart--
				}
			}
		}
		if !(textStart < textEnd) {
			return fmt.Errorf("unexpected %d >= %d", textStart, textEnd)
		}

		if debugMode >= 1 {
			fmt.Println()
			fmt.Printf("VERIFY SAFE TO BREAK : start %d end %d text start %d end %d\n", start, end, textStart, textEnd)
			fmt.Println()
		}

		fragment.Clear()
		copyBufferProperties(fragment, buffer)

		flags := fragment.Flags
		if 0 < textStart {
			flags = (flags & ^Bot)
		}
		if textEnd < len(textBuffer.Info) {
			flags = (flags & ^Eot)
		}
		fragment.Flags = flags

		appendBuffer(fragment, textBuffer, textStart, textEnd)
		features, err := so.parseFeatures()
		if err != nil {
			return err
		}
		fragment.Shape(font, features)
		appendBuffer(reconstruction, fragment, 0, len(fragment.Info))

		start = end
		if forward {
			textStart = textEnd
		} else {
			textEnd = textStart
		}
	}

	diff := bufferDiff(reconstruction, buffer, ^GID(0), 0)
	if diff != bufferDiffFlagEqual {
		/* Return the reconstructed result instead so it can be inspected. */
		buffer.Info = nil
		buffer.Pos = nil
		appendBuffer(buffer, reconstruction, 0, len(reconstruction.Info))

		return fmt.Errorf("safe-to-break test failed: %d", diff)
	}

	return nil
}

func (opts *shapeOptions) parseDirection(s string) error {
	switch toLower(s[0]) {
	case 'l':
		opts.props.Direction = LeftToRight
	case 'r':
		opts.props.Direction = RightToLeft
	case 't':
		opts.props.Direction = TopToBottom
	case 'b':
		opts.props.Direction = BottomToTop
	default:
		return fmt.Errorf("invalid direction %s", s)
	}
	return nil
}

// returns the serialized shaped output
// if `verify` is true, additional check on buffer contents is performed
func (mft testOptions) shape(verify bool) (string, error) {
	buffer := mft.shaper.populateBuffer(mft.input)

	font := mft.fontOpts.getFont()
	if err := mft.shaper.shape(font, buffer, verify); err != nil {
		return "", err
	}

	return buffer.serialize(font, mft.format), nil
}

const featuresUsage = `Comma-separated list of font features

    Features can be enabled or disabled, either globally or limited to
    specific character ranges.  The format for specifying feature settings
    follows.  All valid CSS font-feature-settings values other than 'normal'
    and the global values are also accepted, though not documented below.
    CSS string escapes are not supported.

    The range indices refer to the positions between Unicode characters,
    unless the --utf8-clusters is provided, in which case range indices
    refer to UTF-8 byte indices. The position before the first character
    is always 0.

    The format is Python-esque.  Here is how it all works:

      Syntax:       Value:    Start:    End:

    Setting value:
      "kern"        1         0         ∞         // Turn feature on
      "+kern"       1         0         ∞         // Turn feature on
      "-kern"       0         0         ∞         // Turn feature off
      "kern=0"      0         0         ∞         // Turn feature off
      "kern=1"      1         0         ∞         // Turn feature on
      "aalt=2"      2         0         ∞         // Choose 2nd alternate

    Setting index:
      "kern[]"      1         0         ∞         // Turn feature on
      "kern[:]"     1         0         ∞         // Turn feature on
      "kern[5:]"    1         5         ∞         // Turn feature on, partial
      "kern[:5]"    1         0         5         // Turn feature on, partial
      "kern[3:5]"   1         3         5         // Turn feature on, range
      "kern[3]"     1         3         3+1       // Turn feature on, single char

    Mixing it all:

      "aalt[3:5]=2" 2         3         5         // Turn 2nd alternate on for range
`

func (opts *shapeOptions) parseFeatures() ([]Feature, error) {
	if opts.features == "" {
		return nil, nil
	}
	// remove possible quote
	s := strings.Trim(opts.features, `"`)

	features := strings.Split(s, ",")
	out := make([]Feature, len(features))

	var err error
	for i, feature := range features {
		out[i], err = ParseFeature(feature)
		if err != nil {
			return nil, fmt.Errorf("parsing features %s: %s", opts.features, err)
		}
	}
	return out, nil
}

func (opts *fontOptions) parseFontSize(arg string) error {
	if arg == "upem" {
		opts.fontSizeY = fontSizeUpem
		opts.fontSizeX = fontSizeUpem
		return nil
	}
	n, err := fmt.Sscanf(arg, "%d %d", &opts.fontSizeX, &opts.fontSizeY)
	if err != io.EOF {
		return fmt.Errorf("font-size argument should be one or two space-separated numbers")
	}
	if n == 1 {
		opts.fontSizeY = opts.fontSizeX
	}
	return nil
}

func (opts *fontOptions) parseFontPpem(arg string) error {
	n, err := fmt.Sscanf(arg, "%d %d", &opts.xPpem, &opts.yPpem)
	if err != io.EOF {
		return fmt.Errorf("font-ppem argument should be one or two space-separated integers")
	}
	if n == 1 {
		opts.yPpem = opts.xPpem
	}
	return nil
}

const variationsUsage = `Comma-separated list of font variations

    Variations are set globally. The format for specifying variation settings
    follows.  All valid CSS font-variation-settings values other than 'normal'
    and 'inherited' are also accepted, although not documented below.

    The format is a tag, optionally followed by an equals sign, followed by a
    number. For example:

      "wght=500"
      "slnt=-7.5";
`

// parse the options, written in command line format
func parseOptions(options string) (testOptions, error) {
	flags := flag.NewFlagSet("options", flag.ContinueOnError)

	var fmtOpts formatOptions
	flags.BoolVar(&fmtOpts.hideClusters, "no-clusters", false, "Do not output cluster indices")
	flags.BoolVar(&fmtOpts.hideGlyphNames, "no-glyph-names", false, "Output glyph indices instead of names")
	flags.BoolVar(&fmtOpts.hidePositions, "no-positions", false, "Do not output glyph positions")
	flags.BoolVar(&fmtOpts.hideAdvances, "no-advances", false, "Do not output glyph advances")
	flags.BoolVar(&fmtOpts.showExtents, "show-extents", false, "Output glyph extents")
	flags.BoolVar(&fmtOpts.showFlags, "show-flags", false, "Output glyph flags")

	ned := flags.Bool("ned", false, "No Extra Data; Do not output clusters or advances")

	var shapeOpts shapeOptions
	flags.StringVar(&shapeOpts.features, "features", "", featuresUsage)

	flags.String("list-shapers", "", "(ignored)")
	flags.StringVar(&shapeOpts.shaper, "shaper", "", "Force a shaper")
	flags.String("shapers", "", "(ignored)")
	flags.Func("direction", "Set text direction (default: auto)", shapeOpts.parseDirection)
	flags.Func("language", "Set text language (default: $LANG)", func(s string) error {
		shapeOpts.props.Language = language.NewLanguage(s)
		return nil
	})
	flags.Func("script", "Set text script, as an ISO-15924 tag (default: auto)", func(s string) error {
		var err error
		shapeOpts.props.Script, err = language.ParseScript(s)
		return err
	})
	flags.BoolVar(&shapeOpts.bot, "bot", false, "Treat text as beginning-of-paragraph")
	flags.BoolVar(&shapeOpts.eot, "eot", false, "Treat text as end-of-paragraph")
	flags.BoolVar(&shapeOpts.removeDefaultIgnorables, "remove-default-ignorables", false, "Remove Default-Ignorable characters")
	flags.BoolVar(&shapeOpts.preserveDefaultIgnorables, "preserve-default-ignorables", false, "Preserve Default-Ignorable characters")
	flags.Func("cluster-level", "Cluster merging level (0/1/2, default: 0)", func(s string) error {
		l, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid cluster-level option: %s", err)
		}
		if l < 0 || l > 2 {
			return fmt.Errorf("invalid cluster-level option : %d", l)
		}
		shapeOpts.clusterLevel = ClusterLevel(l)
		return nil
	})

	fontOpts := newFontOptions()

	flags.StringVar(&fontOpts.fontRef.File, "font-file", "", "Set font file-name")
	fontRefIndex := flags.Int("face-index", 0, "Set face index (default: 0)")
	flags.Func("font-size", "Font size", fontOpts.parseFontSize)
	flags.Func("font-ppem", "Set x,y pixels per EM (default: 0; disabled)", fontOpts.parseFontPpem)
	flags.Float64Var(&fontOpts.ptem, "font-ptem", 0, "Set font point-size (default: 0; disabled)")
	flags.Func("variations", variationsUsage, fontOpts.parseVariations)
	flags.String("font-funcs", "", "(ignored)")
	flags.String("ft-load-flags", "", "(ignored)")

	ub := flags.String("unicodes-before", "", "Set Unicode codepoints context before each line")
	ua := flags.String("unicodes-after", "", "Set Unicode codepoints context after each line")

	err := flags.Parse(strings.Split(options, " "))
	if err != nil {
		return testOptions{}, err
	}

	if *ned {
		fmtOpts.hideClusters = true
		fmtOpts.hideAdvances = true
	}
	fontOpts.fontRef.Index = uint16(*fontRefIndex)
	out := testOptions{
		fontOpts: fontOpts,
		format:   fmtOpts,
		shaper:   shapeOpts,
	}

	if *ub != "" {
		out.input.textBefore, err = parseUnicodes(*ub)
		if err != nil {
			return testOptions{}, err
		}
	}
	if *ua != "" {
		out.input.textAfter, err = parseUnicodes(*ua)
		if err != nil {
			return testOptions{}, err
		}
	}

	return out, nil
}

// harfbuzz seems to be OK with an invalid font
// in pratice, it seems useless to do shaping without
// font, so we dont support it, meaning we skip this test
func skipInvalidFontIndex(ft api.FontID) bool {
	f, err := td.Files.ReadFile(ft.File)
	check(err)

	fonts, err := loader.NewLoaders(bytes.NewReader(f))
	check(err)

	if int(ft.Index) >= len(fonts) {
		fmt.Printf("skipping invalid font index %d in font %s\n", ft.Index, ft.File)
		return true
	}
	return false
}

type testAction = func(t *testing.T, driver testOptions, dir, line, glyphsExpected string)

// skipVerify is true when debugging, to reduce stdout clutter
func runShapingTest(t *testing.T, driver testOptions, dir, line, glyphsExpected string, skipVerify bool) {
	verify := glyphsExpected != "*"

	// actual does the shaping
	output, err := driver.shape(!skipVerify && verify)
	if err != nil {
		t.Fatal("line ", line, ":", err)
	}

	if got := strings.TrimSpace(output); verify && glyphsExpected != got {
		t.Fatalf("for dir %s and line\n%s\n, expected :\n%s\n got \n%s", dir, line, glyphsExpected, got)
	}
}

// parses and run one test given as line in .tests files
func parseAndRunTest(t *testing.T, dir, line string, action testAction) {
	chunks := strings.Split(line, ";")
	if L := len(chunks); L != 4 {
		t.Fatalf("invalid test line %s : %d chunks", line, L)
	}
	fontFileHash, options, unicodes, glyphsExpected := chunks[0], chunks[1], chunks[2], chunks[3]

	splitHash := strings.Split(fontFileHash, "@")
	fontFile := filepath.Join(dir, splitHash[0])
	if len(splitHash) >= 2 {
		ff, err := td.Files.ReadFile(fontFile)
		check(err)

		hash := sha1.Sum(ff)
		trimmedHash := strings.TrimSpace(hex.EncodeToString(hash[:]))
		if exp := splitHash[1]; trimmedHash != exp {
			t.Fatalf("invalid font file (%s) hash: expected %s, got %s", fontFile, exp, trimmedHash)
		}
	}

	driver, err := parseOptions(options)
	if err != nil {
		t.Fatalf("invalid test file: line %s: %s", line, err)
	}
	driver.fontOpts.fontRef.File = fontFile

	if skipInvalidFontIndex(driver.fontOpts.fontRef) {
		return
	}

	driver.input.text, err = parseUnicodes(unicodes)
	if err != nil {
		t.Fatalf("invalid test file: line %s: %s", line, err)
	}

	action(t, driver, dir, line, glyphsExpected)
}

// opens and parses a test file containing
// the font file, the unicode input and the expected result
func processHarfbuzzTestFile(t *testing.T, dir, filename string, action testAction) {
	f, err := td.Files.ReadFile(filename)
	check(err)

	for _, line := range strings.Split(string(f), "\n") {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" { // skip comments
			continue
		}

		// special case
		if strings.Contains(line, "--shaper=fallback") {
			// we do not support fallback shaper
			continue
		}
		// fails since the FT and Harfbuzz implementations of GlyphVOrigin differ
		// we prefer to match Harfbuzz implementation, so we replace
		// these tests with same, using Harbufzz font funcs
		if line == "../fonts/191826b9643e3f124d865d617ae609db6a2ce203.ttf;--direction=t --font-funcs=ft;U+300C;[uni300C.vert=0@-512,-578+0,-1024]" {
			line = "../fonts/191826b9643e3f124d865d617ae609db6a2ce203.ttf;--direction=t --font-funcs=ot;U+300C;[uni300C.vert=0@-512,-189+0,-1024]"
		} else if line == "../fonts/f9b1dd4dcb515e757789a22cb4241107746fd3d0.ttf;--direction=t --font-funcs=ft;U+0041,U+0042;[gid1=0@-654,-2128+0,-2789|gid2=1@-665,-2125+0,-2789]" {
			line = "../fonts/f9b1dd4dcb515e757789a22cb4241107746fd3d0.ttf;--direction=t --font-funcs=ot;U+0041,U+0042;[gid1=0@-654,-1468+0,-2048|gid2=1@-665,-1462+0,-2048]"
		}

		parseAndRunTest(t, dir, line, action)
	}
}

func dirFiles(t *testing.T, dir string) []string {
	files, err := td.Files.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var filenames []string
	for _, fi := range files {
		filenames = append(filenames, filepath.Join(dir, fi.Name()))
	}
	return filenames
}

func walkShapeTests(t *testing.T, action testAction) {
	disabledTests := []string{
		// requires proprietary fonts from the system (see the file)
		// "harfbuzz_reference/in-house/tests/macos.tests",

		// already handled in emojis_test.go
		"harfbuzz_reference/in-house/tests/emoji-clusters.tests",

		// disabled by harfbuzz (see harfbuzz/test/shaping/data/text-rendering-tests/DISABLED)
		"harfbuzz_reference/text-rendering-tests/tests/CMAP-3.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHARAN-1.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHBALI-1.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHBALI-2.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHKNDA-2.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHKNDA-3.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-1.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-10.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-2.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-3.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-4.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-5.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-6.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-7.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-8.tests",
		"harfbuzz_reference/text-rendering-tests/tests/SHLANA-9.tests",
	}

	isDisabled := func(file string) bool {
		for _, dis := range disabledTests {
			if file == dis {
				return true
			}
		}
		return false
	}

	for _, file := range dirFiles(t, "harfbuzz_reference/aots/tests") {
		if isDisabled(file) {
			continue
		}

		processHarfbuzzTestFile(t, "harfbuzz_reference/aots/tests", file, action)
	}
	for _, file := range dirFiles(t, "harfbuzz_reference/in-house/tests") {
		if isDisabled(file) {
			continue
		}

		processHarfbuzzTestFile(t, "harfbuzz_reference/in-house/tests", file, action)
	}

	for _, file := range dirFiles(t, "harfbuzz_reference/text-rendering-tests/tests") {
		if isDisabled(file) {
			continue
		}

		processHarfbuzzTestFile(t, "harfbuzz_reference/text-rendering-tests/tests", file, action)
	}
}

func runOneTest(t *testing.T, driver testOptions, dir, line, glyphsExpected string) {
	runShapingTest(t, driver, dir, line, glyphsExpected, false)
}

func TestShapeExpected(t *testing.T) {
	walkShapeTests(t, runOneTest)
}

func TestDebug(t *testing.T) {
	dir := "harfbuzz_reference/in-house"
	testString := `macos/System/Library/Fonts/GeezaPro.ttc@ab26ea45dcaa5e1c5a958e42af10e10d330e7334;--font-funcs ot;U+0627,U+0644,U+0623,U+064E,U+0628,U+0652,U+062C,U+064E,U+062F,U+0650,U+064A,U+064E,U+0651,U+0629,U+0640,U+0627,U+0644,U+0639,U+064E,U+0631,U+064E,U+0628,U+0650,U+064A,U+064E,U+0651,U+0629;[u0629.final.tehMarbuta=26+713|u064e_u0651.shaddaFatha=23@0,-200+0|u064a.medial.yeh=23+656|u0650.kasra=21@80,290+80|u0628.initial.beh=21@-80,0+576|u064e.fatha=19@200,-570+200|u0631.final.reh=19@-200,0+702|u064e.fatha=17@200,-200+200|u0639.medial.ain=17@-200,0+738|u0644.initial.lam=16+515|u0627.final.alef=15+647|u0640.tatweel=14+449|u0629.final.tehMarbuta=13+713|u064e_u0651.shaddaFatha=10@0,-200+0|u064a.initial.yeh=10+656|u0650.kasra=8@80,570+80|u062f.final.dal=8@-80,0+822|u064e.fatha=6@290,-160+290|u062c.medial.jeem=6@-290,0+1069|u0652.sukun=4@0,-200+0|u0628.initial.beh=4+656|u064e.fatha=1@-372,120+-372|u0644_u0623.isolated.lamHamzaOnAlef=1@120,0+1282|u0627.alef=0+647]`

	parseAndRunTest(t, dir, testString, func(t *testing.T, driver testOptions, dir, line, glyphsExpected string) {
		runShapingTest(t, driver, dir, line, glyphsExpected, true)
	})
}

func TestExample(t *testing.T) {
	ft := openFontFileTT(t, "common/NotoSansArabic.ttf")
	buffer := NewBuffer()

	// runes := []rune("This is a line to shape..")
	runes := []rune{0x0633, 0x064F, 0x0644, 0x064E, 0x0651, 0x0627, 0x0651, 0x0650, 0x0645, 0x062A, 0x06CC}
	buffer.AddRunes(runes, 0, -1)

	face := &font.Face{Font: ft}
	font := NewFont(face)
	buffer.GuessSegmentProperties()
	buffer.Shape(font, nil)

	for i, pos := range buffer.Pos {
		info := buffer.Info[i]
		ext, ok := face.GlyphExtents(info.Glyph)
		testutils.AssertC(t, ok, fmt.Sprintf("invalid glyph %d", info.Glyph))

		fmt.Println(pos.XAdvance, pos.XOffset, ext.Width, ext.XBearing)
	}
}
