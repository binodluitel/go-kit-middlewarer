package encoding

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// Encoder is anything that, given an interface, can store an encoding of the
// structure passed into Encode.
type Encoder interface {
	// Encode takes an interface, and should be able to translate it within the
	// given encoding, or it will fail with an error.
	Encode(interface{}) error
}

// GenerateEncoder is a function which takes an io.Writer, and returns an
// Encoder
type GenerateEncoder func(w io.Writer) Encoder

// MakeRequestEncoder takes a GenerateEncoder and returns an
// httptransport.EncodeRequestFunc
func MakeRequestEncoder(gen GenerateEncoder) httptransport.EncodeRequestFunc {
	return func(r *http.Request, request interface{}) error {
		var buf bytes.Buffer
		err := gen(&buf).Encode(request)
		r.Body = ioutil.NopCloser(&buf)
		return err
	}
}

// MakeResponseEncoder takes a GenerateEncoder and returns an
// httpstransport.EncodeResponseFunc
func MakeResponseEncoder(gen GenerateEncoder) httptransport.EncodeResponseFunc {
	return func(w http.ResponseWriter, response interface{}) error {
		return gen(w).Encode(response)
	}
}