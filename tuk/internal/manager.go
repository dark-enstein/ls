package internal

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/jroimartin/gocui"
)

// Manager handles manages the events retrieval and the rendering on the TUI
type Manager struct {
	// g is the gocui instance
	g *gocui.Gui
	// LogChan is the channel for processing logs
	LogChan chan string
	// EventChan is the channel for processing events
	EventChan chan fmt.Stringer
	// Event processor processes the events before rendering in the TUI
	Processor *EventProcessor
	// Watcher is the fsnotify watcher
	w *fsnotify.Watcher
	// Paths are the list of paths to watch
	paths []Path
}

// Path represents a path to watch
type Path struct {
	raw      string
	resolved string
}

func NewPath(raw string) *Path {
	return &Path{
		raw: raw,
	}
}

func (p *Path) Validate() error {
	// Resolve the path
	abs, err := filepath.Abs(p.raw)
	if err != nil {
		return err
	}
	_, err = os.Stat(abs)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exist", abs)
	} else if err != nil {
		return fmt.Errorf("failed to access path %s, due to error: %v", abs, err)
	}

	p.resolved = abs
	return nil
}

func (p *Path) Name() string {
	return p.resolved
}

var (
	gManagers = []gocui.Manager{&DefaultManager{}}
)

// NewManager initializes a new Manager instance
func NewManager(ctx context.Context) (*Manager, error) {
	var err error
	m := &Manager{
		LogChan:   make(chan string),
		EventChan: make(chan fmt.Stringer),
	}
	m.g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	// Initialize keybindings
	err = initKeybindings(m.g)
	if err != nil {
		return nil, err
	}

	// Set managers
	m.SetManagers()

	return m, nil
}

func (m *Manager) Init() error {
	err := m.initWatcher(context.Background())
	if err != nil {
		return err
	}

	m.initProcessor(context.Background())
	return nil
}

func (m *Manager) Listen() {
	m.listenOutput()
	m.listenFileEvents()
}

func (m *Manager) Close() {
	// Close the watcher
	m.w.Close()
	// Close the gocui instance
	m.g.Close()
}

type DefaultManager struct {
}

func (m *DefaultManager) Name() string {
	return "default"
}

func (m *DefaultManager) Layout(g *gocui.Gui) error {
	// Set the layout for the manager
	maxX, maxY := g.Size()
	if v, err := g.SetView(Events, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Events"
		v.Autoscroll = true
	}

	if h, err := g.SetView(Help, maxX-25, 0, maxX-1, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(h, "File system events will be displayed here.")
		fmt.Fprintln(h, "Browse filesystem events using arrow keys.")
		h.Wrap = true
		h.Title = "Help"
		h.FgColor = gocui.ColorGreen
	}

	if h, err := g.SetView(Logs, maxX-50, 10, maxX-1, 25); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		h.Wrap = true
		h.Title = "Logs"
		h.FgColor = gocui.ColorYellow
		h.Autoscroll = true
	}
	return nil
}

func (m *Manager) SetManagers() {
	m.g.SetManager(gManagers...)
}

func (m *Manager) Update(view, msg string, clear bool) {
	m.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return err
		}
		if clear {
			v.Clear()
		}
		fmt.Fprintln(v, msg)
		return nil
	})
}

func (m *Manager) Run(paths ...string) error {
	m.AddPaths(paths...)
	m.Log("Running manager main loop")
	if err := m.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (m *Manager) AddPaths(paths ...string) {
	// m.Log("Adding paths: %v", paths...)
	for i := 0; i < len(paths); i++ {
		err := m.AddPath(paths[i])
		if err != nil {
			m.Log("Error adding path %s: %v", paths[i], err)
		}
	}
}

func (m *Manager) Log(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	m.LogChan <- msg
}

func (m *Manager) Publish(file, eventType string) {
	m.Log("Publishing event from: %s", file)
	buf, err := m.Processor.Process(file, eventType)
	if err != nil {
		m.Log("Error processing event: %v", err)
		return
	}
	m.Log("Successful publishing")
	m.EventChan <- bytes.NewBuffer(buf)
}

func (m *Manager) initWatcher(ctx context.Context) error {
	// Initialize fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		m.Log("error:", err.Error())
		return err
	}

	m.w = watcher
	return nil
}

func (m *Manager) AddPath(path string) error {
	m.Log("Adding path: %s", path)
	newPath := NewPath(path)
	err := newPath.Validate()
	if err != nil {
		m.Log("Error validating path %s: %v", path, err.Error())
		return err
	}
	m.Log("Successfully validated path: %s", path)
	m.w.Add(newPath.resolved)
	m.paths = append(m.paths, *newPath)
	m.Log("Successfully added path to watcher: %s", path)
	return nil
}

func (m *Manager) listenFileEvents() {
	go func() {
		for {
			select {
			case event := <-m.w.Events:
				m.Publish(event.Name, event.Op.String())
			case err := <-m.w.Errors:
				m.Log("error:", err.Error())
			}
		}
	}()
}

func (m *Manager) initProcessor(ctx context.Context) {
	// Initialize fsnotify watcher
	processor := NewEventProcessor(ctx).WithProcessor(Tab)

	m.Processor = processor
}

func (m *Manager) listenOutput() {
	go func() {
		for {
			select {
			case e := <-m.EventChan:
				m.Update(Events, e.String(), false)
			case l := <-m.LogChan:
				m.Update(Logs, l, false)
			}
		}
	}()
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(Events, gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(Events, gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	return nil
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
