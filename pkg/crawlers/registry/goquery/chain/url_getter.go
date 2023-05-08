package goquery

import "personal-feed/pkg/goquerywrapper"

type URLGetter interface {
	Get(url string) (string, error)
}

type DefaultURLGetter struct{}

func (g *DefaultURLGetter) Get(url string) (string, error) {
	return goquerywrapper.GetURL(url)
}
