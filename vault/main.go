package main

import (
	"context"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/dark-enstein/vault/service"
	"github.com/spf13/pflag"
)

type InitConfig struct {
	debug bool
}

func main() {
	var i = InitConfig{}
	pflag.BoolVarP(&i.debug, "debug", "d", true, "enable debug mode")
	pflag.Parse()

	// starting
	ctx := context.Background()
	logger := vlog.New(i.debug)
	srv := service.New(ctx, logger)
	if err := srv.Run(ctx); err != nil {
		logger.Logger().Fatal().Msgf("error while service is starting: %s\n", err.Error())
	}
}
