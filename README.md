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
 - loading, transforming, saving a font file : 
    TODO:

The lower level font opening is provided by the [opentype](opentype/) package. Given an opened font file, its tables may be accessed under three forms (from lower to higher level): a raw slice of bytes, the corresponding binary content represented as a Go struct (defined in [tables](/tables/)) or an high level type optimized for text layout.
TODO:

Note that the two first forms contain exactly the same information whereas the last only retains layout oriented data.
The main interest of the intermediate form is that it provides methods to be serialized back to a slice of bytes, to simplify font transformations.
