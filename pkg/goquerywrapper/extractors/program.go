package extractors

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/exp/maps"
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
	"personal-feed/pkg/goquerywrapper/extractors/registry"
)

type GoQueryProgram struct {
	extractors []abstract.Extractor
}

func (p *GoQueryProgram) Do(in *goquery.Selection) (string, error) {
	return Exec(in, p.extractors)
}

func NewProgramFromString(in string) (*GoQueryProgram, error) {
	var list []map[string]string
	err := json.Unmarshal([]byte(in), &list)
	if err != nil {
		return nil, xerrors.Errorf("unable to unmarshal array, arr:%s, err:%w", in, err)
	}
	extractors := make([]abstract.Extractor, 0)
	for _, el := range list {
		instruction, _ := json.Marshal(el)

		currKeyArr := maps.Keys(el)
		if len(currKeyArr) != 1 {
			return nil, xerrors.Errorf("instruction must have only one k-v pair, instruction:%s", instruction)
		}
		currKey := currKeyArr[0]
		currVal := el[currKey]

		switch currKey {
		case "Query":
			newExtractor, err := registry.NewQuery(currVal)
			if err != nil {
				return nil, xerrors.Errorf("unable to create 'query' extractor from instruction, instruction:%s, err:%w", instruction, err)
			}
			extractors = append(extractors, newExtractor)
		case "Attr":
			newExtractor, err := registry.NewAttr(currVal)
			if err != nil {
				return nil, xerrors.Errorf("unable to create 'attr' extractor from instruction, instruction:%s, err:%w", instruction, err)
			}
			extractors = append(extractors, newExtractor)
		case "Regex":
			newExtractor, err := registry.NewRegex(currVal)
			if err != nil {
				return nil, xerrors.Errorf("unable to create 'regex' extractor from instruction, instruction:%s, err:%w", instruction, err)
			}
			extractors = append(extractors, newExtractor)
		case "Text":
			extractors = append(extractors, registry.NewText())
		}
	}

	if err := Validate(extractors); err != nil {
		return nil, xerrors.Errorf("program is invalid, program:%s, err:%w", in, err)
	}

	return &GoQueryProgram{
		extractors: extractors,
	}, nil
}

func NewProgramFromProgram(in Program) (*GoQueryProgram, error) {
	myStr, err := json.Marshal(in)
	if err != nil {
		return nil, xerrors.Errorf("unable to marshal program, err: %w", err)
	}
	return NewProgramFromString(string(myStr))
}
