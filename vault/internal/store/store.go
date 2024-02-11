package store

import "context"

type Store interface {
	Connect(ctx context.Context) (bool, error)
	Store(id string, token any) error
	Retrieve(id string) (string, error)
	RetrieveAll() (map[string]string, error)
	Delete(id string) (bool, error)
	Patch(id string, token any) (bool, error)
	Close(ctx context.Context) error
}
