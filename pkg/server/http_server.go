package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"personal-feed/pkg/db/pg"
	"strings"
	"sync"
	"syscall"
	"time"
)

type HTTPServer struct {
	config      *Config
	httpServer  http.Server
	shutdownReq chan bool
	once        sync.Once
}

func NewHTTPServer(config *Config) *HTTPServer {
	s := &HTTPServer{
		config: config,
		httpServer: http.Server{
			Addr:         "0.0.0.0:80",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	router := mux.NewRouter()

	router.HandleFunc("/", s.RootHandler)

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

func (s *HTTPServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	pgConfig := pg.NewConfig(s.config.DBUser, s.config.DBPassword, s.config.DBHost, s.config.DBPort, s.config.DBName, true)
	pgClient, err := pg.NewPgClient(pgConfig)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error0::%s", err.Error())))
		return
	}

	tx, err := pgClient.NewTx()
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error1::%s", err.Error())))
		return
	}

	nodes, err := pgClient.TestExtractAllTreeNodes(tx)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error2::%s", err.Error())))
		return
	}

	buf := []string{"RESULT:"}
	for _, el := range nodes {
		elArr, _ := json.Marshal(el)
		buf = append(buf, string(elArr))
	}
	result := strings.Join(buf, "\n")
	_, _ = w.Write([]byte(result))
}
