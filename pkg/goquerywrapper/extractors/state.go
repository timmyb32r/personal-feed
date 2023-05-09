package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
)

type state string

const stateSelection = state("selection")
const stateString = state("string")

type stateObj struct {
	st             state
	stateSelection *goquery.Selection
	stateString    string
}

func (o *stateObj) SetSelection(s *goquery.Selection) {
	o.st = stateSelection
	o.stateSelection = s
	o.stateString = ""
}

func (o *stateObj) GetSelection() (*goquery.Selection, error) {
	if o.st == stateString {
		return nil, xerrors.New("unable to convert string into selection")
	}
	return o.stateSelection, nil
}

func (o *stateObj) SetString(s string) {
	o.st = stateString
	o.stateSelection = nil
	o.stateString = s
}

func (o *stateObj) GetString() (string, error) {
	if o.st == stateString {
		return o.stateString, nil
	}
	result, err := goquery.OuterHtml(o.stateSelection)
	if err != nil {
		return "", xerrors.Errorf("converter selection to string returned an error, err: %w", err)
	}
	return result, nil
}

func newStateObj(s *goquery.Selection) *stateObj {
	return &stateObj{
		st:             stateSelection,
		stateSelection: s,
		stateString:    "",
	}
}
