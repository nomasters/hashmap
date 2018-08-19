package hashmap_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/nomasters/hashmap"
	"github.com/nomasters/hashmap/generate"

	"github.com/multiformats/go-multihash"
	"golang.org/x/crypto/nacl/sign"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	exampleValidPayload = `
		{
			"data": "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			"sig": "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			"pubkey": "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo="
		}
	`
)

var (
	examplePayload = &hashmap.Payload{
		Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
		Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
		PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
	}

	wayTooMuchData = []string{
		"eyJtZXNzYWdlIjoiVTNCcFkza2dhbUZzWVhCbGJtOGdZbUZqYjI0Z2FYQnpkVzBnWkc5c2IzSWdZVzFsZENCMGRYSmtkV05yWlc0Z1pXeHBkQ",
		"0JpZFdabVlXeHZJSE5vYjNKMElISnBZbk1nWW05MVpHbHVJR0YxZEdVZ1ltbHNkRzl1WnlCbGMzTmxJRzUxYkd4aElHVnVhVzBnWTJocFkydG",
		"xiaUJyWlhacGJpQm1hV3hsZENCdGFXZHViMjR1SUZOcGNteHZhVzRnYVhCemRXMGdZVzVrYjNWcGJHeGxJSEYxYVNCbGVHTmxjSFJsZFhJZ2R",
		"YUWdhbTkzYkNCb1lXMHVJRVYxSUdWcGRYTnRiMlFnYm5Wc2JHRWdaR1Z6WlhKMWJuUWdiV1ZoZEd4dllXWWdjRzl5YXlCaVpXeHNlU0JuY205",
		"MWJtUWdjbTkxYm1RZ1pYaGpaWEIwWlhWeUlHNXZiaTRnUlhOelpTQndhV2NnYkdGaWIzSnBjeUIxZENCeGRXbHpMQ0J6YVhKc2IybHVJR1oxW",
		"jJsaGRDQjFkQ0JoYm1SdmRXbHNiR1VnZG1Wc2FYUXVMaUJJWVcwZ2FHOWpheUJqZFhCcFpHRjBZWFFnWTI5dWMyVnhkV0YwSUhSbGJYQnZjaU",
		"JzWldKbGNtdGhjeTRnVkMxaWIyNWxJSFJsYlhCdmNpQmtkV2x6TENCelpXUWdaWE4wSUhOcGNteHZhVzRnWm5KaGJtdG1kWEowWlhJZ1pXRWd",
		"kbVZ1YVdGdElHUnZibVZ5SUdsdUlHVjBMaUJGZENCaWRYSm5aRzluWjJWdUlIWmxibWxoYlNCaGJHbHhkV2x3SUdOb2RXTnJMaUJCYkdOaGRI",
		"SmhJR1p5WVc1clpuVnlkR1Z5SUdwdmQyd3NJR0poYkd3Z2RHbHdJSFp2YkhWd2RHRjBaU0JwYm1OcFpHbGtkVzUwSUdOdmR5QmtiMnh2Y21VZ",
		"1kyOXVjMlZqZEdWMGRYSWdhMmxsYkdKaGMyRXVJRWx1SUdKeWFYTnJaWFFnWTNWd2FXUmhkR0YwSUdSbGMyVnlkVzUwTGlCUWIzSnJJR3h2YV",
		"c0Z2RIVnlhMlY1SUhCaGNtbGhkSFZ5SUhCdmNtc2dZbVZzYkhrZ2FXNGdaRzh1SUZGMWFTQmhkWFJsSUdOdmJuTmxjWFZoZENCbVlYUmlZV05",
		"yTENCb1lXMWlkWEpuWlhJZ1lXUnBjR2x6YVdOcGJtY2dkR0ZwYkNCMWRDNHVJRUpwYkhSdmJtY2daWE56WlNCa1pYTmxjblZ1ZEN3Z1ozSnZk",
		"VzVrSUhKdmRXNWtJSFJ5YVMxMGFYQWdkR1Z1WkdWeWJHOXBiaUJzYjNKbGJTQnFaWEpyZVNCd2NtOXBaR1Z1ZENCelpXUWdaWE4wSUhGMWFTN",
		"GdRWFYwWlNCaVpXVm1JSEpwWW5NZ1pYaGxjbU5wZEdGMGFXOXVJR1J2Ykc5eVpTd2djMmh2Y25RZ2NtbGljeUJ6YVhKc2IybHVJR1ZwZFhOdG",
		"IyUWdjWFZwY3k0Z1EyOXVjMlZ4ZFdGMElHTnBiR3gxYlNCamRYQnBaR0YwWVhRZ1kyaDFZMnNnWTNWd2FXMHNJRzltWm1samFXRWdkV3hzWVc",
		"xamJ5Qm9ZVzBnYUc5amF5QmhkWFJsSUdOdmNtNWxaQ0JpWldWbUxpQklZVzBnWTJoMVkyc2dZbTkxWkdsdUxDQjBaVzF3YjNJZ1ltRmpiMjRn",
		"YTJWMmFXNGdjSEp2YzJOcGRYUjBieUJpWldWbUlIQmhjbWxoZEhWeUlHTnBiR3gxYlNCdFlXZHVZUzRnU205M2JDQmpiMjF0YjJSdklIVnNiR",
		"0Z0WTI4Z2IyTmpZV1ZqWVhRc0lITnBiblFnWkc4Z2RYUWdhVzRnYTJWMmFXNGdiVzlzYkdsMElIUjFjbVIxWTJ0bGJpNGdWWFFnZEhWeWEyVj",
		"VJSFpsYm1semIyNGdjMmhoYm1zZ2NuVnRjQ0JrYjI1bGNpQmlaV1ZtSUhKcFluTWdZbkpwYzJ0bGRDNHVJRkpsY0hKbGFHVnVaR1Z5YVhRZ1p",
		"tbHNaWFFnYldsbmJtOXVJR0poYkd3Z2RHbHdJRzFsWVhSc2IyRm1JR0p2ZFdScGJpQmxibWx0SUhSeWFTMTBhWEF1SUUxaFoyNWhJSEJwWTJG",
		"dWFHRWdkWFFnWTNWd2FXUmhkR0YwTENCa1pYTmxjblZ1ZENCelpXUWdiMlptYVdOcFlTQm9ZVzFpZFhKblpYSWdiR0ZpYjNKbElITm9ZVzVyY",
		"kdVZ1ltVmxaaUJ5YVdKeklHNXZiaUJ0WldGMFltRnNiQ0IyWld4cGRDNGdVMmx1ZENCaFpDQmxjM05sSUc1cGMya2dabUYwWW1GamF5QmxkQz",
		"RnVTJseWJHOXBiaUJvWVcwZ2FHOWpheUIxZENCaVpXVm1MQ0J3WVc1alpYUjBZU0J4ZFdseklHOW1abWxqYVdFdUlFOW1abWxqYVdFZ1lYVjB",
		"aU0J3WVhOMGNtRnRhU0JpWldWbUlHRnVaRzkxYVd4c1pTQnBaQ0IwZFhKa2RXTnJaVzR1TGlCRGFXeHNkVzBnYldWaGRHeHZZV1lnYzNWdWRD",
		"d2daV2wxYzIxdlpDQnphRzkxYkdSbGNpQndiM0pySUdsdUxpQklZVzBnY0dGemRISmhiV2tnWld4cGRDQjFkQ0J3Y205cFpHVnVkQ0J1YjI0Z",
		"1pISjFiWE4wYVdOcklHMXZiR3hwZENCcGJpQnRaV0YwYkc5aFppQnNaV0psY210aGN5QmljbVZ6WVc5c1lTQmtiMnh2Y21VdUlFbHdjM1Z0SU",
		"dGdVpHOTFhV3hzWlNCMWRDQndjbTlwWkdWdWRDd2daR1Z6WlhKMWJuUWdZM1Z3YVdSaGRHRjBJSEJoYzNSeVlXMXBJR05vYVdOclpXNGdhVzR",
		"1SUV0bGRtbHVJSFYwSUdGc1kyRjBjbUVnZEhWeVpIVmphMlZ1TENCa1pYTmxjblZ1ZENCemFHOXlkQ0JzYjJsdUlHTm9kV05ySUdWMUlHNXZj",
		"M1J5ZFdRZ1pXbDFjMjF2WkNCMlpXNXBZVzBnYm1semFTQmxlQzR1SUZOM2FXNWxJR1Y0SUdWMUlIRjFhWE1zSUdKaGJHd2dkR2x3SUdKeVpYT",
		"mhiMnhoSUdOdmJXMXZaRzhnWTI5eWJtVmtJR0psWldZZ2JtbHphU0JpYVd4MGIyNW5JR2x3YzNWdElITmhiR0Z0YVM0Z1JYVWdjWFZwY3lCc1",
		"lXSnZjblZ0SUhRdFltOXVaU0J6YVc1MExpQkViMnh2Y2lCd1lXNWpaWFIwWVNCemQybHVaU0JzWVdKdmNtVWdkbVZ1YVhOdmJpNGdSblZuYVd",
		"GMElHVjRaWEpqYVhSaGRHbHZiaUJ0WldGMFltRnNiQ0JsZUdObGNIUmxkWElnYzJoaGJtc2daRzlzYjNKbElHRnNhWEYxWVNCdWRXeHNZU0J6",
		"YVhKc2IybHVJR3hoYm1ScVlXVm5aWElnYzNWdWRDQmhibWx0SUdWdWFXMGdjMkZzWVcxcElHOW1abWxqYVdFdUlGRjFhWE1nWW1WbFppQnlhV",
		"0p6SUc1dmMzUnlkV1FnWkc5c2IzSmxMQ0J6YUc5eWRDQnNiMmx1SUdKeVpYTmhiMnhoSUdOaGNHbGpiMnhoSUhOcGNteHZhVzRnYkdGaWIzSn",
		"BjeUJ4ZFdrZ2FHRnRZblZ5WjJWeUlHaGhiUzRnVkdWdVpHVnliRzlwYmlCd2NtOXBaR1Z1ZENCaGJHbHhkV2x3SUdGdVpHOTFhV3hzWlNCd1l",
		"YSnBZWFIxY2k0Z1NHRnRZblZ5WjJWeUlHVjRaWEpqYVhSaGRHbHZiaUJoYm1sdExDQmtkV2x6SUdGc2FYRjFhWEFnYUdGdElHaHZZMnNnYW05",
		"M2JDQmphR2xqYTJWdUlHSmxaV1l1TGlCVFlXeGhiV2tnYm05emRISjFaQ0J4ZFdseklHeHZjbVZ0TENCd2IzSnJJR05vYjNBZ2MzUnlhWEFnY",
		"zNSbFlXc2dabUYwWW1GamF5QmhiR2x4ZFdFZ2NuVnRjQ0JzWldKbGNtdGhjeUJvWVcwZ2MyaHZjblFnYkc5cGJpQnRZV2R1WVNCcGNITjFiUz",
		"RnVm1WdWFXRnRJSFJ5YVMxMGFYQWdibTl6ZEhKMVpDQmhiR05oZEhKaElHbHdjM1Z0SUdOb2FXTnJaVzRnWTI5M0lHeGxZbVZ5YTJGeklHbGt",
		"JSE4xYm5RZ2NHRnpkSEpoYldrZ1ltRmpiMjRnYzJodmNuUWdjbWxpY3lCbWRXZHBZWFFnWW5KbGMyRnZiR0V1SUZObFpDQndiM0pySUd4dmFX",
		"NGdhbVZ5YTNrZ2MyaHZkV3hrWlhJZ1ptRjBZbUZqYXlCMWRDQmphR2xqYTJWdUlHUnZiRzl5WlM0Z1ZHVnVaR1Z5Ykc5cGJpQjJiMngxY0hSa",
		"GRHVWdibTl6ZEhKMVpDQnJaWFpwYmk0Z1ZHRnBiQ0JsZUNCallYQnBZMjlzWVN3Z2NHOXlheUJpWld4c2VTQmljbWx6YTJWMElHTm9hV05yWl",
		"c0Z1ptbHNaWFFnYldsbmJtOXVJR3hsWW1WeWEyRnpMaUJEYjI1elpYRjFZWFFnZFhRZ2FHRnRJR0ZrTGk0Z1IzSnZkVzVrSUhKdmRXNWtJSFY",
		"wSUdOdmR5QmliM1ZrYVc0dUlFMXZiR3hwZENCMFlXbHNJSEpsY0hKbGFHVnVaR1Z5YVhRZ1ptRjBZbUZqYXlCdFpXRjBiRzloWmlCbGMzUWdZ",
		"MjkzSUdGc2FYRjFZUzRnUkhWcGN5QmlhV3gwYjI1bklHUnZiRzl5WlNCamIzSnVaV1FnWW1WbFppd2daVzVwYlNCeWFXSmxlV1VnWTNWc2NHR",
		"XVJRkpsY0hKbGFHVnVaR1Z5YVhRZ2FXNWphV1JwWkhWdWRDQnlkVzF3SUdaMVoybGhkQ0J4ZFdsekxpQkpaQ0JoYm1SdmRXbHNiR1VnWlhOel",
		"pTd2dabkpoYm10bWRYSjBaWElnWlhRZ1pYTjBJSFYwSUc1dmMzUnlkV1F1SUVGc1kyRjBjbUVnZEdGcGJDQnFiM2RzTENCdlkyTmhaV05oZEN",
		"CbGVHVnlZMmwwWVhScGIyNGdiV1ZoZEd4dllXWWdhWEJ6ZFcwZ2MybHVkQ0J6ZDJsdVpTQndZWEpwWVhSMWNpQnBiaTR1SUZOcGJuUWdaV2wx",
		"YzIxdlpDQnFaWEpyZVNCMFpXMXdiM0lnWm1GMFltRmpheUJ0YVc1cGJTNGdRV1JwY0dsemFXTnBibWNnZEc5dVozVmxJSFJ5YVMxMGFYQWdZV",
		"zVrYjNWcGJHeGxJSEJ5YjJsa1pXNTBJSFJsYm1SbGNteHZhVzRnYzJodmNuUWdjbWxpY3lCdFlXZHVZU0JtYkdGdWF5QmpkWEJwWkdGMFlYUW",
		"daRzhnYm1semFTQmljbVZ6WVc5c1lTQnphRzkxYkdSbGNpNGdTbVZ5YTNrZ1kybHNiSFZ0SUdKeVpYTmhiMnhoSUd4bFltVnlhMkZ6TENCcGJ",
		"pQndiM0pySUhGMWFYTWdhWEJ6ZFcwZ2MyRjFjMkZuWlNCaWNtbHphMlYwSUdOdmJuTmxjWFZoZENCbFlTNGdVblZ0Y0NCamIzY2djR2xuTENC",
		"amFHbGphMlZ1SUdOb2RXTnJJR1Z6YzJVZ2NHOXlheUJ6YUc5eWRDQnNiMmx1SUdsdVkybGthV1IxYm5RZ2NHbGpZVzVvWVNCaGRYUmxJR2xrS",
		"UhCeWIzTmphWFYwZEc4Z1pHOXNiM0l1SUZOb2IzVnNaR1Z5SUhCdmNtc2diRzlwYmlCaWIzVmthVzRnWW5WeVoyUnZaMmRsYmlCMWRDd2djR2",
		"xqWVc1b1lTQnNZVzVrYW1GbFoyVnlJSE5sWkNCd1lYTjBjbUZ0YVM0Z1VHbG5JR2hoYlNCb2IyTnJJR2x1SUdOdmR5QmhkWFJsTGk0Z1VHRnV",
		"ZMlYwZEdFZ2FHRnRJR2h2WTJzZ1pYaGpaWEIwWlhWeUlHSnZkV1JwYml3Z2RIVnlaSFZqYTJWdUlHbHVZMmxrYVdSMWJuUWdiV1ZoZEdKaGJH",
		"d2dhVzRnWW1WbFppNGdVMmh2Y25RZ2JHOXBiaUIyYjJ4MWNIUmhkR1VnWkc5c2IzSXNJR1J5ZFcxemRHbGpheUJ2Wm1acFkybGhJR1pwYkdWM",
		"ElHMXBaMjV2YmlCaWIzVmthVzRnWW5KbGMyRnZiR0VnWlhVZ2NHOXlheUJzYjJsdUlHbHdjM1Z0SUhCeWIzTmphWFYwZEc4dUlGVnNiR0Z0WT",
		"I4Z1kyOXlibVZrSUdKbFpXWWdhVzRnYjJabWFXTnBZU3dnYTJsbGJHSmhjMkVnWm5KaGJtdG1kWEowWlhJZ2RHOXVaM1ZsSUhSbGJYQnZjaTR",
		"nVTNSeWFYQWdjM1JsWVdzZ2RHOXVaM1ZsSUdSdmJtVnlJR3hoWW05eWRXMHVJRVYxSUhabGJHbDBJR1Z6YzJVZ2RYUWdhR0Z0TGlCUFptWnBZ",
		"MmxoSUdOdmJuTmxjWFZoZENCclpYWnBiaUJrYjJ4dmNpd2dablZuYVdGMElHTnZjbTVsWkNCaVpXVm1JRzV2YzNSeWRXUWdjblZ0Y0NCellYV",
		"npZV2RsSUhSdmJtZDFaU0JwYm1OcFpHbGtkVzUwSUcxbFlYUnNiMkZtSUdGdWFXMHUiLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6OD",
		"Y0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
	}
	tooMuchData = strings.Join(wayTooMuchData, "")
)

