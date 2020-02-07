package model


type MessageChain struct {
	Data []Message `json:"data"`
}

type Message struct {
	DelaySec        int        `json:"delay"`
	CorrectChecksum bool       `json:"correctChecksum"`
	Sentences       []Sentence `json:"sentences"`
}

type Sentence struct {
	Tags   string `json:"tags"`
	Params string `json:"params"`
}
