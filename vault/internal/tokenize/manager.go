package tokenize

import (
	"context"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"sync"
	"time"
	"unsafe"
)

var (
	DefaultCipherLoc          = "./cipher"
	EnvKeyAESCipher           = "CIPHER"
	EnvKeInitializationVector = "IV"
)

type Manager struct {
	store     *sync.Map
	cipher    map[string]string
	cipherLoc string
}

func NewManager(cts context.Context, log *vlog.Logger) *Manager {
	var manager = Manager{}
	manager.cipherLoc = DefaultCipherLoc
	manager.cipher = map[string]string{}
	var err error

	// if cipher file doesn't exist
	if _, err := os.Stat(manager.cipherLoc); err != nil {
		log.Logger().Debug().Msgf("cipher file %s not found, generating it", manager.cipherLoc)
		err = manager.GenerateCipher()
		if err != nil {
			log.Logger().Error().Msgf("error encountered while writing cipher to file: %s\n", err.Error())
		}
	}

	// if cipher file already exists, emvMap is empty, so read from file
	if len(manager.cipher) == 0 {
		manager.cipher, err = godotenv.Read(manager.cipherLoc)
		if err != nil {
			log.Logger().Error().Msgf("error encountered while reading cipher from file %s: %s\n", manager.cipherLoc, err.Error())
		}
	}

	return &manager
}

// GenerateCipher generates a new AES cipher and Initialization Vector pais, and persists it to disk
func (m *Manager) GenerateCipher() error {
	// generate 32 digit key
	m.cipher[EnvKeyAESCipher] = genAlphaNumericString(32)
	// generate 16 digit initialization vector
	m.cipher[EnvKeInitializationVector] = genAlphaNumericString(16)
	// write to file
	return godotenv.Write(m.cipher, m.cipherLoc)
}

// gotten from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go#:~:text=%22Mimicing%22%20strings.Builder%20with%20package%20unsafe
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// genAlphaNumericString mimics strings.Builder with package unsafe. According to the author it is one of the pastest implementation of strings builder.
func genAlphaNumericString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
