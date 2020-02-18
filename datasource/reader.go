package datasource

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
		_, _ = w.Write([]byte("failed to dump MessageReader: " + err.Error() + "\n"))
	}

	_, _ = w.Write(data)
	_, _ = w.Write([]byte("\n"))
}

func (r *MessageReader) ReadNext() (string, time.Duration, bool) {
	if len(r.data) == 0 {
		return "", 0, false
	}

	nextData := r.data[0]
	r.data = r.data[1:]

	delay := time.Duration(nextData.DelaySec) * time.Second

	var line string

	for _, sentence := range nextData.Sentences {

		if len(sentence.Tags) > 0 {
			tagsChecksum := checkSum(sentence.Tags, nextData.CorrectChecksum)
			line += fmt.Sprintf("\\%s*%x\\", sentence.Tags, tagsChecksum)
		}

		if len(sentence.Params) > 0 {
			paramsChecksum := checkSum(sentence.Params, nextData.CorrectChecksum)
			line += fmt.Sprintf("$%s*%x", sentence.Params, paramsChecksum)
		}

		if nextData.EOL {
			line = line + "\r\n"
		}
	}

	fmt.Printf("-- READ NEXT: [%v]\n", line)

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
