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
	loopBack bool
	data     []model.Message
	idxNext  int
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
		loopBack: data.Loopback,
		data:     data.Data,
		idxNext:      0,
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

	if r.idxNext >= len(r.data) {
		if r.loopBack {
			r.idxNext = 0
		} else
		{
			return "", 0, false
		}
	}

		nextData := r.data[r.idxNext]
	delay := time.Duration(nextData.DelaySecMs) * time.Millisecond

	var line string

	for _, sentence := range nextData.Sentences {

		var tagLine string
		if sentence.Tags.AddTime {
			tagLine = fmt.Sprintf("c:%d", time.Now().Unix())
		}

		if len(sentence.Tags.Data) > 0 {
			tagLine += "," + sentence.Tags.Data
		}

		if len(tagLine) > 0 {
			tagsChecksum := checkSum(tagLine, nextData.CorrectChecksum)
			line += fmt.Sprintf("\\%s*%x\\", tagLine, tagsChecksum)
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

	r.idxNext++

	return line, delay, true
}

func checkSum(data string, correctChecksum bool) byte {
	if !correctChecksum {
		return 0
	}

	var sum byte
	for _, b := range []byte(data) {
		sum = sum ^ b
	}

	return sum
}
