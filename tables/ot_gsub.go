package tables

import "fmt"

// GSUB is the Glyph Substitution (GSUB) table.
// It provides data for substition of glyphs for appropriate rendering of scripts,
// such as cursively-connecting forms in Arabic script,
// or for advanced typographic effects, such as ligatures.
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/gsub
type GSUB Layout

type GSUBLookup interface {
	isGSUBLookup()
}

func (SingleSubstitution) isGSUBLookup()             {}
func (MultipleSubstitution) isGSUBLookup()           {}
func (AlternateSubstitution) isGSUBLookup()          {}
func (LigatureSubstitution) isGSUBLookup()           {}
func (ContextualSubstitution) isGSUBLookup()         {}
func (ChainedContextualSubstitution) isGSUBLookup()  {}
func (ExtensionSubstitution) isGSUBLookup()          {}
func (ReverseChainSingleSubstitution) isGSUBLookup() {}

// AsGSUBLookups returns the GSUB lookup subtables
func (lk Lookup) AsGSUBLookups() ([]GSUBLookup, error) {
	var err error
	out := make([]GSUBLookup, len(lk.subtableOffsets))
	for i, offset := range lk.subtableOffsets {
		if L := len(lk.rawData); L < int(offset) {
			return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		switch lk.lookupType {
		case 1: // Single (format 1.1 1.2)	Replace one glyph with one glyph
			out[i], _, err = ParseSingleSubstitution(lk.rawData[offset:])
		case 2: // Multiple (format 2.1)	Replace one glyph with more than one glyph
			out[i], _, err = ParseMultipleSubstitution(lk.rawData[offset:])
		case 3: // Alternate (format 3.1)	Replace one glyph with one of many glyphs
			out[i], _, err = ParseAlternateSubstitution(lk.rawData[offset:])
		case 4: // Ligature (format 4.1)	Replace multiple glyphs with one glyph
			out[i], _, err = ParseLigatureSubstitution(lk.rawData[offset:])
		case 5: // Context (format 5.1 5.2 5.3)	Replace one or more glyphs in context
			out[i], _, err = ParseContextualSubstitution(lk.rawData[offset:])
		case 6: // Chaining Context (format 6.1 6.2 6.3)	Replace one or more glyphs in chained context
			out[i], _, err = ParseChainedContextualSubstitution(lk.rawData[offset:])
		case 7: // Extension Substitution (format 7.1) Extension mechanism for other substitutions
			out[i], _, err = ParseExtensionSubstitution(lk.rawData[offset:])
		case 8: // Reverse chaining context single (format 8.1)
			out[i], _, err = ParseReverseChainSingleSubstitution(lk.rawData[offset:])
		default:
			err = fmt.Errorf("invalid GSUB Loopkup type %d", lk.lookupType)
		}
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