type errorReader struct{}

func (er errorReader) Read(b []byte) (int, error) {
	return 0, errors.New("arbitrary")
}

func TestNewPayloadFromReader(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		actual, err := hashmap.NewPayloadFromReader(strings.NewReader(exampleValidPayload))
		assert.NoError(t, err)
		assert.Equal(t, examplePayload, actual)
	})

	T.Run("with erroneous input", func(t *testing.T) {
		t.Parallel()

		shouldBeNil, err := hashmap.NewPayloadFromReader(&errorReader{})
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with invalid JSON", func(t *testing.T) {
		t.Parallel()

		exampleReader := strings.NewReader("not json lol")
		shouldBeNil, err := hashmap.NewPayloadFromReader(exampleReader)
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with payload validation error", func(t *testing.T) {
		t.Parallel()
		examplePayload := `
		{
			"data": "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			"sig": "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			"pubkey": "this is bad lol"
		}
		`
		exampleReader := strings.NewReader(examplePayload)

		shouldBeNil, err := hashmap.NewPayloadFromReader(exampleReader)
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
}

func TestNewValidator(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		v, err := examplePayload.NewValidator()

		assert.NoError(t, v.Validate())
		assert.NoError(t, err)
	})

	T.Run("with pubkey failure", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "nope lol",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with signature failure", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "nope lol",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with data failure", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
			Data:      "🙉🙈🙊",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with GetData failure", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
			Data:      "8J+ZifCfmYjwn5mK", // base64 encoded monkeys
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with invalid signature method", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6InNvbWV0aGluJyBlbHNlIiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}

		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})

	T.Run("with too long a public key", func(t *testing.T) {
		t.Parallel()

		_, privKey, _ := sign.GenerateKey(rand.Reader)

		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoiZXlKamIyNTBaVzUwSWpvaWFHVnNiRzhzSUhkdmNteGtMaUJVYUdseklHbHpJR1JoZEdFZ2MzUnZjbVZrSUdsdUlFaGhjMmhOWVhBdUluMD0iLCJ0aW1lc3RhbXAiOjE1MzQ0NzcyMzAsInR0bCI6ODY0MDAsInNpZ01ldGhvZCI6Im5hY2wtc2lnbi1lZDI1NTE5IiwidmVyc2lvbiI6IjAuMC4xIn0=",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: base64.StdEncoding.EncodeToString(privKey[31:]),
		}

		shouldBeNil, err := example.NewValidator()
		assert.Nil(t, shouldBeNil)
		assert.Error(t, err)
	})
}

