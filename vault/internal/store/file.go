package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type File struct {
	loc      string
	fd       *os.File
	channels *FileChannels
	logger   *vlog.Logger
	sync.Mutex
}

type FileChannels struct {
	ops chan func(is string) error
}

// NewFile creates a new filestore at loc
func NewFile(loc string, logger *vlog.Logger) *File {
	return &File{
		loc:    loc,
		logger: logger,
	}
}

// Connect attempts to open a filestream to the file at location loc
func (f *File) Connect(ctx context.Context) (bool, error) {
	var abs string
	var err error
	log := f.logger.Logger()
	loc := f.loc

	// some sanity check
	if loc[0] != '/' {
		abs, err = filepath.Abs(loc)
		if err != nil {
			log.Error().Msgf("cannot obtain absolute path to referenced path: %s\n", err.Error())
			return false, fmt.Errorf("cannot obtain absolute path to referenced path: %s\n", err.Error())
		}
	}

	// confirm dir of path exists
	dir := filepath.Dir(abs)
	if _, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			log.Info().Msgf("base dir %s does not exist, creating it\n", dir)
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				log.Info().Msgf("error while creating base dir %s: %s\n", dir, err.Error())
				return false, err
			}
		} else {
			return false, err
		}
	}

	// create file database
	f.fd, err = os.Create(loc)
	if err != nil {
		log.Info().Msgf("error while creating file at location %s: %s\n", loc, err.Error())
		return false, err
	}
	return true, nil
}

// Close closes an open file
func (f *File) Close(ctx context.Context) error {
	return f.fd.Close()
}

// Store persists a new key-value entry in the file store
func (f *File) Store(id string, token any) error {
	log := f.logger.Logger()
	var err error

	// read current contents of the file
	content, err := f.read()
	if err != nil {
		log.Error().Msgf("error encountered while reading from file store: %s\n", err.Error())
		return fmt.Errorf("error encountered while reading from file store: %s\n", err.Error())
	}

	var storeMap = map[string]string{}

	// harvest currently stored values if file store is not empty
	if len(content) > 0 {
		storeMap, err = godotenv.UnmarshalBytes(content)
	}

	var tokenStr string

	// ensure that token type is string
	switch t := token.(type) {
	case string:
		tokenStr = fmt.Sprintf("%s", t)
	default:
		log.Error().Msgf("token of type string required: %s\n", err.Error())
		return fmt.Errorf("token of type string required: %s\n", err.Error())
	}

	// check if ID already exists
	if _, ok := storeMap[id]; ok {
		log.Error().Msgf("key already exists in store, skipping")
		return fmt.Errorf("key already exists in store, skipping")
	}

	// add the new ID
	storeMap[id] = tokenStr

	// write map to file store
	err = f.Write(storeMap)
	if err != nil {
		log.Error().Msgf("error while writing map to file store")
		return fmt.Errorf("error while writing map to file store")
	}

	return nil
}

// read abstracts away the details of reading from a file
func (f *File) read() ([]byte, error) {
	var b bytes.Buffer
	log := f.logger.Logger()
	f.Lock()
	defer f.Unlock()
	_, err := f.fd.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	i, err := b.ReadFrom(f.fd)
	if i != int64(b.Len()) {
		log.Debug().Msgf("number of read bytes %d does not match length in the bytes buffer %d\n", i, b.Len())
		return b.Bytes(), err
	}

	return b.Bytes(), err
}

// Retrieve retrieves a token from the store identified by id
func (f *File) Retrieve(id string) (string, error) {
	log := f.logger.Logger()
	var ok bool

	// read current contents of the file
	content, err := f.read()
	if err != nil {
		log.Error().Msgf("error encountered while reading from file store: %s\n", err.Error())
		return "", fmt.Errorf("error encountered while reading from file store: %s\n", err.Error())
	}

	var storeMap = map[string]string{}

	// if file is empty return empty
	if len(content) == 0 {
		log.Debug().Msgf("token with id %s doesn't exist", id)
		return "", fmt.Errorf("token with id %s doesn't exist", id)
	}

	storeMap, err = godotenv.UnmarshalBytes(content)
	// check err
	if err != nil {
		log.Debug().Msg("error while unmarshalling file store bytes")
		return "", errors.New("error while unmarshalling file store bytes")
	}

	var tokenStr string

	// check if ID already exists
	if tokenStr, ok = storeMap[id]; !ok {
		log.Debug().Msgf("token with id %s doesn't exist", id)
		return "", fmt.Errorf("token with id %s doesn't exist", id)
	}

	return tokenStr, nil
}

