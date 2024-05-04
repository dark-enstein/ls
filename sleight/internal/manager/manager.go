package manager

import (
	"context"
	"sync"

	"github.com/dark-enstein/sleight/internal/jury"
	"github.com/dark-enstein/sleight/internal/trace"
	"github.com/dark-enstein/vault/pkg/store"
)

const (
	DefaultStateStore = "./.store"
)

type Manager struct {
	pCount int8
	logger trace.Monitor
	gob    *store.Gob
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
	lug := m.logger.Logger()
	var err error
	// check and set up bin/state
	lug.Info("running")

	// validate user actions
	m.runValidation()

	// connect to shared memory? mmap
	m.gob, err = store.NewGob(ctx, DefaultStateStore, nil, false)
	if err != nil {
		return 0
	}

	// unmarshall into struct and check state of system
	// determine the last/current state of the system, and simply carry on from there (define states of the system)

	// determine action

	// spawn go routines and other admin tasks // delegate this to life?

	// return
	return jury.ErrSuccess
}

// runValidation validates the user inputs, ensuring consistency
func (m *Manager) runValidation() {}
