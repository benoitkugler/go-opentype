package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from cmap_src.go. DO NOT EDIT

func (item *CmapSubtable0) mustParse(src []byte) {
	_ = src[261] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	item.length = binary.BigEndian.Uint16(src[2:])
	item.language = binary.BigEndian.Uint16(src[4:])
	item.glyphIdArray[0] = src[6]
	item.glyphIdArray[1] = src[7]
	item.glyphIdArray[2] = src[8]
	item.glyphIdArray[3] = src[9]
	item.glyphIdArray[4] = src[10]
	item.glyphIdArray[5] = src[11]
	item.glyphIdArray[6] = src[12]
	item.glyphIdArray[7] = src[13]
	item.glyphIdArray[8] = src[14]
	item.glyphIdArray[9] = src[15]
	item.glyphIdArray[10] = src[16]
	item.glyphIdArray[11] = src[17]
	item.glyphIdArray[12] = src[18]
	item.glyphIdArray[13] = src[19]
	item.glyphIdArray[14] = src[20]
	item.glyphIdArray[15] = src[21]
	item.glyphIdArray[16] = src[22]
	item.glyphIdArray[17] = src[23]
	item.glyphIdArray[18] = src[24]
	item.glyphIdArray[19] = src[25]
	item.glyphIdArray[20] = src[26]
	item.glyphIdArray[21] = src[27]
	item.glyphIdArray[22] = src[28]
	item.glyphIdArray[23] = src[29]
	item.glyphIdArray[24] = src[30]
	item.glyphIdArray[25] = src[31]
	item.glyphIdArray[26] = src[32]
	item.glyphIdArray[27] = src[33]
	item.glyphIdArray[28] = src[34]
	item.glyphIdArray[29] = src[35]
	item.glyphIdArray[30] = src[36]
	item.glyphIdArray[31] = src[37]
	item.glyphIdArray[32] = src[38]
	item.glyphIdArray[33] = src[39]
	item.glyphIdArray[34] = src[40]
	item.glyphIdArray[35] = src[41]
	item.glyphIdArray[36] = src[42]
	item.glyphIdArray[37] = src[43]
	item.glyphIdArray[38] = src[44]
	item.glyphIdArray[39] = src[45]
	item.glyphIdArray[40] = src[46]
	item.glyphIdArray[41] = src[47]
	item.glyphIdArray[42] = src[48]
	item.glyphIdArray[43] = src[49]
	item.glyphIdArray[44] = src[50]
	item.glyphIdArray[45] = src[51]
	item.glyphIdArray[46] = src[52]
	item.glyphIdArray[47] = src[53]
	item.glyphIdArray[48] = src[54]
	item.glyphIdArray[49] = src[55]
	item.glyphIdArray[50] = src[56]
	item.glyphIdArray[51] = src[57]
	item.glyphIdArray[52] = src[58]
	item.glyphIdArray[53] = src[59]
	item.glyphIdArray[54] = src[60]
	item.glyphIdArray[55] = src[61]
	item.glyphIdArray[56] = src[62]
	item.glyphIdArray[57] = src[63]
	item.glyphIdArray[58] = src[64]
	item.glyphIdArray[59] = src[65]
	item.glyphIdArray[60] = src[66]
	item.glyphIdArray[61] = src[67]
	item.glyphIdArray[62] = src[68]
	item.glyphIdArray[63] = src[69]
	item.glyphIdArray[64] = src[70]
	item.glyphIdArray[65] = src[71]
	item.glyphIdArray[66] = src[72]
	item.glyphIdArray[67] = src[73]
	item.glyphIdArray[68] = src[74]
	item.glyphIdArray[69] = src[75]
	item.glyphIdArray[70] = src[76]
	item.glyphIdArray[71] = src[77]
	item.glyphIdArray[72] = src[78]
	item.glyphIdArray[73] = src[79]
	item.glyphIdArray[74] = src[80]
	item.glyphIdArray[75] = src[81]
	item.glyphIdArray[76] = src[82]
	item.glyphIdArray[77] = src[83]
	item.glyphIdArray[78] = src[84]
	item.glyphIdArray[79] = src[85]
	item.glyphIdArray[80] = src[86]
	item.glyphIdArray[81] = src[87]
	item.glyphIdArray[82] = src[88]
	item.glyphIdArray[83] = src[89]
	item.glyphIdArray[84] = src[90]
	item.glyphIdArray[85] = src[91]
	item.glyphIdArray[86] = src[92]
	item.glyphIdArray[87] = src[93]
	item.glyphIdArray[88] = src[94]
	item.glyphIdArray[89] = src[95]
	item.glyphIdArray[90] = src[96]
	item.glyphIdArray[91] = src[97]
	item.glyphIdArray[92] = src[98]
	item.glyphIdArray[93] = src[99]
	item.glyphIdArray[94] = src[100]
	item.glyphIdArray[95] = src[101]
	item.glyphIdArray[96] = src[102]
	item.glyphIdArray[97] = src[103]
	item.glyphIdArray[98] = src[104]
	item.glyphIdArray[99] = src[105]
	item.glyphIdArray[100] = src[106]
	item.glyphIdArray[101] = src[107]
	item.glyphIdArray[102] = src[108]
	item.glyphIdArray[103] = src[109]
	item.glyphIdArray[104] = src[110]
	item.glyphIdArray[105] = src[111]
	item.glyphIdArray[106] = src[112]
	item.glyphIdArray[107] = src[113]
	item.glyphIdArray[108] = src[114]
	item.glyphIdArray[109] = src[115]
	item.glyphIdArray[110] = src[116]
	item.glyphIdArray[111] = src[117]
	item.glyphIdArray[112] = src[118]
	item.glyphIdArray[113] = src[119]
	item.glyphIdArray[114] = src[120]
	item.glyphIdArray[115] = src[121]
	item.glyphIdArray[116] = src[122]
	item.glyphIdArray[117] = src[123]
	item.glyphIdArray[118] = src[124]
	item.glyphIdArray[119] = src[125]
	item.glyphIdArray[120] = src[126]
	item.glyphIdArray[121] = src[127]
	item.glyphIdArray[122] = src[128]
	item.glyphIdArray[123] = src[129]
	item.glyphIdArray[124] = src[130]
	item.glyphIdArray[125] = src[131]
	item.glyphIdArray[126] = src[132]
	item.glyphIdArray[127] = src[133]
	item.glyphIdArray[128] = src[134]
	item.glyphIdArray[129] = src[135]
	item.glyphIdArray[130] = src[136]
	item.glyphIdArray[131] = src[137]
	item.glyphIdArray[132] = src[138]
	item.glyphIdArray[133] = src[139]
	item.glyphIdArray[134] = src[140]
	item.glyphIdArray[135] = src[141]
	item.glyphIdArray[136] = src[142]
	item.glyphIdArray[137] = src[143]
	item.glyphIdArray[138] = src[144]
	item.glyphIdArray[139] = src[145]
	item.glyphIdArray[140] = src[146]
	item.glyphIdArray[141] = src[147]
	item.glyphIdArray[142] = src[148]
	item.glyphIdArray[143] = src[149]
	item.glyphIdArray[144] = src[150]
	item.glyphIdArray[145] = src[151]
	item.glyphIdArray[146] = src[152]
	item.glyphIdArray[147] = src[153]
	item.glyphIdArray[148] = src[154]
	item.glyphIdArray[149] = src[155]
	item.glyphIdArray[150] = src[156]
	item.glyphIdArray[151] = src[157]
	item.glyphIdArray[152] = src[158]
	item.glyphIdArray[153] = src[159]
	item.glyphIdArray[154] = src[160]
	item.glyphIdArray[155] = src[161]
	item.glyphIdArray[156] = src[162]
	item.glyphIdArray[157] = src[163]
	item.glyphIdArray[158] = src[164]
	item.glyphIdArray[159] = src[165]
	item.glyphIdArray[160] = src[166]
	item.glyphIdArray[161] = src[167]
	item.glyphIdArray[162] = src[168]
	item.glyphIdArray[163] = src[169]
	item.glyphIdArray[164] = src[170]
	item.glyphIdArray[165] = src[171]
	item.glyphIdArray[166] = src[172]
	item.glyphIdArray[167] = src[173]
	item.glyphIdArray[168] = src[174]
	item.glyphIdArray[169] = src[175]
	item.glyphIdArray[170] = src[176]
	item.glyphIdArray[171] = src[177]
	item.glyphIdArray[172] = src[178]
	item.glyphIdArray[173] = src[179]
	item.glyphIdArray[174] = src[180]
	item.glyphIdArray[175] = src[181]
	item.glyphIdArray[176] = src[182]
	item.glyphIdArray[177] = src[183]
	item.glyphIdArray[178] = src[184]
	item.glyphIdArray[179] = src[185]
	item.glyphIdArray[180] = src[186]
	item.glyphIdArray[181] = src[187]
	item.glyphIdArray[182] = src[188]
	item.glyphIdArray[183] = src[189]
	item.glyphIdArray[184] = src[190]
	item.glyphIdArray[185] = src[191]
	item.glyphIdArray[186] = src[192]
	item.glyphIdArray[187] = src[193]
	item.glyphIdArray[188] = src[194]
	item.glyphIdArray[189] = src[195]
	item.glyphIdArray[190] = src[196]
	item.glyphIdArray[191] = src[197]
	item.glyphIdArray[192] = src[198]
	item.glyphIdArray[193] = src[199]
	item.glyphIdArray[194] = src[200]
	item.glyphIdArray[195] = src[201]
	item.glyphIdArray[196] = src[202]
	item.glyphIdArray[197] = src[203]
	item.glyphIdArray[198] = src[204]
	item.glyphIdArray[199] = src[205]
	item.glyphIdArray[200] = src[206]
	item.glyphIdArray[201] = src[207]
	item.glyphIdArray[202] = src[208]
	item.glyphIdArray[203] = src[209]
	item.glyphIdArray[204] = src[210]
	item.glyphIdArray[205] = src[211]
	item.glyphIdArray[206] = src[212]
	item.glyphIdArray[207] = src[213]
	item.glyphIdArray[208] = src[214]
	item.glyphIdArray[209] = src[215]
	item.glyphIdArray[210] = src[216]
	item.glyphIdArray[211] = src[217]
	item.glyphIdArray[212] = src[218]
	item.glyphIdArray[213] = src[219]
	item.glyphIdArray[214] = src[220]
	item.glyphIdArray[215] = src[221]
	item.glyphIdArray[216] = src[222]
	item.glyphIdArray[217] = src[223]
	item.glyphIdArray[218] = src[224]
	item.glyphIdArray[219] = src[225]
	item.glyphIdArray[220] = src[226]
	item.glyphIdArray[221] = src[227]
	item.glyphIdArray[222] = src[228]
	item.glyphIdArray[223] = src[229]
	item.glyphIdArray[224] = src[230]
	item.glyphIdArray[225] = src[231]
	item.glyphIdArray[226] = src[232]
	item.glyphIdArray[227] = src[233]
	item.glyphIdArray[228] = src[234]
	item.glyphIdArray[229] = src[235]
	item.glyphIdArray[230] = src[236]
	item.glyphIdArray[231] = src[237]
	item.glyphIdArray[232] = src[238]
	item.glyphIdArray[233] = src[239]
	item.glyphIdArray[234] = src[240]
	item.glyphIdArray[235] = src[241]
	item.glyphIdArray[236] = src[242]
	item.glyphIdArray[237] = src[243]
	item.glyphIdArray[238] = src[244]
	item.glyphIdArray[239] = src[245]
	item.glyphIdArray[240] = src[246]
	item.glyphIdArray[241] = src[247]
	item.glyphIdArray[242] = src[248]
	item.glyphIdArray[243] = src[249]
	item.glyphIdArray[244] = src[250]
	item.glyphIdArray[245] = src[251]
	item.glyphIdArray[246] = src[252]
	item.glyphIdArray[247] = src[253]
	item.glyphIdArray[248] = src[254]
	item.glyphIdArray[249] = src[255]
	item.glyphIdArray[250] = src[256]
	item.glyphIdArray[251] = src[257]
	item.glyphIdArray[252] = src[258]
	item.glyphIdArray[253] = src[259]
	item.glyphIdArray[254] = src[260]
	item.glyphIdArray[255] = src[261]
}

