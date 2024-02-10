package store

import "sync"

func NewSyncMap() *sync.Map {
	return &sync.Map{}
}
