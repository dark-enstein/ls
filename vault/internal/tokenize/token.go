package tokenize

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/pkg/errors"
)

var (
	ErrCipherToken404AES = errors.New("aes cipher not found in cipher map")
	ErrCipherToken404IV  = errors.New("initialization vector (IV) not found in cipher map")
)

type Token struct {
	token string
}

func (t *Token) String() string {
	return t.token
}

func Tokenize(s string, cypher map[string]string) (*Token, error) {
	// resolve aes cipher and initialization vector
	var aesKey, iv string
	var ok bool
	if aesKey, ok = cypher[EnvKeyAESCipher]; !ok {
		return nil, ErrCipherToken404AES
	}

	if iv, ok = cypher[EnvKeInitializationVector]; !ok {
		return nil, ErrCipherToken404IV
	}

	// get request string padded bytes. Padding is done following PKCS #7: https://en.wikipedia.org/wiki/PKCS_7
	padded := getPaddedBlock(s)

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return nil, err
	}

	// generate cipher mode using cipher block
	mode := cipher.NewCBCEncrypter(block, []byte(iv))

	// create a byte block to hold the encrypted bytes. It will be the length of the padded request string
	encryptedBytes := make([]byte, len(padded))
	mode.CryptBlocks(encryptedBytes, padded)

	// encode encryptedBytes using base64 encoding
	token := base64.StdEncoding.EncodeToString(encryptedBytes)
	return &Token{token: token}, nil
}

// getPaddedBlock returns a properly padded bytes block such that it is works with AES encrypting requirements. The padding is done following PKCS #7: https://en.wikipedia.org/wiki/PKCS_7
func getPaddedBlock(s string) []byte {
	// get bytes representation and length of string. this is needed for block checking and AES encryption.
	sBytes := []byte(s)
	length := len(sBytes)

	// calculate the mod 16 of the bytes length, to determing how much is required to make the block a multiple of 16. AES standard. https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
	paddingRequired := 16 - length%16

	// following PKCS #7, the padding to be added (if needed) will be a repetition of the byte representation of the reminder
	// create a mod16 block
	sPaddedBlock := make([]byte, length+paddingRequired)
	// first copy the source bytes into the new mod16 block
	copy(sPaddedBlock, sBytes)

	padding := bytes.Repeat([]byte{byte(paddingRequired)}, paddingRequired)
	// then add the padded bytes to the same mod16 block
	sPaddedBlock = append(sPaddedBlock, padding...)
	return sPaddedBlock
}
