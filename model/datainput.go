package model

type MessageChain struct {
	Loopback bool      `json:"loopback"`
	Data     []Message `json:"data"`
}

type Message struct {
	DelaySecMs      int        `json:"delay"`
	CorrectChecksum bool       `json:"correctChecksum"`
	Sentences       []Sentence `json:"sentences"`
	EOL             bool       `json:"eol"`
}

type Sentence struct {
	Tags   TagsDefinition `json:"tags"`
	Params string         `json:"params"`
}

type TagsDefinition struct {
	AddTime bool   `json:"addTime"`
	Data    string `json:"data"`
}
