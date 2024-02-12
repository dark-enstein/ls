package tokenize

import (
	"context"
	"fmt"
	"github.com/dark-enstein/vault/internal/model"
	"github.com/dark-enstein/vault/internal/store"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists. not overriding")
	ErrKeyDoesNotExists = "key %s does not exist"
	ErrDuplicateKeys    = errors.New("key already exists in request. accepted only the first one")
)

var (
	DefaultCipherLoc           = "./.cipher"
	EnvKeyAESCipher            = "CIPHER"
	EnvKeyInitializationVector = "IV"
	KeyDelimiter               = "__"
)

type Manager struct {
	store     *sync.Map
	cipher    map[string]string
	cipherLoc string
	log       *vlog.Logger
}

// NewManager creates a new instance of Manager. It manages token operations (retrieval, storage, servicing) throughout the lifetime of the server.
func NewManager(ctx context.Context, log *vlog.Logger) *Manager {
	var manager = Manager{}
	manager.log = log
	manager.cipherLoc = DefaultCipherLoc
	manager.cipher = map[string]string{}
	manager.store = store.NewSyncMap()
	//manager.store = store.NewSyncMap()
	var err error

	// if cipher file doesn't exist
	if _, err := os.Stat(manager.cipherLoc); err != nil {
		manager.log.Logger().Debug().Msgf("cipher file %s not found, generating it", manager.cipherLoc)
		err = manager.GenerateCipher()
		if err != nil {
			manager.log.Logger().Error().Msgf("error encountered while writing cipher to file: %s\n", err.Error())
		}
	}

	// if cipher file already exists, emvMap is empty, so read from file
	if len(manager.cipher) == 0 {
		manager.cipher, err = godotenv.Read(manager.cipherLoc)
		if err != nil {
			manager.log.Logger().Error().Msgf("error encountered while reading cipher from file %s: %s\n", manager.cipherLoc, err.Error())
		}
	}

	return &manager
}

// GenerateCipher generates a new AES cipher and Initialization Vector pais, and persists it to disk
func (m *Manager) GenerateCipher() error {
	// generate 32 digit key
	m.cipher[EnvKeyAESCipher] = genAlphaNumericString(32)
	// generate 16 digit initialization vector
	m.cipher[EnvKeyInitializationVector] = genAlphaNumericString(16)
	// write to file
	return godotenv.Write(m.cipher, m.cipherLoc)
}

// GetTokenByID returns the token owned by a specific ID/Key
func (m *Manager) GetTokenByID(id string) (*model.Tokenize, error) {
	log := m.log.Logger()
	var tokenStr string
	// pass a range func over the contents of the store and get the contents
	if val, ok := m.store.Load(id); !ok {
		return nil, fmt.Errorf(ErrKeyDoesNotExists, id)
	} else {
		tokenStr = fmt.Sprint(val)
	}
	log.Debug().Msg("successfully ranged over sync.Map store")

	ss := strings.Split(id, KeyDelimiter)

	log.Debug().Msg("found token in store")
	return &model.Tokenize{
		ID: ss[0],
		Data: []model.Child{
			{
				Key:   ss[1],
				Value: tokenStr,
			},
		},
	}, nil
}

// GetAllTokens returns all tokens currently in the store
func (m *Manager) GetAllTokens() ([]*model.Tokenize, error) {
	log := m.log.Logger()
	allTokens := map[string]*model.Tokenize{}

	// create a bucket for all the tokens
	var allTokenMap = map[string]string{}
	// pass a range func over the contents of the store and get the contents
	m.store.Range(func(key, value interface{}) bool {
		allTokenMap[fmt.Sprint(key)] = fmt.Sprint(value)
		return true
	})
	log.Debug().Msg("successfully ranged over sync.Map store")

	// parse all tokens into a slice of model.Tokenize
	for k, v := range allTokenMap {
		ss := strings.Split(k, KeyDelimiter)
		if val, ok := allTokens[k]; ok {
			if len(val.ID) == 0 {
				log.Debug().Msgf("token with id %s is already stored. continuing.", val.ID)
			}
			val.ID = ss[0]
			val.Data = append(val.Data, model.Child{
				Key:   ss[1],
				Value: v,
			})
			continue
		}
		allTokens[k] = &model.Tokenize{
			ID: ss[0],
		}
		allTokens[k].Data = append(allTokens[k].Data, model.Child{
			Key:   ss[1],
			Value: v,
		})
	}

	respTokens := []*model.Tokenize{}
	for _, v := range allTokens {
		respTokens = append(respTokens, v)
	}

	log.Debug().Msg("successfully parsed all tokens into a token array")
	return respTokens, nil
}

