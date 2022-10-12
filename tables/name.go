package tables

//go:generate ../../binarygen/cmd/generator name.go

// Naming table
// See https://learn.microsoft.com/en-us/typography/opentype/spec/name
type Name struct {
	version     uint16
	count       uint16
	stringData  []byte       `offset-size:"16" len:"_toEnd"`
	nameRecords []nameRecord `len:"count"`
}

type nameRecord struct {
	platformID   PlatformID
	encodingID   EncodingID
	languageID   LanguageID
	nameID       NameID
	length       uint16
	stringOffset uint16
}

// TODO : we are not yet handling language tags
