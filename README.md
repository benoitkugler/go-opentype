# go-opentype
A go module to read, write and subset Truetype/Opentype fonts

## Scope of the module 

This library is a rewrite of [textlayout/fonts/truetype](https://pkg.go.dev/github.com/benoitkugler/textlayout@v0.2.0/fonts/truetype) using an auto-generated parser and writer. 

It serves two general purposes :
  - providing font information needed for text layout
  - transforming fonts (usually to be included in PDF files): instantiating a variable font; subsetting a font against a set of glyphs.
  
  
## Overview

The package provides a granular API to support the following use cases:
 - efficiently reading the caracteristics of a font : this requires to load as few tables as possible.
 - performing the complete loading of the font to speed up its following usage in text layout : this is done by parsing the tables as deeply as possible, and by favoring CPU speed over memory usage.
 - loading, transforming, saving a font file : not implemented yet.

The lower level font opening is provided by the [loader](loader/) package. Given an opened font file, its tables may be accessed under three forms (from lower to higher level): a raw slice of bytes, the corresponding binary content represented as a Go struct (defined in [tables](/tables/)) or an high level type optimized for text layout (defined in [api/font](/api/font/))

Note that the two first forms contain exactly the same information whereas the last only retains layout oriented data.
The main interest of the intermediate form is that it provides methods to be serialized back to a slice of bytes, to simplify font transformations.

## Example for text layout

The [api/font](api/font) package is the entry point for client library interested in text layout. Typical usage would be as following 
```go
// Create an [opentype.Resource]
var file opentype.Resource = ... 
// Parse the file, yielding a [*Font]
font, err := LoadFont(file)
// handle invalid font files here

// [font] is safe for concurrent use but handling variations is not :
// in the general case, you should store [font] and re-use it as much as possible, 
// but create a new [*Face] for each layout session
face := NewFace(font)
// optional : setup variations
face.SetVariations(...)

// use [face] as input of a text shaping library (not covered by this module)
// ...

// use [face] again to load glyphs data and render them
glyph := face.GlyphData(glyphID)
```