func ParseCmap(src []byte) (Cmap, int, error) {
	var item Cmap
	n := 0
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading Cmap: "+"EOF: expected length: 4, got %d", L)
		}
		_ = src[3] // early bound checking
		item.version = binary.BigEndian.Uint16(src[0:])
		item.numTables = binary.BigEndian.Uint16(src[2:])
		n += 4
	}
	{
		arrayLength := int(item.numTables)

		if L := len(src); L < 4+arrayLength*8 {
			return item, 0, fmt.Errorf("reading Cmap: "+"EOF: expected length: %d, got %d", 4+arrayLength*8, L)
		}

		item.records = make([]encodingRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.records {
			item.records[i].mustParse(src[4+i*8:])
		}
		n += arrayLength * 8
	}
	{

		read, err := item.customParseSubtables(src[:])
		if err != nil {
			return item, 0, fmt.Errorf("reading Cmap: %s", err)
		}
		n = read
	}
	return item, n, nil
}

func ParseCmapSubtable0(src []byte) (CmapSubtable0, int, error) {
	var item CmapSubtable0
	n := 0
	if L := len(src); L < 262 {
		return item, 0, fmt.Errorf("reading CmapSubtable0: "+"EOF: expected length: 262, got %d", L)
	}
	item.mustParse(src)
	n += 262
	return item, n, nil
}

