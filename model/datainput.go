package model

type MessageChain struct {
	Data []Message `json:"data"`
}

type Message struct {
	DelaySec        int    `json:"delay"`
	Tags            string `json:"tags"`
	Params          string `json:"params"`
	CorrectChecksum bool   `json:"correctChecksum"`
}
