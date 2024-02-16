package store

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"sync"
)

// TODO: implement btrees for storage

type Gob struct {
	loc string
	// basin is a temporary store, since it will be written to many times by separate method. It should be treated the same.
	basin  *Map
	fd     *os.File
	logger *vlog.Logger
	sync.RWMutex
}

func NewGob(ctx context.Context, loc string, logger *vlog.Logger) (*Gob, error) {
	log := logger.Logger()

	// gob encode
	err := IsValidFile(loc, log)
	if err != nil {
		log.Error().Msgf("gob store destination: %s invalid: error: %s\n", loc, err.Error())
		return nil, err
	}

	// open file
	fd, err := os.Create(loc)
	if err != nil {
		log.Info().Msgf("error while creating file at location %s: %s\n", loc, err.Error())
		return nil, err
	}
	return &Gob{loc, NewSyncMap(ctx, logger), fd, logger, sync.RWMutex{}}, nil
}

func (g *Gob) Connect(ctx context.Context) (bool, error) {
	return g.basin.Connect(ctx)
}

func (g *Gob) Store(ctx context.Context, id string, token any) error {
	log := g.logger.Logger()

	// first refresh in-memory map
	err := g.MapRefresh(ctx)
	if err != nil {
		log.Error().Msgf("error while refresh gob persistent storage: error: %s\n", err.Error())
		return err
	}

	// store new key value pair in the in-memory store
	err = g.basin.Store(ctx, id, token)
	if err != nil {
		return err
	}

	// persist the in-memory store to disk
	i, err := g.persist(ctx, false)
	if err != nil {
		return err
	}

	if i == 0 {
		return fmt.Errorf("wrote 0 bytes to gob persistent storage")
	}

	return nil
}

func (g *Gob) Patch(ctx context.Context, id string, token any) (bool, error) {
	log := g.logger.Logger()

	// first refresh in-memory map
	err := g.MapRefresh(ctx)
	if err != nil {
		log.Error().Msgf("error while refresh gob persistent storage: error: %s\n", err.Error())
		return false, err
	}

	// check if key exists
	b, err := g.basin.Patch(ctx, id, token)
	if err != nil || !b {
		log.Debug().Msgf("error with patching entry with id: %s\n", id)
		return false, fmt.Errorf("error with patching entry with id: %s\n", id)
	}

	// persist in-memory map
	i, err := g.persist(ctx, true)
	if err != nil {
		log.Debug().Msgf("error while persisting patched entry with id: %s\n", id)
		return false, fmt.Errorf("error while persisting patched entry with id: %s\n", id)
	}

	if i == 0 {
		return false, fmt.Errorf("wrote 0 bytes to gob persistent storage")
	}

	return true, nil
}

func (g *Gob) Retrieve(ctx context.Context, id string) (string, error) {
	log := g.logger.Logger()
	// first refresh in-memory map
	err := g.MapRefresh(ctx)
	if err != nil {
		log.Error().Msgf("error while refresh gob persistent storage: error: %s\n", err.Error())
		return "", err
	}

	// retrieve value if it exists in store map
	value, err := g.basin.Retrieve(ctx, id)
	if err != nil {
		log.Debug().Msgf("error while retrieving value with id: %s: %s\n", id, err.Error())
		return "", fmt.Errorf("error while retrieving value with id: %s: %s\n", id, err.Error())
	}

	return value, nil
}

func (g *Gob) RetrieveAll(ctx context.Context) (map[string]string, error) {
	log := g.logger.Logger()

	// first refresh in-memory map
	err := g.MapRefresh(ctx)
	if err != nil {
		log.Error().Msgf("error while refresh gob persistent storage: error: %s\n", err.Error())
		return nil, err
	}

	m, err := g.basin.RetrieveAll(ctx)
	if err != nil {
		log.Debug().Msgf("error while retrieving all entries: %s\n", err.Error())
		return nil, fmt.Errorf("error while retrieving all entries: %s\n", err.Error())
	}

	return m, err
}

func (g *Gob) Delete(ctx context.Context, id string) (bool, error) {
	log := g.logger.Logger()

	// first refresh in-memory map
	err := g.MapRefresh(ctx)
	if err != nil {
		log.Error().Msgf("error while refresh gob persistent storage: error: %s\n", err.Error())
		return false, err
	}

	// delete id from map
	b, err := g.basin.Delete(ctx, id)
	if err != nil || !b {
		log.Debug().Msgf("error while deleting entry with id: %s : %s\n", id, err.Error())
		return false, err
	}

	// replace store with new map
	i, err := g.persist(ctx, true)
	if err != nil {
		log.Debug().Msgf("error while persisting entry with id: %s : %s\n", id, err.Error())
		return false, err
	}

	if i == 0 {
		return false, fmt.Errorf("wrote 0 bytes to gob persistent storage")
	}

	return true, nil
}

