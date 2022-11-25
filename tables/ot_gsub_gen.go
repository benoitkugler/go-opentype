package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from ot_gsub_src.go. DO NOT EDIT

func ParseAlternateSet(src []byte) (AlternateSet, int, error) {
	var item AlternateSet
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading AlternateSet: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthItemAlternateGlyphIDs := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthItemAlternateGlyphIDs*2 {
			return item, 0, fmt.Errorf("reading AlternateSet: "+"EOF: expected length: %d, got %d", 2+arrayLengthItemAlternateGlyphIDs*2, L)
		}

		item.AlternateGlyphIDs = make([]GlyphID, arrayLengthItemAlternateGlyphIDs) // allocation guarded by the previous check
		for i := range item.AlternateGlyphIDs {
			item.AlternateGlyphIDs[i] = GlyphID(binary.BigEndian.Uint16(src[2+i*2:]))
		}
		n += arrayLengthItemAlternateGlyphIDs * 2
	}
	return item, n, nil
}

func ParseAlternateSubs(src []byte) (AlternateSubs, int, error) {
	var item AlternateSubs
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading AlternateSubs: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.substFormat = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverageOffset := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemAlternateSets := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemAlternateSets*2 {
			return item, 0, fmt.Errorf("reading AlternateSubs: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemAlternateSets*2, L)
		}

		item.AlternateSets = make([]AlternateSet, arrayLengthItemAlternateSets) // allocation guarded by the previous check
		for i := range item.AlternateSets {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading AlternateSubs: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.AlternateSets[i], _, err = ParseAlternateSet(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading AlternateSubs: %s", err)
			}

		}
		n += arrayLengthItemAlternateSets * 2
	}
	{

		if offsetItemCoverageOffset != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverageOffset {
				return item, 0, fmt.Errorf("reading AlternateSubs: "+"EOF: expected length: %d, got %d", offsetItemCoverageOffset, L)
			}

			var (
				err  error
				read int
			)
			item.CoverageOffset, read, err = ParseCoverage(src[offsetItemCoverageOffset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading AlternateSubs: %s", err)
			}
			offsetItemCoverageOffset += read
		}
	}
	return item, n, nil
}

