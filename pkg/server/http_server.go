package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"personal-feed/pkg/config"
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

	handler := newStaticHandler()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTPURI(w, "/index.html")
	})
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTPURI(w, "/favicon.ico")
	})
	router.PathPrefix("/index.").Handler(handler)

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

//func (s *HTTPServer) RootHandler(w http.ResponseWriter, r *http.Request) {
//	repoClient, err := repo.NewRepo(r.Context(), s.config.Repo, s.logger)
//	if err != nil {
//		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error0::%s", err.Error())))
//		return
//	}
//
//	tx, err := repoClient.NewTx(r.Context())
//	if err != nil {
//		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error1::%s", err.Error())))
//		return
//	}
//	defer tx.Rollback(r.Context())
//
//	nodes, err := repoClient.TestExtractAllTreeNodes(tx, r.Context())
//	if err != nil {
//		_, _ = w.Write([]byte(fmt.Sprintf("HTTPServer::RootHandler::error2::%s", err.Error())))
//		return
//	}
//
//	buf := []string{"RESULT:"}
//	for _, el := range nodes {
//		elArr, _ := json.Marshal(el)
//		buf = append(buf, string(elArr))
//	}
//	result := strings.Join(buf, "\n")
//	_, _ = w.Write([]byte(result))
//}

func (s *HTTPServer) DistHandler(w http.ResponseWriter, r *http.Request) {

}
