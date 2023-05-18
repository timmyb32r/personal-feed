package server

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed dist
var distFiles embed.FS

type staticHandler struct {
}

func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathInDist := "dist" + r.RequestURI
	buf, err := distFiles.ReadFile(pathInDist)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("staticHandler::ServeHTTP::error1::%s", err.Error())))
		w.WriteHeader(404)
		return
	}
	_, _ = w.Write(buf)
}

func newStaticHandler() *staticHandler {
	return &staticHandler{}
}