// RetrieveAll retrieves all the tokens from the store
func (f *File) RetrieveAll() (map[string]string, error) {
	log := f.logger.Logger()

	// read current contents of the file
	content, err := f.read()
	if err != nil {
		log.Error().Msgf("error encountered while reading from file store: %s\n", err.Error())
		return nil, fmt.Errorf("error encountered while reading from file store: %s\n", err.Error())
	}

	var storeMap = map[string]string{}

	// if file is empty return empty
	if len(content) == 0 {
		log.Debug().Msg("file store empty")
		return nil, errors.New("file store empty")
	}

	storeMap, err = godotenv.UnmarshalBytes(content)
	// check err
	if err != nil {
		log.Debug().Msg("error while unmarshalling file store bytes")
		return nil, errors.New("error while unmarshalling file store bytes")
	}

	return storeMap, nil
}

// Delete removes a token from the file store
func (f *File) Delete(id string) (bool, error) {
	log := f.logger.Logger()

	// read current contents of the file
	content, err := f.read()
	if err != nil {
		log.Error().Msgf("error encountered while reading from file store: %s\n", err.Error())
		return true, fmt.Errorf("error encountered while reading from file store: %s\n", err.Error())
	}

	var storeMap = map[string]string{}

	// if file is empty return empty
	if len(content) == 0 {
		log.Debug().Msg("file store empty")
		return true, errors.New("file store empty")
	}

	storeMap, err = godotenv.UnmarshalBytes(content)
	// check err
	if err != nil {
		log.Debug().Msg("error while unmarshalling file store bytes")
		return true, errors.New("error while unmarshalling file store bytes")
	}

	// delete from store
	delete(storeMap, id)

	// write to file
	// write map to file store
	err = f.Write(storeMap)
	if err != nil {
		log.Error().Msgf("error while writing map to file store")
		return false, fmt.Errorf("error while writing map to file store")
	}

	return true, nil
}

// Patch only updates a token in the file store, identified by id
func (f *File) Patch(id string, token any) (bool, error) {
	log := f.logger.Logger()

	// read current contents of the file
	content, err := f.read()
	if err != nil {
		log.Error().Msgf("error encountered while reading from file store: %s\n", err.Error())
		return true, fmt.Errorf("error encountered while reading from file store: %s\n", err.Error())
	}

	var storeMap = map[string]string{}

	// if file is empty return empty
	if len(content) == 0 {
		log.Debug().Msg("file store empty")
		return true, errors.New("file store empty")
	}

	storeMap, err = godotenv.UnmarshalBytes(content)
	// check err
	if err != nil {
		log.Debug().Msg("error while unmarshalling file store bytes")
		return true, errors.New("error while unmarshalling file store bytes")
	}

	var tokenStr string

	// ensure that token type is string
	switch t := token.(type) {
	case string:
		tokenStr = fmt.Sprintf("%s", t)
	default:
		log.Error().Msgf("token of type string required: %s\n", err.Error())
		return false, fmt.Errorf("token of type string required: %s\n", err.Error())
	}

	// check if ID exists
	if _, ok := storeMap[id]; !ok {
		log.Debug().Msgf("token with id %s doesn't exist", id)
		return false, fmt.Errorf("token with id %s doesn't exist", id)
	}

	// edit map
	storeMap[id] = tokenStr

	// write map to file store
	err = f.Write(storeMap)
	if err != nil {
		log.Error().Msgf("error while writing map to file store")
		return false, fmt.Errorf("error while writing map to file store")
	}

	return true, nil
}

func (f *File) Loop() {

}

// Write persists the map to disk using godotenv.Write
func (f *File) Write(m map[string]string) error {
	f.Lock()
	err := godotenv.Write(m, f.loc)
	f.Unlock()
	return err
}
