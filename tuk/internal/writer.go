package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	tabWriter "text/tabWriter"
)

// Format represents an enum of the format of the output
type Format int

const (
	Tab Format = iota
	Json
)

// processorMap is the map of the format to the formatter
var processorMap = map[Format]Formatter{
	Tab:  NewTabWriter(),
	Json: NewJsonWriter(),
}

// EventData represents the fields of an event object
type EventData struct {
	File string `json:"file"`
	Type string `json:"type"`
}

// NewEventData creates a new event data object
func NewEventData(file, eventType string) *EventData {
	return &EventData{
		File: file,
		Type: eventType,
	}
}

// EventProcessor processes the events before rendering in the TUI
type EventProcessor struct {
	// formatter is the formatter to use for processing events
	formatter Formatter
	// buf is the buffer to write the formatted event output (by default, unbuffered, but buffer size can be set with WithEventBufferOptions)
	buf chan *bytes.Buffer
	// bufferSize is the size of the event buffer
	bufferSize int
}

// ProcessorOptions is the option to configure the event processor
type ProcessorOptions func(*EventProcessor)

// WithEventBufferOptions sets the buffer size of the event processor
func WithEventBufferOptions(size int) ProcessorOptions {
	return func(ep *EventProcessor) {
		ep.bufferSize = size
		ep.buf = make(chan *bytes.Buffer, size)
	}
}

// NewEventProcessor creates a new event processor instance
func NewEventProcessor(ctx context.Context) *EventProcessor {
	// init a buffered channel for queueing events
	return &EventProcessor{}
}

// WithProcessor sets the formatter to use for processing events
func (ep *EventProcessor) WithProcessor(procInt Format, opts ...ProcessorOptions) *EventProcessor {
	ep.formatter = processorMap[procInt]
	for i := range opts {
		opts[i](ep)
	}
	return ep
}

// Process processes the event
func (ep *EventProcessor) Process(file, eventType string) ([]byte, error) {
	return ep.formatter.Output(NewEventData(file, eventType))
}

// Formatter represents the interface for formatting the event data
type Formatter interface {
	Output(*EventData) ([]byte, error)
}

// TabWriter represents the tab writer formatter
type TabWriter struct {
	w   *tabWriter.Writer
	buf *bytes.Buffer
}

// NewTabWriter creates a new tab writer instance
func NewTabWriter() *TabWriter {
	tw := new(tabWriter.Writer)
	return &TabWriter{
		w:   tw.Init(os.Stdout, 0, 8, 8, '\t', 0),
		buf: new(bytes.Buffer),
	}
}

// Output formats the event data into a tab separated string
func (tw *TabWriter) Output(e *EventData) ([]byte, error) {
	tw.buf.Reset()
	fmt.Fprintf(tw.buf, "%s\t%s", e.Type, e.File)
	return tw.buf.Bytes(), nil
}

// JsonWriter represents the json writer formatter
type JsonWriter struct {
	buf *bytes.Buffer
}

// NewJsonWriter creates a new json writer instance
func NewJsonWriter() *JsonWriter {
	return &JsonWriter{
		buf: new(bytes.Buffer),
	}
}

// Output formats the event data into a json string
func (jw *JsonWriter) Output(e *EventData) ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
