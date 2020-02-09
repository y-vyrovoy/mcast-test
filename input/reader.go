package input

import (
	"encoding/json"
	"fmt"
	"io"
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

func (r *MessageReader) Dump(w io.Writer) {
	data, err := json.MarshalIndent(r.data, "", "  ")

	if err != nil {
		w.Write([]byte("failed to dump MessageReader: " + err.Error() + "\n"))
	}

	_, _ = w.Write(data)
	_, _ = w.Write([]byte("\n"))
}

func (r *MessageReader) ReadNext() (string, time.Duration, bool) {
	if len(r.data) == 0 {
		return "", 0, false
	}

	delay := time.Duration(r.data[0].DelaySec) * time.Second

	var line string

	for _, sentence := range r.data[0].Sentences{

		if len(sentence.Tags) > 0{
			tagsChecksum := checkSum(sentence.Tags, r.data[0].CorrectChecksum)
			line += fmt.Sprintf("\\%s*%x\\", sentence.Tags, tagsChecksum)
		}

		if len(sentence.Params) > 0{
			paramsChecksum := checkSum(sentence.Params, r.data[0].CorrectChecksum)
			line += fmt.Sprintf("$%s*%x", sentence.Params, paramsChecksum)
		}

		line = line + "\r\n"
	}

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