// ValidateResponse holds the error response from validation and the associated key.
type ValidateResponse struct {
	Key string
	Err error
}

// Validate is the high level api for validating all the user provided data
func (m *Manager) Validate(token *model.Tokenize) ([]*ValidateResponse, bool) {
	keysValidationResp, ok := m.ValidateKeys(token)
	if !ok {
		m.log.Logger().Error().Msgf("error while validating keys")
		return keysValidationResp, ok
	}

	//valValidationResp, ok := m.ValidateValues(token)
	return keysValidationResp, true
}

// ValidateKeys validates the Keys used in the request, ensuring it doesn't already exist, and that it conforms with the standards.
func (m *Manager) ValidateKeys(token *model.Tokenize) ([]*ValidateResponse, bool) {
	tempMap := make(map[string]bool, len(token.Data))
	valResp := []*ValidateResponse{}
	var verdict = true
	parentKey := token.ID
	for i := 0; i < len(token.Data); i++ {
		childKey := token.Data[i].Key
		combinedKeyName := GetCombinedKey(parentKey, childKey)
		// check that key doesn't already exist
		var err error
		if tempMap, err = keysIsPresent(combinedKeyName, tempMap, m.store); err != nil {
			verdict = false
			valResp = append(valResp, &ValidateResponse{combinedKeyName, err})
		}
	}
	return valResp, verdict
}

func keysIsPresent(key string, tempStore map[string]bool, store *sync.Map) (map[string]bool, error) {
	if _, ok := tempStore[key]; ok {
		return tempStore, ErrDuplicateKeys
	}
	if _, ok := store.Load(key); ok {
		return tempStore, ErrKeyAlreadyExists
	}
	tempStore[key] = true
	return tempStore, nil
}

// Tokenize manages the tokenization, and stores generated tokens in an internal store, for easy retrieval
func (m *Manager) Tokenize(key, val string) (string, error) {
	if _, ok := m.store.Load(key); ok {
		m.log.Logger().Error().Msg(ErrKeyAlreadyExists.Error())
		return "", ErrKeyAlreadyExists
	}
	token, err := tokenize(val, m.cipher)
	if err != nil {
		m.log.Logger().Error().Msgf("error occured while generating token: %s\n", err.Error())
		return "", err
	}
	m.store.Store(key, token.token)
	return token.token, nil
}

// Detokenize retrieves the value represented by a particular token, identified by the particular key
func (m *Manager) Detokenize(key, val string) (bool, string, error) {
	// check if key exists
	if _, ok := m.store.Load(key); !ok {
		m.log.Logger().Error().Msg("token with key id does not exist")
		return false, "", errors.New("token with key id does not exist")
	}

	decryptedStr, err := detokenize(val, m.cipher)
	if err != nil {
		m.log.Logger().Error().Msgf("error occured while decrypting token: %s\n", err.Error())
		return false, "", err
	}
	m.store.Store(key, decryptedStr)
	return true, decryptedStr, nil
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

// IsErrKeyAlreadyExist enables easy checking of error
func IsErrKeyAlreadyExist(err error) bool {
	if err == ErrKeyAlreadyExists {
		return true
	}
	return false
}

// GetCombinedKey creates a key string unique to every value in the request object. This key string is a concatenation of all the parent keys that constitute the request data
func GetCombinedKey(s ...string) (cs string) {
	for i := 0; i < len(s); i++ {
		delimiter := ""
		if i > 0 {
			delimiter = KeyDelimiter
		}
		cs += delimiter + s[i]
	}
	return
}