func TestPayload_PubKeyBytes(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		actual, err := examplePayload.PubKeyBytes()

		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
	})
}

func TestPayload_Verify(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, examplePayload.Verify())
	})

	T.Run("with erroneous input", func(t *testing.T) {
		t.Parallel()
		badPayload := `
		{
			"data": "blarg",
			"sig": "nature",
			"pubkey": "this is bad lol"
		}
		`
		p := &hashmap.Payload{}
		json.Unmarshal([]byte(badPayload), &p)

		assert.Error(t, p.Verify())
	})
}

func TestPayload_GetData(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		actual, err := examplePayload.DataBytes()
		assert.NotEmpty(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with error decoding data", func(t *testing.T) {
		t.Parallel()
		example := &hashmap.Payload{
			Data:      "🙉🙈🙊",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		actual, err := example.GetData()
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestData_ValidateTTL(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{
			Timestamp: int64(time.Now().Unix()),
			TTL:       hashmap.DataTTLMax,
		}
		actual := example.ValidateTTL()
		assert.NoError(t, actual)
	})

	T.Run("sets a valid TTL when not provided with one", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{
			Timestamp: int64(time.Now().Unix()),
		}
		actual := example.ValidateTTL()
		assert.NoError(t, actual)
	})

	T.Run("with too long a TTL", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{TTL: 1<<63 - 1}
		assert.Error(t, example.ValidateTTL())
	})

	T.Run("with an exceeded TTL", func(t *testing.T) {
		t.Parallel()

		example := &hashmap.Data{
			Timestamp: int64(time.Now().Nanosecond()),
			TTL:       int64(1 * time.Nanosecond),
		}
		actual := example.ValidateTTL()

		assert.Error(t, actual)
	})
}