func ParseCmapSubtable10(src []byte) (CmapSubtable10, int, error) {
	var item CmapSubtable10
	n := 0
	{
		if L := len(src); L < 20 {
			return item, 0, fmt.Errorf("reading CmapSubtable10: "+"EOF: expected length: 20, got %d", L)
		}
		_ = src[19] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.reserved = binary.BigEndian.Uint16(src[2:])
		item.length = binary.BigEndian.Uint32(src[4:])
		item.language = binary.BigEndian.Uint32(src[8:])
		item.startCharCode = binary.BigEndian.Uint32(src[12:])
		item.numChars = binary.BigEndian.Uint32(src[16:])
		n += 20
	}
	{
		arrayLength := int(item.numChars)

		if L := len(src); L < 20+arrayLength*2 {
			return item, 0, fmt.Errorf("reading CmapSubtable10: "+"EOF: expected length: %d, got %d", 20+arrayLength*2, L)
		}

		item.glyphIdArray = make([]GlyphID, arrayLength) // allocation guarded by the previous check
		for i := range item.glyphIdArray {
			item.glyphIdArray[i] = GlyphID(binary.BigEndian.Uint16(src[20+i*2:]))
		}
		n += arrayLength * 2
	}
	return item, n, nil
}

func ParseCmapSubtable12(src []byte) (CmapSubtable12, int, error) {
	var item CmapSubtable12
	n := 0
	{
		if L := len(src); L < 16 {
			return item, 0, fmt.Errorf("reading CmapSubtable12: "+"EOF: expected length: 16, got %d", L)
		}
		_ = src[15] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.reserved = binary.BigEndian.Uint16(src[2:])
		item.length = binary.BigEndian.Uint32(src[4:])
		item.language = binary.BigEndian.Uint32(src[8:])
		item.numGroups = binary.BigEndian.Uint32(src[12:])
		n += 16
	}
	{
		arrayLength := int(item.numGroups)

		if L := len(src); L < 16+arrayLength*12 {
			return item, 0, fmt.Errorf("reading CmapSubtable12: "+"EOF: expected length: %d, got %d", 16+arrayLength*12, L)
		}

		item.groups = make([]sequentialMapGroup, arrayLength) // allocation guarded by the previous check
		for i := range item.groups {
			item.groups[i].mustParse(src[16+i*12:])
		}
		n += arrayLength * 12
	}
	return item, n, nil
}

