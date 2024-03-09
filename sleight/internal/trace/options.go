package trace

type Config func(logger *Zap) error

// WithLogFile sets the log file for the logger
func WithLogFile(s string) Config {
	return func(l *Zap) error {
		l.loc = s
		return nil
	}
}
