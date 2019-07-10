package analyze

type KeySetAnalysis struct {
	Hash    string
	Signers []Signer
	Valid   bool
	Errors  []string
}

type Signer struct {
	Type  string
	Count string
	PQR   bool
	Valid bool
}