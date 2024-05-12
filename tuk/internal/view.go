package internal

import (
	"context"
	"fmt"
	"strings"
)

const (
	Events = "events"
	Help   = "help"
	Logs   = "logs"
)

type Plane struct {
	// Manager handles manages the events lifetime and the rendering on the TUI
	m *Manager
	// paths are the paths to watch
	ps []string
}

type PlaneOption func(*Plane)

func WithPaths(path string) PlaneOption {
	return func(p *Plane) {
		if splices := strings.Split(path, ","); len(splices) > 1 {
			p.ps = append(p.ps, splices...)
			return
		}
	}
}

func NewPlane(ctx context.Context, opts ...PlaneOption) (*Plane, error) {
	// Initialize a new manager instance
	m, err := NewManager(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize the watcher and processor
	err = m.Init()
	if err != nil {
		return nil, err
	}

	p := &Plane{
		m: m,
	}

	for i := range opts {
		opts[i](p)
	}

	return p, err
}

func (p *Plane) Listen() {
	p.m.Listen()
}

func (p *Plane) Run() {
	p.m.Run(p.ps...)
}

func (p *Plane) Close() {
	p.m.Close()
}

func (p *Plane) Log(msg string, args ...string) {
	p.m.Log(fmt.Sprintf(msg, args))
}
