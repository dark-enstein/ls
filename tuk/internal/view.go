package internal

import (
	"context"
	"fmt"
	"os"
	"strings"
	"tuk/internal/config"
)

const (
	Events = "events"
	Help   = "help"
	Logs   = "logs"
)

// Plane represents the plane of the TUI
type Plane struct {
	// Manager handles manages the events lifetime and the rendering on the TUI
	m *Manager
	// paths are the paths to watch
	ps []string
}

// PlaneOption is the option to configure the plane
type PlaneOption func(*Plane)

// WithPaths sets the paths to watch
func WithPaths(path string) PlaneOption {
	return func(p *Plane) {
		if splices := strings.Split(path, ","); len(splices) > 1 {
			p.ps = append(p.ps, splices...)
			return
		}
	}
}

// WithConfig sets the configuration for the plane
func WithConfig(config *config.Config) PlaneOption {
	return func(p *Plane) {
		// paths := make([]string, 0, len(config.Paths)+len(config.Args.Path))
		paths := []string{}
		// funnel config paths
		for _, path := range config.Paths {
			paths = append(paths, path.Raw)
		}

		// funnel args paths into slice
		if config.Args.Path != "" {
			splice := strings.Split(config.Args.Path, ",")
			paths = append(paths, splice...)
		}

		// consolidate paths
		p.ps = append(p.ps, paths...)
	}
}

// NewPlane initializes a new plane instance
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

// Listen starts listening for events
func (p *Plane) Listen() {
	p.m.Listen()
}

// Run starts the TUI
func (p *Plane) Run() {
	// Clear Stdout
	fmt.Fprintf(os.Stdout, "\033[2J")
	// Run
	p.m.Run(p.ps...)
}

// Close closes the plane and its components: watcher and processor
func (p *Plane) Close() {
	p.m.Close()
}

// Log logs a message to the TUI log view
func (p *Plane) Log(msg string, args ...string) {
	p.m.Log(fmt.Sprintf(msg, args))
}
