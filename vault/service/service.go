package service

import (
	"context"
	"github.com/dark-enstein/vault/internal/vlog"
	"net/http"
	"time"
)

const (
	port = "8080"
)

type Service struct {
	sc  *StartConfig
	srv *http.Server
	mux *http.ServeMux
	log *vlog.Logger
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
	log.Logger().Debug().Msgf("settings:\n\taddress: %v\n\tread timeout: %v\n\t write timeout\n", ":"+port, readTimeout, writeTimeout)
	return &Service{sc: &StartConfig{port: port}, srv: hs, mux: http.NewServeMux(), log: log}
}

type StartConfig struct {
	port string
}

func (s *Service) Port() string {
	return s.sc.port
}

func (s *Service) LoadHandlers() {
	for k, v := range *NewVaultHandler() {
		s.mux.HandleFunc(k, v)
	}
}

func (s *Service) Run() error {
	// load handlers into mux
	s.LoadHandlers()
	// set mux into server
	s.srv.Handler = s.mux
	// start server
	return s.srv.ListenAndServe()
}
