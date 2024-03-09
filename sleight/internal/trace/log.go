package trace

import (
	"github.com/dark-enstein/sleight/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const (
	DefaultLogLevel = zapcore.InfoLevel
)

type Level int8

const (
	InfoLevel Level = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	InvalidLevel
)

func (l *Level) String() string {
	switch *l {
	case DebugLevel:
		return config.DebugString
	case WarnLevel:
		return config.WarnString
	case ErrorLevel:
		return config.ErrorString
	case FatalLevel:
		return config.FatalString
	case PanicLevel:
		return config.PanicString
	case InfoLevel:
		return config.InfoString
	default:
		return config.InvalidErrTypeString
	}
}

type Zap struct {
	loc   string
	level Level
	l     *zap.Logger
}

const (
	defaultLogLoc = "./log/s.log"
)

// New instantiates a new logger instance which logs to both stdout and to a file
func New(l Level, opts ...Config) (*Zap, error) {
	var r Zap

	// define defaults
	r.loc = defaultLogLoc
	r.level = l

	// apply options
	for i := 0; i < len(opts); i++ {
		if err := opts[i](&r); err != nil {
			return nil, err
		}
	}
	// TODO: also apply config from yaml

	// define file descriptors
	stdout := zapcore.AddSync(os.Stdout)

	zFd := zapcore.AddSync(&lumberjack.Logger{
		Filename:   r.loc,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	// design custom zap writers
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.TimeKey = "timestamp"
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.TimeKey = "timestamp"
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// resolve perfered logging level
	level := zap.NewAtomicLevelAt(resolveLogger(l))

	// complete zap config encoding
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, zFd, level),
	)

	// init zap instance
	r.l = zap.New(core)

	return &r, nil
}

// Level returns the configured logging level
func (l *Zap) Level() string {
	return l.level.String()
}

// resolveLogger resolves the user-configured log level for initializing the Zap
func resolveLogger(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// Logger returns the underlying zap.logger instance for use
func (l *Zap) Logger() *zap.Logger {
	return l.l
}

func Tmp() *zap.Logger {
	return zap.Must(zap.NewProduction())
}
