package service

import (
	"context"
	"github.com/dark-enstein/vault/internal/store"
	"github.com/dark-enstein/vault/internal/tokenize"
	"github.com/dark-enstein/vault/internal/vlog"
	"net/http"
	"time"
)

const (
	port = "8080"
)

type Service struct {
	sc      *StartConfig
	manager *tokenize.Manager
	srv     *http.Server
	mux     *http.ServeMux
	log     *vlog.Logger
}

func New(ctx context.Context, log *vlog.Logger) *Service {
	log.Logger().Debug().Msg("generating service config")
	readTimeout := 10 * time.Second
	writeTimeout := 10 * time.Second
	hs := &http.Server{
		Addr:                         ":" + port,
		Handler:                      nil,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  readTimeout,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 writeTimeout,
		IdleTimeout:                  0,
		MaxHeaderBytes:               1 << 20,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	store, err := store.NewRedis(store.DefaultRedisConnectionString, log)
	if err != nil {
		log.Logger().Debug().Msgf("could not initialize redis backend: error: %s\n", err.Error())
	}

	log.Logger().Debug().Msgf("initialized service with settings:\n\taddress: %v\n\tread timeout: %v\n\twrite timeout: %v\n", ":"+port, readTimeout, writeTimeout)
	return &Service{sc: &StartConfig{port: port}, manager: tokenize.NewManager(ctx, log, tokenize.WithStore(store)), srv: hs, mux: http.NewServeMux(), log: log}
}

type StartConfig struct {
	port string
}

func (s *Service) Port() string {
	return s.sc.port
}

func (s *Service) LoadHandlers(ctx context.Context) {
	for k, v := range *NewVaultHandler(ctx, s) {
		s.mux.HandleFunc(k, v)
	}
}

func (s *Service) Run(ctx context.Context) error {
	// load handlers into mux
	s.LoadHandlers(ctx)
	// set mux into server
	s.srv.Handler = s.mux
	// start server
	return s.srv.ListenAndServe()
}
