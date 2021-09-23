package crypto

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
)

// EncodeDecode func to encode and decode request
func EncodeDecode(rawData []byte, operationType string) ([]byte, error) {
	aesBlock, err := aes.NewCipher([]byte(os.Getenv("CIPHER_KEY")))
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := []byte(os.Getenv("NONCE_KEY"))
	if err != nil {
		return nil, err
	}

	if operationType == "encode" {
		encodedData := aesGcm.Seal(nil, nonce, rawData, nil)
		return encodedData, nil
	} else if operationType == "decode" {
		decodedData, err := aesGcm.Open(nil, nonce, rawData, nil)
		if err != nil {
			return nil, err
		}
		return decodedData, nil
	}
	return nil, fmt.Errorf("can't make operation")
}

// Compress func to compress types.Stats
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

// Decompress func to decompress types.Stats
func Decompress(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			return
		}
	}(r)

	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