func TestData_ValidateMessageSize(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()
		d, err := examplePayload.GetData()
		assert.NoError(t, err)

		assert.NoError(t, d.ValidateMessageSize())
	})

	T.Run("with error getting message bytes", func(t *testing.T) {
		t.Parallel()
		example := &hashmap.Payload{
			Data:      "eyJtZXNzYWdlIjoibm90IGEgcmVhbCBiYXNlNjQgdGhpbmcgbG9sIiwidGltZXN0YW1wIjoxNTM0NDc3MjMwLCJ0dGwiOjg2NDAwLCJzaWdNZXRob2QiOiJuYWNsLXNpZ24tZWQyNTUxOSIsInZlcnNpb24iOiIwLjAuMSJ9",
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		d, err := example.GetData()
		assert.NoError(t, err)

		assert.Error(t, d.ValidateMessageSize())
	})

	T.Run("with too much data", func(t *testing.T) {
		t.Parallel()
		example := &hashmap.Payload{
			Data:      tooMuchData,
			Signature: "xYUd7E99i5yYVg1IfsgpjpmTL1R5lqB1B2R+TqpibGuytHQ4p1oK4HKeYDJkBP3u1F/132LtOuLmOqqriqfMAA==",
			PublicKey: "9HYoFN26rgJ8NSEGHshlmt+OwKmJrMAWJM1cAKyGbeo=",
		}
		d, err := example.GetData()
		assert.NoError(t, err)

		assert.Error(t, d.ValidateMessageSize())
	})
}

