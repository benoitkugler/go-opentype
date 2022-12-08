package harfbuzz

import (
	"fmt"
	"math/bits"

	"github.com/benoitkugler/go-opentype/api"
	"github.com/benoitkugler/go-opentype/api/font"
	"github.com/benoitkugler/go-opentype/tables"
)

const maxContextLength = 64

var _ layoutLookup = lookupGSUB{}

// implements layoutLookup
type lookupGSUB font.GSUBLookup

func (l lookupGSUB) collectCoverage(dst *setDigest) {
	for _, table := range l.Subtables {
		dst.collectCoverage(table.Coverage)
	}
}

func (l lookupGSUB) dispatchSubtables(ctx *getSubtablesContext) {
	for _, table := range l.Subtables {
		*ctx = append(*ctx, newGSUBApplicable(table))
	}
}

func (l lookupGSUB) dispatchApply(ctx *otApplyContext) bool {
	for _, table := range l.Subtables {
		if gsubSubtable(table).apply(ctx) {
			return true
		}
	}
	return false
}

func (l lookupGSUB) wouldApply(ctx *wouldApplyContext, accel *otLayoutLookupAccelerator) bool {
	if len(ctx.glyphs) == 0 {
		return false
	}
	if !accel.digest.mayHave(ctx.glyphs[0]) {
		return false
	}
	// dispatch on subtables
	for _, table := range l.Subtables {
		if gsubSubtable(table).wouldApply(ctx) {
			return true
		}
	}
	return false
}

func (l lookupGSUB) isReverse() bool { return l.Type == tables.GSUBReverse }

func applyRecurseGSUB(c *otApplyContext, lookupIndex uint16) bool {
	gsub := c.font.face.GSUB
	l := lookupGSUB(gsub.Lookups[lookupIndex])
	return c.applyRecurseLookup(lookupIndex, l)
}

// implements `hb_apply_func_t`
type gsubSubtable struct {
	tables.GSUBLookup
}

// matchesLigature tests if the ligature should be applied on `glyphsFromSecond`,
// which starts from the second glyph.
func matchesLigature(l tables.Ligature, glyphsFromSecond []api.GID) bool {
	if len(glyphsFromSecond) != len(l.ComponentGlyphIDs) {
		return false
	}
	for i, g := range glyphsFromSecond {
		if g != api.GID(l.ComponentGlyphIDs[i]) {
			return false
		}
	}
	return true
}

// return `true` is we should apply this lookup to the glyphs in `c`,
// which are assumed to be non empty
func (table gsubSubtable) wouldApply(c *wouldApplyContext) bool {
	index, ok := table.Coverage().Index(tables.GlyphID(c.glyphs[0]))
	switch data := table.GSUBLookup.(type) {
	case tables.SingleSubs, tables.MultipleSubs, tables.AlternateSubs, tables.ReverseChainSingleSubs:
		return len(c.glyphs) == 1 && ok

	case tables.LigatureSubs:
		if !ok {
			return false
		}
		ligatureSet := data.LigatureSets[index].Ligatures
		glyphsFromSecond := c.glyphs[1:]
		for _, ligature := range ligatureSet {
			if matchesLigature(ligature, glyphsFromSecond) {
				return true
			}
		}
		return false

	case tables.GSUBContext1:
		return c.wouldApplyLookupContext1(tables.LookupContext1(data), index)
	case tables.GSUBContext2:
		return c.wouldApplyLookupContext2(tables.LookupContext2(data), index, c.glyphs[0])
	case tables.GSUBContext3:
		return c.wouldApplyLookupContext3(tables.LookupContext3(data), index)
	case tables.GSUBChainedContext1:
		return c.wouldApplyLookupChainedContext1(tables.LookupChainedContext1(data), index)
	case tables.GSUBChainedContext2:
		return c.wouldApplyLookupChainedContext2(tables.LookupChainedContext2(data), index, c.glyphs[0])
	case tables.GSUBChainedContext3:
		return c.wouldApplyLookupChainedContext3(tables.LookupChainedContext3(data), index)
	}
	return false
}

