package font

import "github.com/benoitkugler/go-opentype/tables"

type GSUB struct {
	Lookups []GSUBLookup
}

type LookupOptions struct {
	// Lookup qualifiers.
	Flag uint16
	// Index (base 0) into GDEF mark glyph sets structure,
	// meaningfull only if UseMarkFilteringSet is set.
	MarkFilteringSet uint16
}

type GSUBLookup struct {
	LookupOptions
	Subtables []tables.GSUBLookup
}

func newGSUB(table tables.Layout) (GSUB, error) {
	out := GSUB{
		Lookups: make([]GSUBLookup, len(table.LookupList.Lookups)),
	}
	for i, lk := range table.LookupList.Lookups {
		subtables, err := lk.AsGSUBLookups()
		if err != nil {
			return GSUB{}, err
		}
		for j, subtable := range subtables {
			// start by resolving extension
			if ext, isExt := subtable.(tables.ExtensionSubs); isExt {
				subtables[j], err = ext.Resolve()
				if err != nil {
					return GSUB{}, err
				}
			}

			// sanitize each lookup
			switch subtable := subtable.(type) {
			case tables.MultipleSubs:
				err = subtable.Sanitize()
			case tables.LigatureSubs:
				err = subtable.Sanitize()
			case tables.ContextualSubs:
				err = subtable.Sanitize(uint16(len(out.Lookups)))
			case tables.ReverseChainSingleSubs:
				err = subtable.Sanitize()
			}
			if err != nil {
				return GSUB{}, err
			}
		}
		out.Lookups[i] = GSUBLookup{
			LookupOptions: LookupOptions{
				Flag:             lk.LookupFlag,
				MarkFilteringSet: lk.MarkFilteringSet,
			},
			Subtables: subtables,
		}
	}
	return out, nil
}

type GPOS struct {
	Lookups []GPOSLookup
}

type GPOSLookup struct {
	LookupOptions
	Subtables []tables.GPOSLookup
}

func newGPOS(table tables.Layout) (GPOS, error) {
	out := GPOS{
		Lookups: make([]GPOSLookup, len(table.LookupList.Lookups)),
	}
	for i, lk := range table.LookupList.Lookups {
		subtables, err := lk.AsGPOSLookups()
		if err != nil {
			return GPOS{}, err
		}
		for j, subtable := range subtables {
			// start by resolving extension
			if ext, isExt := subtable.(tables.ExtensionPos); isExt {
				subtables[j], err = ext.Resolve()
				if err != nil {
					return GPOS{}, err
				}
			}

			// sanitize each lookup
			switch subtable := subtable.(type) {
			case tables.SinglePos:
				err = subtable.Sanitize()
			case tables.PairPos:
				err = subtable.Sanitize()
			case tables.MarkBasePos:
				err = subtable.Sanitize()
			case tables.MarkLigPos:
				err = subtable.Sanitize()
			case tables.ContextualPos:
				err = subtable.Sanitize(uint16(len(out.Lookups)))
			}
			if err != nil {
				return GPOS{}, err
			}
		}
		out.Lookups[i] = GPOSLookup{
			LookupOptions: LookupOptions{
				Flag:             lk.LookupFlag,
				MarkFilteringSet: lk.MarkFilteringSet,
			},
			Subtables: subtables,
		}
	}
	return out, nil
}
