package extractors

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/xerrors"
	"personal-feed/pkg/goquerywrapper/extractors/abstract"
)

// Validate - checks one rule - 'string' can't be converted into 'selection'
func Validate(extractors []abstract.Extractor) error {
	currState := stateSelection
	step := 0
	for _, currExtractor := range extractors {
		var newState state
		switch currExtractor.(type) {
		case abstract.ExtractorSelectionToString:
			{
				newState = stateString
				if currState == stateString {
					return xerrors.Errorf("error on step %d: 'string' state can't be converted into 'selected'", step)
				}
			}
		case abstract.ExtractorSelectionToSelectionMaybeError:
			{
				newState = stateSelection
				if currState == stateString {
					return xerrors.Errorf("error on step %d: 'string' state can't be converted into 'selected'", step)
				}
			}
		case abstract.ExtractorStringToString:
			{
				newState = stateString
			}
		}
		currState = newState
		step++
	}
	return nil
}

func Exec(in *goquery.Selection, extractors []abstract.Extractor) (string, error) {
	if err := Validate(extractors); err != nil {
		return "", xerrors.Errorf("extractors is invalid, err: %w", err)
	}

	currState := newStateObj(in)
	step := 0
	for _, currExtractor := range extractors {
		switch extractor := currExtractor.(type) {
		case abstract.ExtractorSelectionToString:
			{
				selection, err := currState.GetSelection()
				if err != nil {
					return "", xerrors.Errorf("currState.GetSelection() returned an error, err: %w", err)
				}
				str := extractor.Do(selection)
				if str == nil {
					return "", xerrors.Errorf("error on step %d: string not found", step)
				}
				currState.SetString(*str)
			}
		case abstract.ExtractorSelectionToSelectionMaybeError:
			{
				selection, err := currState.GetSelection()
				if err != nil {
					return "", xerrors.Errorf("currState.GetSelection() returned an error, err: %w", err)
				}
				newSelection, err := extractor.Do(selection)
				if err != nil {
					return "", xerrors.Errorf("extractor returned an error, err: %w", err)
				}
				if newSelection == nil {
					return "", xerrors.Errorf("error on step %d: selection not found", step)
				}
				currState.SetSelection(newSelection)
			}
		case abstract.ExtractorStringToString:
			{
				str, err := currState.GetString()
				if err != nil {
					return "", xerrors.Errorf("currState.GetString() returned an error, err: %w", err)
				}
				newStr := extractor.Do(str)
				if newStr == nil {
					return "", xerrors.Errorf("error on step %d: string not found", step)
				}
				currState.SetString(*newStr)
			}
		}
		step++
	}
	result, err := currState.GetString()
	if err != nil {
		return "", xerrors.Errorf("currState.GetString() returned an error, err: %w", err)
	}
	return result, nil
}
