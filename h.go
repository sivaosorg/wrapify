package wrapify

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"

	"github.com/sivaosorg/unify4g"
)

// CalculateSize calculates the size of the marshaled data.
// It uses unify4g.MarshalN to marshal the data and returns the length of the resulting byte slice.
// If an error occurs during marshaling, it returns 0.
func CalculateSize(data any) int {
	_bytes, err := unify4g.MarshalN(data)
	if err != nil {
		return 0
	}
	return len(_bytes)
}

// Compress compresses the given data using gzip and encodes it in base64.
// It first marshals the data using unify4g.MarshalN, then compresses the resulting byte slice
// using gzip. The compressed data is then encoded in base64 and returned as a string.
// If any error occurs during marshaling or compression, it returns an empty string.
func Compress(data any) string {
	_bytes, err := unify4g.MarshalN(data)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(_bytes)
	if err != nil {
		return ""
	}
	err = gz.Close()
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
