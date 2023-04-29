package goquerywrapper

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

func URLToDoc(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(buf.String()))
	if err != nil {
		return nil, err
	}
	return doc, nil
}
