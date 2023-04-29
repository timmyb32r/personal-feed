package clients

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func ChannelIDByDoc(doc *goquery.Document) (string, error) {
	val, exists := doc.Find("body meta[itemprop=channelId]").Attr("content")
	if !exists {
		return "", xerrors.Errorf("unable to find channelID")
	}
	return val, nil
}

func ChannelIDByHtml(body string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}
	return ChannelIDByDoc(doc)
}

func ChannelIDByURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", err
	}
	return ChannelIDByHtml(buf.String())
}

//---

var reYoutubeChannel = regexp.MustCompile(`^https://www.youtube.com/[0-9a-zA-Z@]+$`)

func ValidateLinkToChannel(url string) bool {
	return reYoutubeChannel.Match([]byte(url))
}
