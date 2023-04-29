package server

import (
	"github.com/sirupsen/logrus"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/db/pg"
	"personal-feed/pkg/engine"
)

type Server struct {
	config     *Config
	httpServer *HTTPServer
}

func NewServer(config *Config) *Server {
	httpServer := NewHTTPServer(config)

	go func() {
		httpServer.ListenAndServe()
	}()

	return &Server{
		config:     config,
		httpServer: httpServer,
	}
}

func (s *Server) Close() {
	s.httpServer.Close()
}

func (s *Server) RunIteration(logger *logrus.Logger) error {
	pgConfig := pg.NewConfig(s.config.DBUser, s.config.DBPassword, s.config.DBHost, s.config.DBPort, s.config.DBName, true)
	pgClient, err := pg.NewPgClient(pgConfig)
	if err != nil {
		return err
	}

	tx, err := pgClient.NewTx()
	if err != nil {
		return err
	}
	sources, err := pgClient.ListSources(tx)
	if err != nil {
		return err
	}

	// TODO - add scheduler

	for _, source := range sources {
		currCrawler, err := crawlers.NewCrawler(source, logger)
		if err != nil {
			return err
		}
		currEngine := engine.NewEngine(source, currCrawler, pgClient)
		err = currEngine.RunOnce()
		if err != nil {
			return err
		}
	}
	return nil
}
