package wrapify

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"

	"github.com/sivaosorg/unify4g"
)

// calculateSize calculates the size of the marshaled data.
// It uses unify4g.MarshalN to marshal the data and returns the length of the resulting byte slice.
// If an error occurs during marshaling, it returns 0.
func calculateSize(data any) int {
	_bytes, err := unify4g.MarshalN(data)
	if err != nil {
		return 0
	}
	return len(_bytes)
}

// compress compresses the given data using gzip and encodes it in base64.
// It first marshals the data using unify4g.MarshalN, then compresses the resulting byte slice
// using gzip. The compressed data is then encoded in base64 and returned as a string.
// If any error occurs during marshaling or compression, it returns an empty string.
func compress(data any) string {
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

// chunk takes a response represented as a map and returns a slice of byte slices,
// where each byte slice is a chunk of the JSON representation of the response.
// This is useful for streaming large responses in smaller segments.
// If the JSON encoding fails, it returns nil.
func chunk(data map[string]any) [][]byte {
	_bytes, err := unify4g.MarshalN(data)
	if err != nil {
		return nil
	}
	var chunks [][]byte
	for i := 0; i < len(_bytes); i += DefaultChunkSize {
		end := i + DefaultChunkSize
		if end > len(_bytes) {
			end = len(_bytes)
		}
		// Create a copy of the chunk to avoid referencing the underlying array.
		// This is important to ensure that each chunk is independent and can be
		// processed separately without affecting the others.
		chunk := make([]byte, end-i)
		copy(chunk, _bytes[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}
