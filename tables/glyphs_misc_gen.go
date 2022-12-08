package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from glyphs_misc_src.go. DO NOT EDIT

func ParseSVG(src []byte) (SVG, int, error) {
	var item SVG
	n := 0
	if L := len(src); L < 10 {
		return item, 0, fmt.Errorf("reading SVG: "+"EOF: expected length: 10, got %d", L)
	}
	_ = src[9] // early bound checking
	item.version = binary.BigEndian.Uint16(src[0:])
	offsetSVGDocumentList := int(binary.BigEndian.Uint32(src[2:]))
	item.reserved = binary.BigEndian.Uint32(src[6:])
	n += 10

	{

		if offsetSVGDocumentList != 0 { // ignore null offset
			if L := len(src); L < offsetSVGDocumentList {
				return item, 0, fmt.Errorf("reading SVG: "+"EOF: expected length: %d, got %d", offsetSVGDocumentList, L)
			}

			var err error
			item.SVGDocumentList, _, err = ParseSVGDocumentList(src[offsetSVGDocumentList:])
			if err != nil {
				return item, 0, fmt.Errorf("reading SVG: %s", err)
			}

		}
	}
	return item, n, nil
}

func ParseSVGDocumentList(src []byte) (SVGDocumentList, int, error) {
	var item SVGDocumentList
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading SVGDocumentList: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthDocumentRecords := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthDocumentRecords*12 {
			return item, 0, fmt.Errorf("reading SVGDocumentList: "+"EOF: expected length: %d, got %d", 2+arrayLengthDocumentRecords*12, L)
		}

		item.DocumentRecords = make([]SVGDocumentRecord, arrayLengthDocumentRecords) // allocation guarded by the previous check
		for i := range item.DocumentRecords {
			item.DocumentRecords[i].mustParse(src[2+i*12:])
		}
		n += arrayLengthDocumentRecords * 12
	}
	{

		item.SVGRawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func ParseVORG(src []byte) (VORG, int, error) {
	var item VORG
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading VORG: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.majorVersion = binary.BigEndian.Uint16(src[0:])
	item.minorVersion = binary.BigEndian.Uint16(src[2:])
	item.DefaultVertOriginY = int16(binary.BigEndian.Uint16(src[4:]))
	arrayLengthVertOriginYMetrics := int(binary.BigEndian.Uint16(src[6:]))
	n += 8

	{

		if L := len(src); L < 8+arrayLengthVertOriginYMetrics*4 {
			return item, 0, fmt.Errorf("reading VORG: "+"EOF: expected length: %d, got %d", 8+arrayLengthVertOriginYMetrics*4, L)
		}

		item.VertOriginYMetrics = make([]VertOriginYMetric, arrayLengthVertOriginYMetrics) // allocation guarded by the previous check
		for i := range item.VertOriginYMetrics {
			item.VertOriginYMetrics[i].mustParse(src[8+i*4:])
		}
		n += arrayLengthVertOriginYMetrics * 4
	}
	return item, n, nil
}

func (item *SVGDocumentRecord) mustParse(src []byte) {
	_ = src[11] // early bound checking
	item.StartGlyphID = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.EndGlyphID = GlyphID(binary.BigEndian.Uint16(src[2:]))
	item.SvgDocOffset = Offset32(binary.BigEndian.Uint32(src[4:]))
	item.SvgDocLength = binary.BigEndian.Uint32(src[8:])
}

func (item *VertOriginYMetric) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.GlyphIndex = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.VertOriginY = int16(binary.BigEndian.Uint16(src[2:]))
}
