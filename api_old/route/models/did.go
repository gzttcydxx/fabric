package models

type DID struct {
	Scheme     string
	Method     string
	ChainID    string
	SpecificID string
	Fragment   string
}
