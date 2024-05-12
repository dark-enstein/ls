package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	tabWriter "text/tabWriter"
)

type Format int

const (
	Tab Format = iota
	Json
)

var processorMap = map[Format]Formatter{
	Tab:  NewTabWriter(),
	Json: NewJsonWriter(),
}

type EventData struct {
	File string `json:"file"`
	Type string `json:"type"`
}

func NewEventData(file, eventType string) *EventData {
	return &EventData{
		File: file,
		Type: eventType,
	}
}

type EventProcessor struct {
	// formatter is the formatter to use for processing events
	formatter Formatter
	// buf is the buffer to write the formatted event output (by default, unbuffered, but buffer size can be set with WithEventBufferOptions)
	buf chan *bytes.Buffer
	// bufferSize is the size of the event buffer
	bufferSize int
}

type ProcessorOptions func(*EventProcessor)

func WithEventBufferOptions(size int) ProcessorOptions {
	return func(ep *EventProcessor) {
		ep.bufferSize = size
		ep.buf = make(chan *bytes.Buffer, size)
	}
}

func NewEventProcessor(ctx context.Context) *EventProcessor {
	// init a buffered channel for queueing events
	return &EventProcessor{}
}

func (ep *EventProcessor) WithProcessor(procInt Format, opts ...ProcessorOptions) *EventProcessor {
	ep.formatter = processorMap[procInt]
	for i := range opts {
		opts[i](ep)
	}
	return ep
}

func (ep *EventProcessor) Process(file, eventType string) ([]byte, error) {
	return ep.formatter.Output(NewEventData(file, eventType))
}

type Formatter interface {
	Output(*EventData) ([]byte, error)
}

type TabWriter struct {
	w   *tabWriter.Writer
	buf *bytes.Buffer
}

func NewTabWriter() *TabWriter {
	tw := new(tabWriter.Writer)
	return &TabWriter{
		w:   tw.Init(os.Stdout, 0, 8, 8, '\t', 0),
		buf: new(bytes.Buffer),
	}
}

func (tw *TabWriter) Output(e *EventData) ([]byte, error) {
	tw.buf.Reset()
	fmt.Fprintf(tw.buf, "%s\t%s", e.Type, e.File)
	return tw.buf.Bytes(), nil
}

type JsonWriter struct {
	buf *bytes.Buffer
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{
		buf: new(bytes.Buffer),
	}
}

func (jw *JsonWriter) Output(e *EventData) ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