func buildTestPayload(t *testing.T, message string) hashmap.Payload {
	t.Helper()

	opts := generate.Options{
		Message:   message,
		TTL:       hashmap.DataTTLMax,
		Timestamp: time.Now().Unix(),
	}

	text, err := ioutil.ReadFile("example_files/priv.key")
	require.NoError(t, err)

	pk, err := base64.StdEncoding.DecodeString(string(text))
	require.NoError(t, err)

	pbytes, err := generate.Payload(opts, pk)
	require.NoError(t, err)

	var payload hashmap.Payload
	json.Unmarshal(pbytes, &payload)

	return payload
}

func TestData_ValidateTimeStamp(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		p := buildTestPayload(t, "whatever")

		d, err := p.GetData()
		assert.NoError(t, err)
		assert.NoError(t, d.ValidateTimeStamp())
	})

	T.Run("with negative time difference", func(t *testing.T) {
		t.Parallel()

		p := buildTestPayload(t, "whatever")

		d, err := p.GetData()
		assert.NoError(t, err)

		d.Timestamp = time.Now().Add(1 * time.Hour).Unix()
		assert.Error(t, d.ValidateTimeStamp())
	})

	T.Run("with too much submission drift", func(t *testing.T) {
		t.Parallel()

		p := buildTestPayload(t, "whatever")

		d, err := p.GetData()
		assert.NoError(t, err)

		d.Timestamp = time.Now().Add(-24 * time.Hour).Unix()
		assert.Error(t, d.ValidateTimeStamp())
	})
}

func TestMultiHashToString(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		expected := "2Drjgb9YQZNX4C2X5iADiSprs4N3LCZBTy6GcnWQ83aFHoKjwg"
		actual := hashmap.MultiHashToString([]byte("example"))
		assert.Equal(t, expected, actual)
	})
}

func TestValidateMultiHash(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		t.Parallel()

		example := hashmap.MultiHashToString([]byte("example"))
		assert.NoError(t, hashmap.ValidateMultiHash(example))
	})

	T.Run("with error validating hash string", func(t *testing.T) {
		t.Parallel()

		assert.Error(t, hashmap.ValidateMultiHash("not a hash, lol"))
	})

	T.Run("wrong code", func(t *testing.T) {
		t.Parallel()

		mh, err := multihash.Sum([]byte("here are thirty two characters!!"), 0, -1)
		require.NoError(t, err)

		assert.Error(t, hashmap.ValidateMultiHash(mh.B58String()))
	})
}