func ParseCmapSubtable13(src []byte) (CmapSubtable13, int, error) {
	var item CmapSubtable13
	n := 0
	{
		if L := len(src); L < 16 {
			return item, 0, fmt.Errorf("reading CmapSubtable13: "+"EOF: expected length: 16, got %d", L)
		}
		_ = src[15] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.reserved = binary.BigEndian.Uint16(src[2:])
		item.length = binary.BigEndian.Uint32(src[4:])
		item.language = binary.BigEndian.Uint32(src[8:])
		item.numGroups = binary.BigEndian.Uint32(src[12:])
		n += 16
	}
	{
		arrayLength := int(item.numGroups)

		if L := len(src); L < 16+arrayLength*12 {
			return item, 0, fmt.Errorf("reading CmapSubtable13: "+"EOF: expected length: %d, got %d", 16+arrayLength*12, L)
		}

		item.groups = make([]sequentialMapGroup, arrayLength) // allocation guarded by the previous check
		for i := range item.groups {
			item.groups[i].mustParse(src[16+i*12:])
		}
		n += arrayLength * 12
	}
	return item, n, nil
}

func ParseCmapSubtable14(src []byte) (CmapSubtable14, int, error) {
	var item CmapSubtable14
	n := 0
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading CmapSubtable14: "+"EOF: expected length: 10, got %d", L)
		}
		_ = src[9] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.length = binary.BigEndian.Uint32(src[2:])
		item.numVarSelectorRecords = binary.BigEndian.Uint32(src[6:])
		n += 10
	}
	{
		arrayLength := int(item.numVarSelectorRecords)

		if L := len(src); L < 10+arrayLength*11 {
			return item, 0, fmt.Errorf("reading CmapSubtable14: "+"EOF: expected length: %d, got %d", 10+arrayLength*11, L)
		}

		item.varSelectors = make([]variationSelector, arrayLength) // allocation guarded by the previous check
		for i := range item.varSelectors {
			item.varSelectors[i].mustParse(src[10+i*11:])
		}
		n += arrayLength * 11
	}
	{

		item.rawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func ParseCmapSubtable2(src []byte) (CmapSubtable2, int, error) {
	var item CmapSubtable2
	n := 0
	{
		if L := len(src); L < 2 {
			return item, 0, fmt.Errorf("reading CmapSubtable2: "+"EOF: expected length: 2, got %d", L)
		}
		item.format = binary.BigEndian.Uint16(src[0:])
		n += 2
	}
	{

		item.rawData = src[2:]
		n = len(src)
	}
	return item, n, nil
}

func ParseCmapSubtable4(src []byte) (CmapSubtable4, int, error) {
	var item CmapSubtable4
	n := 0
	{
		if L := len(src); L < 14 {
			return item, 0, fmt.Errorf("reading CmapSubtable4: "+"EOF: expected length: 14, got %d", L)
		}
		_ = src[13] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.length = binary.BigEndian.Uint16(src[2:])
		item.language = binary.BigEndian.Uint16(src[4:])
		item.segCountX2 = binary.BigEndian.Uint16(src[6:])
		item.searchRange = binary.BigEndian.Uint16(src[8:])
		item.entrySelector = binary.BigEndian.Uint16(src[10:])
		item.rangeShift = binary.BigEndian.Uint16(src[12:])
		n += 14
	}
	{
		arrayLength := int(item.segCountX2 / 2)

		if L := len(src); L < 14+arrayLength*2 {
			return item, 0, fmt.Errorf("reading CmapSubtable4: "+"EOF: expected length: %d, got %d", 14+arrayLength*2, L)
		}

		item.endCode = make([]uint16, arrayLength) // allocation guarded by the previous check
		for i := range item.endCode {
			item.endCode[i] = binary.BigEndian.Uint16(src[14+i*2:])
		}
		n += arrayLength * 2
	}
	{
		if L := len(src); L < n+2 {
			return item, 0, fmt.Errorf("reading CmapSubtable4: "+"EOF: expected length: n + 2, got %d", L)
		}
		item.reservedPad = binary.BigEndian.Uint16(src[n:])
		n += 2
	}
	{
		arrayLength := int(item.segCountX2 / 2)

		if L := len(src); L < n+arrayLength*2 {
			return item, 0, fmt.Errorf("reading CmapSubtable4: "+"EOF: expected length: %d, got %d", n+arrayLength*2, L)
		}

		item.startCode = make([]uint16, arrayLength) // allocation guarded by the previous check
		for i := range item.startCode {
			item.startCode[i] = binary.BigEndian.Uint16(src[n+i*2:])
		}
		n += arrayLength * 2
	}
	{
		arrayLength := int(item.segCountX2 / 2)

		if L := len(src); L < n+arrayLength*2 {
			return item, 0, fmt.Errorf("reading CmapSubtable4: "+"EOF: expected length: %d, got %d", n+arrayLength*2, L)
		}

		item.idDelta = make([]int16, arrayLength) // allocation guarded by the previous check
		for i := range item.idDelta {
			item.idDelta[i] = int16(binary.BigEndian.Uint16(src[n+i*2:]))
		}
		n += arrayLength * 2
	}
	{
		arrayLength := int(item.segCountX2 / 2)

		if L := len(src); L < n+arrayLength*2 {
			return item, 0, fmt.Errorf("reading CmapSubtable4: "+"EOF: expected length: %d, got %d", n+arrayLength*2, L)
		}

		item.idRangeOffsets = make([]uint16, arrayLength) // allocation guarded by the previous check
		for i := range item.idRangeOffsets {
			item.idRangeOffsets[i] = binary.BigEndian.Uint16(src[n+i*2:])
		}
		n += arrayLength * 2
	}
	{

		item.glyphIDArray = src[n:]
		n = len(src)
	}
	return item, n, nil
}

func ParseCmapSubtable6(src []byte) (CmapSubtable6, int, error) {
	var item CmapSubtable6
	n := 0
	{
		if L := len(src); L < 10 {
			return item, 0, fmt.Errorf("reading CmapSubtable6: "+"EOF: expected length: 10, got %d", L)
		}
		_ = src[9] // early bound checking
		item.format = binary.BigEndian.Uint16(src[0:])
		item.length = binary.BigEndian.Uint16(src[2:])
		item.language = binary.BigEndian.Uint16(src[4:])
		item.firstCode = binary.BigEndian.Uint16(src[6:])
		item.entryCount = binary.BigEndian.Uint16(src[8:])
		n += 10
	}
	{
		arrayLength := int(item.entryCount)

		if L := len(src); L < 10+arrayLength*2 {
			return item, 0, fmt.Errorf("reading CmapSubtable6: "+"EOF: expected length: %d, got %d", 10+arrayLength*2, L)
		}

		item.glyphIdArray = make([]GlyphID, arrayLength) // allocation guarded by the previous check
		for i := range item.glyphIdArray {
			item.glyphIdArray[i] = GlyphID(binary.BigEndian.Uint16(src[10+i*2:]))
		}
		n += arrayLength * 2
	}
	return item, n, nil
}

func ParseDefaultUVSTable(src []byte) (DefaultUVSTable, int, error) {
	var item DefaultUVSTable
	n := 0
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading DefaultUVSTable: "+"EOF: expected length: 4, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint32(src[0:]))
		n += 4

		if L := len(src); L < 4+arrayLength*4 {
			return item, 0, fmt.Errorf("reading DefaultUVSTable: "+"EOF: expected length: %d, got %d", 4+arrayLength*4, L)
		}

		item.ranges = make([]unicodeRange, arrayLength) // allocation guarded by the previous check
		for i := range item.ranges {
			item.ranges[i].mustParse(src[4+i*4:])
		}
		n += arrayLength * 4
	}
	return item, n, nil
}