func (g *Gob) trunc(i int64) error {
	g.Lock()
	defer g.Unlock()
	err := g.fd.Truncate(i)
	if err != nil {
		log.Info().Msgf("error while emptying store: %s\n", err.Error())
		return err
	}
	return nil
}

func (g *Gob) persist(ctx context.Context, replace bool) (int64, error) {
	log := g.logger.Logger()

	// pre-checks and pre-reqs
	if replace {
		// first, truncate file to zero. clean the file.
		err := g.trunc(0)
		if err != nil {
			log.Debug().Msgf("error while cleaning gob persistent store: %s\n", err.Error())
			return 0, err
		}

	}

	// core persist
	// dump
	err := g.MapDump(ctx)
	if err != nil {
		log.Error().Msgf("error while dumping in-memory map : error: %s\n", err.Error())
		return 0, err
	}

	// post checks
	// confirm bytes len written to file
	f, err := g.fd.Stat()
	if err != nil {
		log.Error().Msgf("error retrieving file stat: error: %s\n", err.Error())
		return 0, err
	}

	return f.Size(), nil
}

// Close closes the redis connection
func (g *Gob) Close(ctx context.Context) error {
	return g.fd.Close()
}

// MapRefresh refreshes the sync.Map in-memory store with the latest updates from the persistent store
func (g *Gob) MapRefresh(ctx context.Context) error {
	// sets up a temporary store for the
	var m = map[string]string{}

	// lock
	g.RLock()
	defer g.RUnlock()

	// set up gob decoder with file descriptor to gob persistent store
	dec := gob.NewDecoder(g.fd)

	// encode map and write to io.Writer
	err := dec.Decode(&m)
	if err != io.EOF && err != nil {
		log.Error().Msgf("error while encoding into map from gob persistent storage: error: %s\n", err.Error())
		return fmt.Errorf("error while encoding into map from gob persistent storage: error: %s\n", err.Error())
	}

	// first empty sync map
	// if any error is received from reading from file, the internal sync.Map isn't cleared
	_, err = g.basin.Flush(ctx)
	if err != nil {
		log.Error().Msgf("error flushing in-memory store: %s\n", err.Error())
		return fmt.Errorf("error flushing in-memory store: %s\n", err.Error())
	}

	// unfurl map into sync map
	syncM := g.basin.Map()
	for k, v := range m {
		syncM.Store(k, v)
	}

	return nil
}

// MapDump persists the current in-memory data to the persistent store
func (g *Gob) MapDump(ctx context.Context) error {

	g.Lock()
	defer g.Unlock()
	// sets up a temporary store for
	var m = map[string]string{}
	// set up gob encoder with file descriptor to gob persistent store
	dec := gob.NewEncoder(g.fd)

	// store all the sync map contents into the temporary map
	g.basin.scaffold.Range(func(id, value interface{}) bool {
		m[fmt.Sprint(id)] = fmt.Sprint(value)
		return true
	})
	log.Debug().Msg("successfully ranged over sync.Map store")

	// encode map and write to io.Writer || fd
	err := dec.Encode(&m)
	if err != nil {
		log.Error().Msgf("error while encoding into map from gob persistent storage: error: %s\n", err.Error())
		return err
	}
	log.Info().Msg("successfully persisted in-memory map to disk")

	return nil
}

// Flush empties the internal sync.Map and the persistent gob store
func (g *Gob) Flush(ctx context.Context) (bool, error) {
	log := g.logger.Logger()

	// first empty sync map
	b, err := g.basin.Flush(ctx)
	log.Debug().Msgf("flushing gob store")
	if !b || err != nil {
		log.Error().Msgf("error occurred while flushing gob store: %s\n", err.Error())
		return b, err
	}

	// delete persistent gob store last
	err = g.trunc(0)
	log.Debug().Msgf("flushing gob persistent store")
	if err != nil {
		log.Error().Msgf("error occurred while flushing persistent gob store: %s\n", err.Error())
		return false, err
	}
	return true, nil
}
