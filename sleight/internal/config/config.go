package config

import "github.com/dark-enstein/sleight/internal/trace"

const (
	FlagLevel      = "level"
	FlagShortLevel = "l"
)

const (
	DebugString          = "debug"
	WarnString           = "warn"
	ErrorString          = "error"
	FatalString          = "fatal"
	PanicString          = "panic"
	InfoString           = "info"
	InvalidErrTypeString = ""
)

type Config struct {
	Level     string
	validated bool
}

func (c *Config) resolveLevel() trace.Level {
	switch c.Level {
	case DebugString:
		return trace.DebugLevel
	case WarnString:
		return trace.WarnLevel
	case ErrorString:
		return trace.ErrorLevel
	case FatalString:
		return trace.FatalLevel
	case PanicString:
		return trace.PanicLevel
	case InfoString:
		return trace.InfoLevel
	default:
		return trace.InvalidLevel
	}
}

func (c *Config) Validate() error {
	// resolve tracer loglevel
	c.resolveLevel()
	return nil
}

func (c *Config) LogLevel() (trace.Level, error) {
	if !c.validated {
		// always validate first, if not validated yet
		if err := c.Validate(); err != nil {
			return trace.InvalidLevel, err
		}
	}
	return c.resolveLevel(), nil
}
