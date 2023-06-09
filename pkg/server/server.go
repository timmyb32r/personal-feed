package server

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/config"
	"personal-feed/pkg/crawlers"
	"personal-feed/pkg/engine"
	"personal-feed/pkg/model"
	"personal-feed/pkg/operation"
	"personal-feed/pkg/repo"
	"time"
)

type Server struct {
	config     *config.Config
	logger     *logrus.Logger
	httpServer *HTTPServer
}

func NewServer(config *config.Config, logger *logrus.Logger) *Server {
	httpServer := NewHTTPServer(config, logger)

	go func() {
		httpServer.ListenAndServe()
	}()

	return &Server{
		config:     config,
		logger:     logger,
		httpServer: httpServer,
	}
}

func (s *Server) Close() {
	s.httpServer.Close()
}

func (s *Server) runIteration(ctx context.Context, currRepo repo.Repo, source *model.Source, lastRunTime *time.Time, currentTime time.Time) error {
	schedule, err := cron.ParseStandard(source.Schedule)
	if err != nil {
		return xerrors.Errorf("unable to parse cron expression %s, err: %w", source.Schedule, err)
	}

	nextTime := schedule.Next(*lastRunTime)
	if nextTime.After(currentTime) {
		return nil // it's not time for you next run, bro
	}

	s.logger.Infof("started to handle source %d by the schedule", source.ID)
	defer s.logger.Infof("finished to handle source %d by the schedule", source.ID)

	numNotMatchedNotifier := func(crawlerDescr string, expected *int, real int) {
		if expected == nil {
			return
		}
		if *expected != real {
			s.logger.Warnf("NumNotMatched, crawler: %s, expected: %d, real: %d", crawlerDescr, *expected, real)
		}
	}

	currCrawler, err := crawlers.NewCrawler(*source, s.logger)
	if err != nil {
		return xerrors.Errorf("unable to create new crawler, err: %w", err)
	}
	currEngine, err := engine.NewEngine(source, numNotMatchedNotifier, currCrawler, currRepo, s.logger)
	if err != nil {
		return xerrors.Errorf("unable to create engine, err: %w", err)
	}
	op := operation.Operation{
		OperationType: operation.OpTypeRegularUpdate,
	}
	errRunOnce := currEngine.RunOnce(ctx, op)
	if errRunOnce != nil { // we are not react on this error immediately - to not get to endless loop if some source lead to error
		s.logger.Warnf("engine.RunOnce() returned an error, err: %s", errRunOnce)
	}
	if err := currRepo.SetCronLastRunTime(ctx, currentTime); err != nil {
		return xerrors.Errorf("unable to set cron time, err: %w", err)
	}
	if errRunOnce != nil {
		return xerrors.Errorf("currEngine.RunOnce returned an error, err: %w", err)
	}
	return nil
}

func (s *Server) HandleAllSources(ctx context.Context) error {
	currRepo, err := repo.NewRepo(ctx, s.config.Repo, s.logger)
	if err != nil {
		return xerrors.Errorf("unable to create repo: %w", err)
	}

	sources, err := currRepo.ListSources(ctx)
	if err != nil {
		return err
	}

	lastRunTime, currentTime, err := currRepo.GetNextCronPeriod(ctx)
	if err != nil {
		return xerrors.Errorf("GetNextCronPeriod returned an error, err: %w", err)
	}

	if lastRunTime == nil || currentTime.Before(*lastRunTime) {
		if lastRunTime == nil {
			s.logger.Infof("initialize last cron time by current time")
		} else if currentTime.Before(*lastRunTime) {
			s.logger.Warnf("wrong last cron time, last_time(UTC) '%v', current_time(UTC) '%v'", lastRunTime, currentTime)
		}
		if err := currRepo.SetCronLastRunTime(ctx, currentTime); err != nil {
			return xerrors.Errorf("unable to set cron time, err: %w", err)
		}
		return nil
	}

	for _, source := range sources {
		err := s.runIteration(ctx, currRepo, &source, lastRunTime, currentTime)
		if err != nil {
			s.logger.Errorf("handling source %d returned an error, err: %s", source.ID, err)
		}
	}
	return nil
}
