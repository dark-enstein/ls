package trace

import "go.uber.org/zap"

type Monitor interface {
	Logger() *zap.Logger
	Level() string
}
