package util

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"strings"
)

func GetURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", xerrors.Errorf("unable to get http page by url, err: %w", err)
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func HTMLToDoc(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func URLToDoc(url string) (*goquery.Document, error) {
	page, err := GetURL(url)
	if err != nil {
		return nil, xerrors.Errorf("unable to convert url to string, err: %w", err)
	}
	return HTMLToDoc(page)
}
