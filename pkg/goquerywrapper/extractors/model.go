package extractors

type Instruction struct {
	Query string `json:"Query,omitempty"`
	Attr  string `json:"Attr,omitempty"`
	Regex string `json:"Regex,omitempty"`
	Text  string `json:"Text,omitempty"`
}

type Program []Instruction
