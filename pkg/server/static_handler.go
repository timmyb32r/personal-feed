package server

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
)

//go:embed dist
var distFiles embed.FS

type staticHandler struct {
}

func (h staticHandler) ServeHTTPURI(w http.ResponseWriter, uri string) {
	pathInDist := "dist" + uri
	buf, err := distFiles.ReadFile(pathInDist)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("staticHandler::ServeHTTP::error1::%s", err.Error())))
		w.WriteHeader(404)
		return
	}
	if strings.HasSuffix(pathInDist, ".js") {
		w.Header().Add("Content-Type", "text/javascript") // https://stackoverflow.com/a/9664327
	}
	_, _ = w.Write(buf)
}

func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPURI(w, r.RequestURI)
}

func newStaticHandler() *staticHandler {
	return &staticHandler{}
}