func ParseUVSMappingTable(src []byte) (UVSMappingTable, int, error) {
	var item UVSMappingTable
	n := 0
	{
		if L := len(src); L < 4 {
			return item, 0, fmt.Errorf("reading UVSMappingTable: "+"EOF: expected length: 4, got %d", L)
		}
		arrayLength := int(binary.BigEndian.Uint32(src[0:]))
		n += 4

		if L := len(src); L < 4+arrayLength*5 {
			return item, 0, fmt.Errorf("reading UVSMappingTable: "+"EOF: expected length: %d, got %d", 4+arrayLength*5, L)
		}

		item.ranges = make([]uvsMappingRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.ranges {
			item.ranges[i].mustParse(src[4+i*5:])
		}
		n += arrayLength * 5
	}
	return item, n, nil
}

func (item *encodingRecord) mustParse(src []byte) {
	_ = src[7] // early bound checking
	item.platformID = binary.BigEndian.Uint16(src[0:])
	item.encodingID = binary.BigEndian.Uint16(src[2:])
	item.subtableOffset = binary.BigEndian.Uint32(src[4:])
}

func (item *sequentialMapGroup) mustParse(src []byte) {
	_ = src[11] // early bound checking
	item.startCharCode = binary.BigEndian.Uint32(src[0:])
	item.endCharCode = binary.BigEndian.Uint32(src[4:])
	item.startGlyphID = binary.BigEndian.Uint32(src[8:])
}

func (item *unicodeRange) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.startUnicodeValue[0] = src[0]
	item.startUnicodeValue[1] = src[1]
	item.startUnicodeValue[2] = src[2]
	item.additionalCount = src[3]
}

func (item *uvsMappingRecord) mustParse(src []byte) {
	_ = src[4] // early bound checking
	item.unicodeValue[0] = src[0]
	item.unicodeValue[1] = src[1]
	item.unicodeValue[2] = src[2]
	item.glyphID = binary.BigEndian.Uint16(src[3:])
}

func (item *variationSelector) mustParse(src []byte) {
	_ = src[10] // early bound checking
	item.varSelector[0] = src[0]
	item.varSelector[1] = src[1]
	item.varSelector[2] = src[2]
	item.defaultUVSOffset = binary.BigEndian.Uint32(src[3:])
	item.nonDefaultUVSOffset = binary.BigEndian.Uint32(src[7:])
}
