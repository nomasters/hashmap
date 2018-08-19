package hashmap

import (
	"errors"
)

const (
	examplePrivateKeyPath = "example_files/priv.key"
	exampleValidPayload   = `
		{
			"data": "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			"sig": "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			"pubkey": "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo="
		}
	`
)

var (
	examplePayload = &Payload{
		Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
		Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
		PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
	}
)

type errorReader struct{}

func (er errorReader) Read(b []byte) (int, error) {
	return 0, errors.New("arbitrary")
}
