# go-opentype
A go library to read, write and subset Truetype/Opentype fonts

## Scope of the package 

This library is a rewrite of [textlayout/fonts/truetype](https://pkg.go.dev/github.com/benoitkugler/textlayout@v0.2.0/fonts/truetype) using an auto-generated parser and writer. 

It serves two purposes :
  - providing font information needed for text layout
  - transforming fonts (usually to be included in PDF): instantiating a variable font; subsetting.
  
  