func ParseChainedContextualSubs(src []byte) (ChainedContextualSubs, int, error) {
	var item ChainedContextualSubs
	n := 0
	{
		var (
			err  error
			read int
		)
		item.Data, read, err = ParseChainedContextualSubsITF(src[0:])
		if err != nil {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func ParseChainedContextualSubs1(src []byte) (ChainedContextualSubs1, int, error) {
	var item ChainedContextualSubs1
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubs1: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemChainedSeqRuleSet := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemChainedSeqRuleSet*2 {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs1: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemChainedSeqRuleSet*2, L)
		}

		item.ChainedSeqRuleSet = make([]ChainedSequenceRuleSet, arrayLengthItemChainedSeqRuleSet) // allocation guarded by the previous check
		for i := range item.ChainedSeqRuleSet {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs1: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.ChainedSeqRuleSet[i], _, err = ParseChainedSequenceRuleSet(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs1: %s", err)
			}

		}
		n += arrayLengthItemChainedSeqRuleSet * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs1: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs1: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	return item, n, nil
}

func ParseChainedContextualSubs2(src []byte) (ChainedContextualSubs2, int, error) {
	var item ChainedContextualSubs2
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	offsetItemBacktrackClassDef := int(binary.BigEndian.Uint16(src[4:]))
	offsetItemInputClassDef := int(binary.BigEndian.Uint16(src[6:]))
	offsetItemLookaheadClassDef := int(binary.BigEndian.Uint16(src[8:]))
	arrayLengthItemChainedClassSeqRuleSet := int(binary.BigEndian.Uint16(src[10:]))
	n += 12

	{

		if L := len(src); L < 12+arrayLengthItemChainedClassSeqRuleSet*2 {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: %d, got %d", 12+arrayLengthItemChainedClassSeqRuleSet*2, L)
		}

		item.ChainedClassSeqRuleSet = make([]ChainedClassSequenceRuleSet, arrayLengthItemChainedClassSeqRuleSet) // allocation guarded by the previous check
		for i := range item.ChainedClassSeqRuleSet {
			offset := int(binary.BigEndian.Uint16(src[12+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.ChainedClassSeqRuleSet[i], _, err = ParseChainedClassSequenceRuleSet(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: %s", err)
			}

		}
		n += arrayLengthItemChainedClassSeqRuleSet * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	{

		if offsetItemBacktrackClassDef != 0 { // ignore null offset
			if L := len(src); L < offsetItemBacktrackClassDef {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: %d, got %d", offsetItemBacktrackClassDef, L)
			}

			var (
				err  error
				read int
			)
			item.BacktrackClassDef, read, err = ParseClassDef(src[offsetItemBacktrackClassDef:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: %s", err)
			}
			offsetItemBacktrackClassDef += read
		}
	}
	{

		if offsetItemInputClassDef != 0 { // ignore null offset
			if L := len(src); L < offsetItemInputClassDef {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: %d, got %d", offsetItemInputClassDef, L)
			}

			var (
				err  error
				read int
			)
			item.InputClassDef, read, err = ParseClassDef(src[offsetItemInputClassDef:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: %s", err)
			}
			offsetItemInputClassDef += read
		}
	}
	{

		if offsetItemLookaheadClassDef != 0 { // ignore null offset
			if L := len(src); L < offsetItemLookaheadClassDef {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: "+"EOF: expected length: %d, got %d", offsetItemLookaheadClassDef, L)
			}

			var (
				err  error
				read int
			)
			item.LookaheadClassDef, read, err = ParseClassDef(src[offsetItemLookaheadClassDef:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs2: %s", err)
			}
			offsetItemLookaheadClassDef += read
		}
	}
	return item, n, nil
}

func ParseChainedContextualSubs3(src []byte) (ChainedContextualSubs3, int, error) {
	var item ChainedContextualSubs3
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: 4, got %d", L)
	}
	_ = src[3] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	arrayLengthItemBacktrackCoverages := int(binary.BigEndian.Uint16(src[2:]))
	n += 4

	{

		if L := len(src); L < 4+arrayLengthItemBacktrackCoverages*2 {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", 4+arrayLengthItemBacktrackCoverages*2, L)
		}

		item.BacktrackCoverages = make([]Coverage, arrayLengthItemBacktrackCoverages) // allocation guarded by the previous check
		for i := range item.BacktrackCoverages {
			offset := int(binary.BigEndian.Uint16(src[4+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.BacktrackCoverages[i], _, err = ParseCoverage(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs3: %s", err)
			}

		}
		n += arrayLengthItemBacktrackCoverages * 2
	}
	if L := len(src); L < n+2 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: n + 2, got %d", L)
	}
	arrayLengthItemInputCoverages := int(binary.BigEndian.Uint16(src[n:]))
	n += 2

	{

		if L := len(src); L < n+arrayLengthItemInputCoverages*2 {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", n+arrayLengthItemInputCoverages*2, L)
		}

		item.InputCoverages = make([]Coverage, arrayLengthItemInputCoverages) // allocation guarded by the previous check
		for i := range item.InputCoverages {
			offset := int(binary.BigEndian.Uint16(src[n+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.InputCoverages[i], _, err = ParseCoverage(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs3: %s", err)
			}

		}
		n += arrayLengthItemInputCoverages * 2
	}
	if L := len(src); L < n+2 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: n + 2, got %d", L)
	}
	arrayLengthItemLookaheadCoverages := int(binary.BigEndian.Uint16(src[n:]))
	n += 2

	{

		if L := len(src); L < n+arrayLengthItemLookaheadCoverages*2 {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", n+arrayLengthItemLookaheadCoverages*2, L)
		}

		item.LookaheadCoverages = make([]Coverage, arrayLengthItemLookaheadCoverages) // allocation guarded by the previous check
		for i := range item.LookaheadCoverages {
			offset := int(binary.BigEndian.Uint16(src[n+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.LookaheadCoverages[i], _, err = ParseCoverage(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ChainedContextualSubs3: %s", err)
			}

		}
		n += arrayLengthItemLookaheadCoverages * 2
	}
	if L := len(src); L < n+2 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: n + 2, got %d", L)
	}
	arrayLengthItemSeqLookupRecords := int(binary.BigEndian.Uint16(src[n:]))
	n += 2

	{

		if L := len(src); L < n+arrayLengthItemSeqLookupRecords*4 {
			return item, 0, fmt.Errorf("reading ChainedContextualSubs3: "+"EOF: expected length: %d, got %d", n+arrayLengthItemSeqLookupRecords*4, L)
		}

		item.SeqLookupRecords = make([]SequenceLookupRecord, arrayLengthItemSeqLookupRecords) // allocation guarded by the previous check
		for i := range item.SeqLookupRecords {
			item.SeqLookupRecords[i].mustParse(src[n+i*4:])
		}
		n += arrayLengthItemSeqLookupRecords * 4
	}
	return item, n, nil
}

func ParseChainedContextualSubsITF(src []byte) (ChainedContextualSubsITF, int, error) {
	var item ChainedContextualSubsITF

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading ChainedContextualSubsITF: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 1:
		item, read, err = ParseChainedContextualSubs1(src[0:])
	case 2:
		item, read, err = ParseChainedContextualSubs2(src[0:])
	case 3:
		item, read, err = ParseChainedContextualSubs3(src[0:])
	default:
		err = fmt.Errorf("unsupported ChainedContextualSubsITF format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading ChainedContextualSubsITF: %s", err)
	}

	return item, read, nil
}

func ParseContextualSubs(src []byte) (ContextualSubs, int, error) {
	var item ContextualSubs
	n := 0
	{
		var (
			err  error
			read int
		)
		item.Data, read, err = ParseContextualSubsITF(src[0:])
		if err != nil {
			return item, 0, fmt.Errorf("reading ContextualSubs: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func ParseContextualSubs1(src []byte) (ContextualSubs1, int, error) {
	var item ContextualSubs1
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading ContextualSubs1: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemSeqRuleSet := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemSeqRuleSet*2 {
			return item, 0, fmt.Errorf("reading ContextualSubs1: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemSeqRuleSet*2, L)
		}

		item.SeqRuleSet = make([]SequenceRuleSet, arrayLengthItemSeqRuleSet) // allocation guarded by the previous check
		for i := range item.SeqRuleSet {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ContextualSubs1: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.SeqRuleSet[i], _, err = ParseSequenceRuleSet(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ContextualSubs1: %s", err)
			}

		}
		n += arrayLengthItemSeqRuleSet * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading ContextualSubs1: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ContextualSubs1: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	return item, n, nil
}

func ParseContextualSubs2(src []byte) (ContextualSubs2, int, error) {
	var item ContextualSubs2
	n := 0
	if L := len(src); L < 8 {
		return item, 0, fmt.Errorf("reading ContextualSubs2: "+"EOF: expected length: 8, got %d", L)
	}
	_ = src[7] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	offsetItemClassDef := int(binary.BigEndian.Uint16(src[4:]))
	arrayLengthItemClassSeqRuleSet := int(binary.BigEndian.Uint16(src[6:]))
	n += 8

	{

		if L := len(src); L < 8+arrayLengthItemClassSeqRuleSet*2 {
			return item, 0, fmt.Errorf("reading ContextualSubs2: "+"EOF: expected length: %d, got %d", 8+arrayLengthItemClassSeqRuleSet*2, L)
		}

		item.ClassSeqRuleSet = make([]ClassSequenceRuleSet, arrayLengthItemClassSeqRuleSet) // allocation guarded by the previous check
		for i := range item.ClassSeqRuleSet {
			offset := int(binary.BigEndian.Uint16(src[8+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ContextualSubs2: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.ClassSeqRuleSet[i], _, err = ParseClassSequenceRuleSet(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ContextualSubs2: %s", err)
			}

		}
		n += arrayLengthItemClassSeqRuleSet * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading ContextualSubs2: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ContextualSubs2: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	{

		if offsetItemClassDef != 0 { // ignore null offset
			if L := len(src); L < offsetItemClassDef {
				return item, 0, fmt.Errorf("reading ContextualSubs2: "+"EOF: expected length: %d, got %d", offsetItemClassDef, L)
			}

			var (
				err  error
				read int
			)
			item.ClassDef, read, err = ParseClassDef(src[offsetItemClassDef:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ContextualSubs2: %s", err)
			}
			offsetItemClassDef += read
		}
	}
	return item, n, nil
}

func ParseContextualSubs3(src []byte) (ContextualSubs3, int, error) {
	var item ContextualSubs3
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading ContextualSubs3: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	item.glyphCount = binary.BigEndian.Uint16(src[2:])
	item.seqLookupCount = binary.BigEndian.Uint16(src[4:])
	n += 6

	{
		arrayLength := int(item.glyphCount)

		if L := len(src); L < 6+arrayLength*2 {
			return item, 0, fmt.Errorf("reading ContextualSubs3: "+"EOF: expected length: %d, got %d", 6+arrayLength*2, L)
		}

		item.CoverageOffsets = make([]Coverage, arrayLength) // allocation guarded by the previous check
		for i := range item.CoverageOffsets {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ContextualSubs3: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.CoverageOffsets[i], _, err = ParseCoverage(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ContextualSubs3: %s", err)
			}

		}
		n += arrayLength * 2
	}
	{
		arrayLength := int(item.seqLookupCount)

		if L := len(src); L < n+arrayLength*4 {
			return item, 0, fmt.Errorf("reading ContextualSubs3: "+"EOF: expected length: %d, got %d", n+arrayLength*4, L)
		}

		item.SeqLookupRecords = make([]SequenceLookupRecord, arrayLength) // allocation guarded by the previous check
		for i := range item.SeqLookupRecords {
			item.SeqLookupRecords[i].mustParse(src[n+i*4:])
		}
		n += arrayLength * 4
	}
	return item, n, nil
}

func ParseContextualSubsITF(src []byte) (ContextualSubsITF, int, error) {
	var item ContextualSubsITF

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading ContextualSubsITF: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 1:
		item, read, err = ParseContextualSubs1(src[0:])
	case 2:
		item, read, err = ParseContextualSubs2(src[0:])
	case 3:
		item, read, err = ParseContextualSubs3(src[0:])
	default:
		err = fmt.Errorf("unsupported ContextualSubsITF format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading ContextualSubsITF: %s", err)
	}

	return item, read, nil
}

func ParseExtensionSubs(src []byte) (ExtensionSubs, int, error) {
	var item ExtensionSubs
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading ExtensionSubs: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.substFormat = binary.BigEndian.Uint16(src[0:])
	item.ExtensionLookupType = binary.BigEndian.Uint16(src[2:])
	item.ExtensionOffset = Offset32(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		item.RawData = src[0:]
		n = len(src)
	}
	return item, n, nil
}

func ParseLigature(src []byte) (Ligature, int, error) {
	var item Ligature
	n := 0
	if L := len(src); L < 4 {
		return item, 0, fmt.Errorf("reading Ligature: "+"EOF: expected length: 4, got %d", L)
	}
	_ = src[3] // early bound checking
	item.LigatureGlyph = GlyphID(binary.BigEndian.Uint16(src[0:]))
	item.componentCount = binary.BigEndian.Uint16(src[2:])
	n += 4

	{
		arrayLength := int(item.componentCount - 1)

		if L := len(src); L < 4+arrayLength*2 {
			return item, 0, fmt.Errorf("reading Ligature: "+"EOF: expected length: %d, got %d", 4+arrayLength*2, L)
		}

		item.ComponentGlyphIDs = make([]GlyphID, arrayLength) // allocation guarded by the previous check
		for i := range item.ComponentGlyphIDs {
			item.ComponentGlyphIDs[i] = GlyphID(binary.BigEndian.Uint16(src[4+i*2:]))
		}
		n += arrayLength * 2
	}
	return item, n, nil
}

func ParseLigatureSet(src []byte) (LigatureSet, int, error) {
	var item LigatureSet
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading LigatureSet: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthItemLigatures := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthItemLigatures*2 {
			return item, 0, fmt.Errorf("reading LigatureSet: "+"EOF: expected length: %d, got %d", 2+arrayLengthItemLigatures*2, L)
		}

		item.Ligatures = make([]Ligature, arrayLengthItemLigatures) // allocation guarded by the previous check
		for i := range item.Ligatures {
			offset := int(binary.BigEndian.Uint16(src[2+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading LigatureSet: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.Ligatures[i], _, err = ParseLigature(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading LigatureSet: %s", err)
			}

		}
		n += arrayLengthItemLigatures * 2
	}
	return item, n, nil
}

func ParseLigatureSubs(src []byte) (LigatureSubs, int, error) {
	var item LigatureSubs
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading LigatureSubs: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.substFormat = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemLigatureSets := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemLigatureSets*2 {
			return item, 0, fmt.Errorf("reading LigatureSubs: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemLigatureSets*2, L)
		}

		item.LigatureSets = make([]LigatureSet, arrayLengthItemLigatureSets) // allocation guarded by the previous check
		for i := range item.LigatureSets {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading LigatureSubs: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.LigatureSets[i], _, err = ParseLigatureSet(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading LigatureSubs: %s", err)
			}

		}
		n += arrayLengthItemLigatureSets * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading LigatureSubs: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading LigatureSubs: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	return item, n, nil
}

func ParseMultipleSubs(src []byte) (MultipleSubs, int, error) {
	var item MultipleSubs
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading MultipleSubs: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.substFormat = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverageOffset := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemSequences := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemSequences*2 {
			return item, 0, fmt.Errorf("reading MultipleSubs: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemSequences*2, L)
		}

		item.Sequences = make([]Sequence, arrayLengthItemSequences) // allocation guarded by the previous check
		for i := range item.Sequences {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading MultipleSubs: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.Sequences[i], _, err = ParseSequence(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading MultipleSubs: %s", err)
			}

		}
		n += arrayLengthItemSequences * 2
	}
	{

		if offsetItemCoverageOffset != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverageOffset {
				return item, 0, fmt.Errorf("reading MultipleSubs: "+"EOF: expected length: %d, got %d", offsetItemCoverageOffset, L)
			}

			var (
				err  error
				read int
			)
			item.CoverageOffset, read, err = ParseCoverage(src[offsetItemCoverageOffset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading MultipleSubs: %s", err)
			}
			offsetItemCoverageOffset += read
		}
	}
	return item, n, nil
}

func ParseReverseChainSingleSubs(src []byte) (ReverseChainSingleSubs, int, error) {
	var item ReverseChainSingleSubs
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.substFormat = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemBacktrackCoverageOffsets := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemBacktrackCoverageOffsets*2 {
			return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemBacktrackCoverageOffsets*2, L)
		}

		item.BacktrackCoverageOffsets = make([]Coverage, arrayLengthItemBacktrackCoverageOffsets) // allocation guarded by the previous check
		for i := range item.BacktrackCoverageOffsets {
			offset := int(binary.BigEndian.Uint16(src[6+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.BacktrackCoverageOffsets[i], _, err = ParseCoverage(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: %s", err)
			}

		}
		n += arrayLengthItemBacktrackCoverageOffsets * 2
	}
	if L := len(src); L < n+2 {
		return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: n + 2, got %d", L)
	}
	arrayLengthItemLookaheadCoverageOffsets := int(binary.BigEndian.Uint16(src[n:]))
	n += 2

	{

		if L := len(src); L < n+arrayLengthItemLookaheadCoverageOffsets*2 {
			return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: %d, got %d", n+arrayLengthItemLookaheadCoverageOffsets*2, L)
		}

		item.LookaheadCoverageOffsets = make([]Coverage, arrayLengthItemLookaheadCoverageOffsets) // allocation guarded by the previous check
		for i := range item.LookaheadCoverageOffsets {
			offset := int(binary.BigEndian.Uint16(src[n+i*2:]))
			// ignore null offsets
			if offset == 0 {
				continue
			}

			if L := len(src); L < offset {
				return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: %d, got %d", offset, L)
			}

			var err error
			item.LookaheadCoverageOffsets[i], _, err = ParseCoverage(src[offset:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: %s", err)
			}

		}
		n += arrayLengthItemLookaheadCoverageOffsets * 2
	}
	if L := len(src); L < n+2 {
		return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: n + 2, got %d", L)
	}
	arrayLengthItemSubstituteGlyphIDs := int(binary.BigEndian.Uint16(src[n:]))
	n += 2

	{

		if L := len(src); L < n+arrayLengthItemSubstituteGlyphIDs*2 {
			return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: %d, got %d", n+arrayLengthItemSubstituteGlyphIDs*2, L)
		}

		item.SubstituteGlyphIDs = make([]GlyphID, arrayLengthItemSubstituteGlyphIDs) // allocation guarded by the previous check
		for i := range item.SubstituteGlyphIDs {
			item.SubstituteGlyphIDs[i] = GlyphID(binary.BigEndian.Uint16(src[n+i*2:]))
		}
		n += arrayLengthItemSubstituteGlyphIDs * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading ReverseChainSingleSubs: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	return item, n, nil
}

func ParseSequence(src []byte) (Sequence, int, error) {
	var item Sequence
	n := 0
	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading Sequence: "+"EOF: expected length: 2, got %d", L)
	}
	arrayLengthItemSubstituteGlyphIDs := int(binary.BigEndian.Uint16(src[0:]))
	n += 2

	{

		if L := len(src); L < 2+arrayLengthItemSubstituteGlyphIDs*2 {
			return item, 0, fmt.Errorf("reading Sequence: "+"EOF: expected length: %d, got %d", 2+arrayLengthItemSubstituteGlyphIDs*2, L)
		}

		item.SubstituteGlyphIDs = make([]GlyphID, arrayLengthItemSubstituteGlyphIDs) // allocation guarded by the previous check
		for i := range item.SubstituteGlyphIDs {
			item.SubstituteGlyphIDs[i] = GlyphID(binary.BigEndian.Uint16(src[2+i*2:]))
		}
		n += arrayLengthItemSubstituteGlyphIDs * 2
	}
	return item, n, nil
}

func ParseSingleSubs(src []byte) (SingleSubs, int, error) {
	var item SingleSubs
	n := 0
	{
		var (
			err  error
			read int
		)
		item.Data, read, err = ParseSingleSubstData(src[0:])
		if err != nil {
			return item, 0, fmt.Errorf("reading SingleSubs: %s", err)
		}
		n += read
	}
	return item, n, nil
}

func ParseSingleSubstData(src []byte) (SingleSubstData, int, error) {
	var item SingleSubstData

	if L := len(src); L < 2 {
		return item, 0, fmt.Errorf("reading SingleSubstData: "+"EOF: expected length: 2, got %d", L)
	}
	format := uint16(binary.BigEndian.Uint16(src[0:]))
	var (
		read int
		err  error
	)
	switch format {
	case 1:
		item, read, err = ParseSingleSubstData1(src[0:])
	case 2:
		item, read, err = ParseSingleSubstData2(src[0:])
	default:
		err = fmt.Errorf("unsupported SingleSubstData format %d", format)
	}
	if err != nil {
		return item, 0, fmt.Errorf("reading SingleSubstData: %s", err)
	}

	return item, read, nil
}

func ParseSingleSubstData1(src []byte) (SingleSubstData1, int, error) {
	var item SingleSubstData1
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading SingleSubstData1: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	item.DeltaGlyphID = int16(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading SingleSubstData1: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading SingleSubstData1: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	return item, n, nil
}

func ParseSingleSubstData2(src []byte) (SingleSubstData2, int, error) {
	var item SingleSubstData2
	n := 0
	if L := len(src); L < 6 {
		return item, 0, fmt.Errorf("reading SingleSubstData2: "+"EOF: expected length: 6, got %d", L)
	}
	_ = src[5] // early bound checking
	item.format = binary.BigEndian.Uint16(src[0:])
	offsetItemCoverage := int(binary.BigEndian.Uint16(src[2:]))
	arrayLengthItemSubstituteGlyphIDs := int(binary.BigEndian.Uint16(src[4:]))
	n += 6

	{

		if L := len(src); L < 6+arrayLengthItemSubstituteGlyphIDs*2 {
			return item, 0, fmt.Errorf("reading SingleSubstData2: "+"EOF: expected length: %d, got %d", 6+arrayLengthItemSubstituteGlyphIDs*2, L)
		}

		item.SubstituteGlyphIDs = make([]GlyphID, arrayLengthItemSubstituteGlyphIDs) // allocation guarded by the previous check
		for i := range item.SubstituteGlyphIDs {
			item.SubstituteGlyphIDs[i] = GlyphID(binary.BigEndian.Uint16(src[6+i*2:]))
		}
		n += arrayLengthItemSubstituteGlyphIDs * 2
	}
	{

		if offsetItemCoverage != 0 { // ignore null offset
			if L := len(src); L < offsetItemCoverage {
				return item, 0, fmt.Errorf("reading SingleSubstData2: "+"EOF: expected length: %d, got %d", offsetItemCoverage, L)
			}

			var (
				err  error
				read int
			)
			item.Coverage, read, err = ParseCoverage(src[offsetItemCoverage:])
			if err != nil {
				return item, 0, fmt.Errorf("reading SingleSubstData2: %s", err)
			}
			offsetItemCoverage += read
		}
	}
	return item, n, nil
}
