package input

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"main/model"
)

type MessageReader struct {
	data []model.Message
}

func New(fname string) (*MessageReader, error) {

	fileBody, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	var data model.MessageChain

	if err := json.Unmarshal(fileBody, &data); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file %s", fname)
	}

	return &MessageReader{
		data: data.Data,
	}, nil

}

func (r *MessageReader) ReadNext() (string, time.Duration, bool) {
	if len(r.data) == 0 {
		return "", 0, false
	}

	var line string

	if len(r.data[0].Tags) > 0{
		tagsChecksum := checkSum(r.data[0].Tags, r.data[0].CorrectChecksum)
		line = fmt.Sprintf("\\%s*%x\\", r.data[0].Tags, tagsChecksum)
	}

	if len(r.data[0].Params) > 0{
		paramsChecksum := checkSum(r.data[0].Params, r.data[0].CorrectChecksum)
		line += fmt.Sprintf("$%s*%x", r.data[0].Params, paramsChecksum)
	}

	line = line + "\r\n"

	delay := time.Duration(r.data[0].DelaySec) * time.Second

	r.data = r.data[1:]

	return line, delay, true
}

func checkSum(data string, correctChecksum bool) rune {
	if !correctChecksum {
		return 12
	}

	var sum rune
	for _, b := range data {
		sum = sum ^ b
	}

	return sum
}