// return `true` is the subsitution found a match and was applied
func (table gsubSubtable) apply(c *otApplyContext) bool {
	glyph := c.buffer.cur(0)
	glyphID := glyph.Glyph
	index, ok := table.Coverage.Index(glyphID)
	if !ok {
		return false
	}

	if debugMode >= 2 {
		fmt.Printf("\tAPPLY - type %T at index %d\n", table.Data, c.buffer.idx)
	}

	switch data := table.GSUBLookup.(type) {
	case tables.GSUBSingle1:
		/* According to the Adobe Annotated OpenType Suite, result is always
		* limited to 16bit. */
		glyphID = api.GID(uint16(int(glyphID) + int(data)))
		c.replaceGlyph(glyphID)
	case tables.GSUBSingle2:
		if index >= len(data) { // index is not sanitized in tables.Parse
			return false
		}
		c.replaceGlyph(data[index])

	case tables.MultipleSubs:
		c.applySubsSequence(data.Sequences[index].SubstituteGlyphIDs)

	case tables.GSUBAlternate1:
		alternates := data[index]
		return c.applySubsAlternate(alternates)

	case tables.GSUBLigature1:
		ligatureSet := data[index]
		return c.applySubsLigature(ligatureSet)

	case tables.GSUBContext1:
		return c.applyLookupContext1(tables.LookupContext1(data), index)
	case tables.GSUBContext2:
		return c.applyLookupContext2(tables.LookupContext2(data), index, glyphID)
	case tables.GSUBContext3:
		return c.applyLookupContext3(tables.LookupContext3(data), index)
	case tables.GSUBChainedContext1:
		return c.applyLookupChainedContext1(tables.LookupChainedContext1(data), index)
	case tables.GSUBChainedContext2:
		return c.applyLookupChainedContext2(tables.LookupChainedContext2(data), index, glyphID)
	case tables.GSUBChainedContext3:
		return c.applyLookupChainedContext3(tables.LookupChainedContext3(data), index)

	case tables.GSUBReverseChainedContext1:
		if c.nestingLevelLeft != maxNestingLevel {
			return false // no chaining to this type
		}
		lB, lL := len(data.Backtrack), len(data.Lookahead)
		hasMatch, startIndex := c.matchBacktrack(get1N(&c.indices, 0, lB), matchCoverage(data.Backtrack))
		if !hasMatch {
			return false
		}

		hasMatch, endIndex := c.matchLookahead(get1N(&c.indices, 0, lL), matchCoverage(data.Lookahead), 1)
		if !hasMatch {
			return false
		}

		c.buffer.unsafeToBreakFromOutbuffer(startIndex, endIndex)
		c.setGlyphProps(data.Substitutes[index])
		c.buffer.cur(0).Glyph = data.Substitutes[index]
		/* Note: We DON'T decrease buffer.idx.  The main loop does it
		 * for us.  This is useful for preventing surprises if someone
		 * calls us through a Context lookup. */

	}

	return true
}

func (c *otApplyContext) applySubsSequence(seq []tables.GlyphID) {
	/* Special-case to make it in-place and not consider this
	 * as a "multiplied" substitution. */
	switch len(seq) {
	case 1:
		c.replaceGlyph(seq[0])
	case 0:
		/* Spec disallows this, but Uniscribe allows it.
		 * https://github.com/harfbuzz/harfbuzz/issues/253 */
		c.buffer.deleteGlyph()
	default:
		var klass uint16
		if c.buffer.cur(0).isLigature() {
			klass = tables.BaseGlyph
		}
		ligID := c.buffer.cur(0).getLigID()
		for i, g := range seq {
			/* If is attached to a ligature, don't disturb that.
			 * https://github.com/harfbuzz/harfbuzz/issues/3069 */
			if ligID == 0 {
				c.buffer.cur(0).setLigPropsForMark(0, uint8(i))
			}
			c.setGlyphPropsExt(g, klass, false, true)
			c.buffer.outputGlyphIndex(g)
		}
		c.buffer.skipGlyph()
	}
}

func (c *otApplyContext) applySubsAlternate(alternates []api.GID) bool {
	count := uint32(len(alternates))
	if count == 0 {
		return false
	}

	glyphMask := c.buffer.cur(0).Mask
	lookupMask := c.lookupMask

	/* Note: This breaks badly if two features enabled this lookup together. */

	shift := bits.TrailingZeros32(lookupMask)
	altIndex := (lookupMask & glyphMask) >> shift

	/* If altIndex is MAX_VALUE, randomize feature if it is the rand feature. */
	if altIndex == otMapMaxValue && c.random {
		// Maybe we can do better than unsafe-to-break all; but since we are
		// changing random state, it would be hard to track that.  Good 'nough.
		c.buffer.unsafeToBreak(0, len(c.buffer.Info))
		altIndex = c.randomNumber()%count + 1
	}

	if altIndex > count || altIndex == 0 {
		return false
	}

	c.replaceGlyph(alternates[altIndex-1])
	return true
}

func (c *otApplyContext) applySubsLigature(ligatureSet []tables.LigatureGlyph) bool {
	for _, lig := range ligatureSet {
		count := len(lig.Components) + 1

		/* Special-case to make it in-place and not consider this
		 * as a "ligated" substitution. */
		if count == 1 {
			c.replaceGlyph(lig.Glyph)
			return true
		}

		var matchPositions [maxContextLength]int

		ok, matchLength, totalComponentCount := c.matchInput(lig.Components, matchGlyph, &matchPositions)
		if !ok {
			continue
		}
		c.ligateInput(count, matchPositions, matchLength, lig.Glyph, totalComponentCount)

		return true
	}
	return false
}
