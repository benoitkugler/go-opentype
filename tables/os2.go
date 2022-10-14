package tables

//go:generate ../../binarygen/cmd/generator os2.go

// OS/2 and Windows Metrics Table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/os2
type Os2 struct {
	version             uint16
	xAvgCharWidth       uint16
	uSWeightClass       uint16
	uSWidthClass        uint16
	fSType              uint16
	ySubscriptXSize     int16
	ySubscriptYSize     int16
	ySubscriptXOffset   int16
	ySubscriptYOffset   int16
	ySuperscriptXSize   int16
	ySuperscriptYSize   int16
	ySuperscriptXOffset int16
	ySuperscriptYOffset int16
	yStrikeoutSize      int16
	yStrikeoutPosition  int16
	sFamilyClass        int16
	panose              [10]byte
	ulCharRange         [4]uint32
	achVendID           Tag
	fsSelection         uint16
	uSFirstCharIndex    uint16
	uSLastCharIndex     uint16
	sTypoAscender       int16
	sTypoDescender      int16
	sTypoLineGap        int16
	usWinAscent         uint16
	usWinDescent        uint16
	higherVersionData   []byte `len:"__toEnd"`
}
