package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dark-enstein/sleight/internal/config"
	"github.com/dark-enstein/sleight/internal/jury"
	"github.com/dark-enstein/sleight/internal/manager"
	"github.com/dark-enstein/sleight/internal/trace"
	"github.com/spf13/pflag"
)

func _flags() (*config.Config, int) {
	lug := trace.Tmp()
	var c = config.Config{}

	pflag.StringVarP(&c.Level, config.FlagLevel, config.FlagShortLevel, config.InfoString, "one of the following log level: info, debug, warn, error, fatal, panic")
	pflag.Parse()

	if err := c.Validate(); err != nil {
		lug.Info(fmt.Sprintf("error validating flags: %v", err))
		return nil, jury.ErrRequest
	}

	return &c, jury.ErrSuccess
}

func main() {
	// dubidu
	fmt.Println(string(config.HeadingText))

	// init flags
	cfg, exit := _flags()
	if exit > jury.ErrSuccess {
		trace.Tmp().Fatal(fmt.Sprintln("exit:", exit))
	}
	llevel, _ := cfg.LogLevel()

	// init logger
	logger, err := trace.New(trace.Level(llevel))
	if err != nil {
		trace.Tmp().Fatal(fmt.Sprintln("exit:", jury.ErrInternal))
	}

	lug := logger.Logger()
	// init manager
	man := manager.NewManager(logger)
	lug.Debug("manager initialized")

	// exitChan for graceful shutdown
	exitChan := make(chan struct{}, 1)
	ctx := context.Background()

	lug.Debug("setting up Run")
	e := man.Run(ctx, exitChan)
	if e != jury.ErrSuccess {
		log.Fatalln("error running sleight. exiting")
	}
}
