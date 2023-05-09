package registry

import (
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
	"regexp"
)

type Regex struct {
	abstract.Extractor
	Regex string
	re    *regexp.Regexp
}

var _ abstract.ExtractorStringToString = (*Regex)(nil)

func (e *Regex) Do(in string) *string { // result is optional: nil is absence of value
	subMatch := e.re.FindStringSubmatch(in)
	if len(subMatch) != 2 {
		return nil
	}
	return &subMatch[1]
}

func NewRegex(regex string) (*Regex, error) {
	if regex == "" {
		return nil, xerrors.New("regex shouldn't be empty. it's optional entity - if you don't need it - just don't use it")
	}
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, xerrors.Errorf("unable to compile regex, regex:%s, err:%s", regex, err)
	}
	if re.NumSubexp() != 1 {
		return nil, xerrors.Errorf("number of 'capturing groups' must be equal to 1, regex:%s", regex)
	}
	return &Regex{
		Regex: regex,
		re:    re,
	}, nil
}
