package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"personal-feed/pkg/config"
	"personal-feed/pkg/repo"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type HTTPServer struct {
	config      *config.Config
	logger      *logrus.Logger
	httpServer  http.Server
	shutdownReq chan bool
	once        sync.Once
}

// DEBUG

type sourceIDHandlerFunc func(http.ResponseWriter, *http.Request)
type sourceIDHandler struct {
	sourceIDHandlerFuncField sourceIDHandlerFunc
}

func (q *sourceIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q.sourceIDHandlerFuncField(w, r)
}

func newSourceIDHandler(sourceIDHandlerFuncIn sourceIDHandlerFunc) *sourceIDHandler {
	return &sourceIDHandler{
		sourceIDHandlerFuncField: sourceIDHandlerFuncIn,
	}
}

// DEBUG

func NewHTTPServer(config *config.Config, logger *logrus.Logger) *HTTPServer {
	s := &HTTPServer{
		config: config,
		logger: logger,
		httpServer: http.Server{
			Addr:         "0.0.0.0:80",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	router := mux.NewRouter()

	handler := newStaticHandler(logger)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTPURI(w, "/index.html")
	})
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTPURI(w, "/favicon.ico")
	})
	router.PathPrefix("/index.").Handler(handler)

	// DEBUG
	router.HandleFunc("/api/source_ids", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("[source_ids] got request on URI: %s", r.RequestURI)
		s.sourcesHandler(w, r)
	})
	router.PathPrefix("/api/source_id/").Handler(newSourceIDHandler(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("[source_id] got request on URI: %s", r.RequestURI)
		s.sourceIDHandler(w, r)
	}))
	// DEBUG

	s.httpServer.Handler = router

	return s
}

func (s *HTTPServer) Close() {
	s.shutdown()
}

func (s *HTTPServer) ListenAndServe() {
	s.httpServer.ListenAndServe()
}

func (s *HTTPServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}
}

func (s *HTTPServer) shutdown() {
	s.once.Do(func() {
		log.Printf("Stoping http server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			log.Printf("Shutdown request error: %v", err)
		}
	})
}

// DEBUG

func (s *HTTPServer) sourcesHandler(w http.ResponseWriter, r *http.Request) {
	repoClient, err := repo.NewRepo(r.Context(), s.config.Repo, s.logger)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error0::%s", err.Error())))
		return
	}

	tx, err := repoClient.NewTx(r.Context())
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::sourcesHandler::error1::%s", err.Error())))
		return
	}
	defer tx.Rollback(r.Context())

	sources, err := repoClient.ListSources(r.Context())
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::sourcesHandler::error2::%s", err.Error())))
		return
	}

	type idAndTitle struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	result := make([]idAndTitle, 0)

	for _, el := range sources {
		result = append(result, idAndTitle{
			ID:    strconv.Itoa(el.ID),
			Title: el.Description,
		})
	}
	resultArr, _ := json.Marshal(result)
	s.logger.Infof("[sources] made response for URI: %s, response: %s", r.RequestURI, string(resultArr))
	_, _ = w.Write(resultArr)
}

func (s *HTTPServer) sourceIDHandler(w http.ResponseWriter, r *http.Request) {
	repoClient, err := repo.NewRepo(r.Context(), s.config.Repo, s.logger)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error0::%s", err.Error())))
		return
	}

	tx, err := repoClient.NewTx(r.Context())
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::sourceIDHandler::error1::%s", err.Error())))
		return
	}
	defer tx.Rollback(r.Context())

	nodes, err := repoClient.TestExtractAllTreeNodes(tx, r.Context())
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::sourceIDHandler::error2::%s", err.Error())))
		return
	}

	type feedEvent struct {
		ID          string    `json:"id"`
		At          time.Time `json:"at"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}
	type feed struct {
		Events []feedEvent `json:"events"`
	}
	result := feed{
		Events: make([]feedEvent, 0),
	}

	for _, el := range nodes {
		result.Events = append(result.Events, feedEvent{
			ID:          el.CurrentFullKey,
			At:          el.BusinessTime,
			Title:       "my-title-stub",
			Description: "", // temporary is empty - for debugging
		})
	}
	resultArr, _ := json.Marshal(result)
	s.logger.Infof("[source_id] made response for URI: %s, response: %s", r.RequestURI, string(resultArr))
	_, _ = w.Write(resultArr)
}

// DEBUG
