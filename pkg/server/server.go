package server

import (
	"github.com/sirupsen/logrus"
	"personal-feed/pkg/config"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/engine"
	"personal-feed/pkg/repo"
)

type Server struct {
	config     *config.Config
	httpServer *HTTPServer
}

func NewServer(config *config.Config) *Server {
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
	//pgConfig := pg2.NewConfig(s.config.User, s.config.Password, s.config.Host, s.config.Port, s.config.Name, true)
	currRepo, err := repo.NewRepo(s.config.Repo)
	if err != nil {
		return err
	}

	tx, err := currRepo.NewTx()
	if err != nil {
		return err
	}
	sources, err := currRepo.ListSources(tx)
	if err != nil {
		return err
	}

	// TODO - add scheduler

	for _, source := range sources {
		logger.Infof("RunIteration::start id:%d\n", source.ID)

		currCrawler, err := crawlers.NewCrawler(source, logger)
		if err != nil {
			return err
		}
		currEngine := engine.NewEngine(source, currCrawler, currRepo)
		err = currEngine.RunOnce()
		if err != nil {
			return err
		}

		logger.Infof("RunIteration::end id:%d\n", source.ID)
	}
	return nil
}
