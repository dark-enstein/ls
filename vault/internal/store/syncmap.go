package store

import (
	"context"
	"fmt"
	"github.com/dark-enstein/vault/internal/vlog"
	"sync"
)

type Map struct {
	scaffold *sync.Map
	logger   *vlog.Logger
}

//func NewSyncMap() *sync.Map {
//	return &sync.Map{}
//}

func NewSyncMap(ctx context.Context, logger *vlog.Logger) *Map {
	return &Map{
		&sync.Map{},
		logger,
	}
}

func (m *Map) Connect(ctx context.Context) (bool, error) {
	return true, nil
}

func (m *Map) Store(ctx context.Context, id string, token any) (bool, error) {
	log := m.logger.Logger()
	// check if key exists
	if m.IsExist(id) {
		log.Error().Msgf("key %s already exists, aborting\n", id)
		return false, fmt.Errorf("key %s already exists, aborting\n", id)
	}

	m.scaffold.Store(id, token)
	if _, ok := m.scaffold.Load(id); !ok {
		log.Debug().Msgf("error occurred while confirming insertion")
		return false, fmt.Errorf("error occurred while confirming insertion")
	}
	return true, nil
}

func (m *Map) Retrieve(ctx context.Context, id string) (string, error) {
	log := m.logger.Logger()
	val, ok := m.scaffold.Load(id)
	if !ok {
		log.Debug().Msgf("error occured while retrieving value from store using key id: %s\n", id)
		return "", fmt.Errorf("error occured while retrieving value from store using key id: %s\n", id)
	}

	return val.(string), nil
}

func (m *Map) RetrieveAll(ctx context.Context) (map[string]string, error) {
	log := m.logger.Logger()

	// create a bucket for all the tokens
	var allTokenMap = map[string]string{}
	// pass a range func over the contents of the store and get the contents
	m.scaffold.Range(func(id, value interface{}) bool {
		allTokenMap[fmt.Sprint(id)] = fmt.Sprint(value)
		return true
	})
	log.Debug().Msg("successfully ranged over sync.Map store")

	return allTokenMap, nil
}

func (m *Map) Delete(ctx context.Context, id string) (bool, error) {
	log := m.logger.Logger()

	// delete key from map
	m.scaffold.Delete(id)

	// check if key still exists
	if m.IsExist(id) {
		log.Error().Msgf("key with id %s did not delete successfully\n", id)
		return false, fmt.Errorf("key with id: %s did not delete successfully\n", id)
	}

	log.Debug().Msgf("successfully deleted key with id: %s\n", id)

	return true, nil
}

func (m *Map) Patch(ctx context.Context, id string, token any) (bool, error) {
	log := m.logger.Logger()

	// check if key exists
	if m.IsExist(id) {
		log.Error().Msgf("key with id exists, patching", id)
	}

	// patch key in map
	m.scaffold.Store(id, token)
	log.Debug().Msgf("successfully updated key with id: %s\n", id)

	return true, nil
}

func (m *Map) Flush(ctx context.Context) (bool, error) {

	// pass a range func over the contents of the store and delete the contents
	m.scaffold = &sync.Map{}

	return true, nil
}

func (m *Map) Close(ctx context.Context) error {
	return nil
}

// IsExist checks whether a key exists in a sync map
func (m *Map) IsExist(key string) bool {
	_, ok := m.scaffold.Load(key)
	return ok
}
