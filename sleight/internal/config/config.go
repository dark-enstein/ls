package config

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

func (c *Config) resolveLevel() int {
	switch c.Level {
	case DebugString:
		return DebugLevel
	case WarnString:
		return WarnLevel
	case ErrorString:
		return ErrorLevel
	case FatalString:
		return FatalLevel
	case PanicString:
		return PanicLevel
	case InfoString:
		return InfoLevel
	default:
		return InvalidLevel
	}
}

func (c *Config) Validate() error {
	// resolve tracer loglevel
	c.resolveLevel()
	return nil
}

func (c *Config) LogLevel() (int, error) {
	if !c.validated {
		// always validate first, if not validated yet
		if err := c.Validate(); err != nil {
			return InvalidLevel, err
		}
	}
	return c.resolveLevel(), nil
}

// duplicate of what's in trace/log
const (
	InfoLevel = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	InvalidLevel
)
