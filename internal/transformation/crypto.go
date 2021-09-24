package transformation

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
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
