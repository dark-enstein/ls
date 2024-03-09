package manager

import (
	"context"
	"github.com/dark-enstein/sleight/internal/jury"
	"github.com/dark-enstein/sleight/internal/trace"
	"sync"
)

type Manager struct {
	pCount int8
	logger trace.Monitor
	// comm and management channels
	cmdChan  chan []byte
	exitChan chan struct{}
	// goroutine managments
	wg sync.WaitGroup
	rw sync.RWMutex
	w  sync.Mutex
}

func NewManager(m trace.Monitor) *Manager {
	log := m.Logger()

	log.Debug("creating Sleight Manager instance")
	return &Manager{
		pCount:   0,
		logger:   m,
		cmdChan:  make(chan []byte),
		exitChan: make(chan struct{}, 1),
		wg:       sync.WaitGroup{},
		rw:       sync.RWMutex{},
		w:        sync.Mutex{},
	}
}

// Run does the core of Manager tasks
func (m *Manager) Run(ctx context.Context, exitChan chan struct{}) int {
	// check and set up bin/state

	// spawn go routines and other admin tasks

	// return
	return jury.ErrSuccess
}
