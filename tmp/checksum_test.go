package tmp

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {

	nmeaSentence := "GPGGA,050214.56,1333.7380,N,05027.3577,E,2,07,01,+0051,M,-034,M,,"
	//nmeaSentence := "g:1-2-33,b:bbb"
	//nmeaSentence := "g:2-2-33,a:aaaa,s:s9182837465"
	//nmeaSentence := "GPGGA"

	var sum rune
	for _, b := range nmeaSentence {
		sum = sum ^ b
	}

	fmt.Printf("checksum of %s: 0x%x (%d)\n", nmeaSentence, sum, sum)
}

func Test_UnixTime(t *testing.T) {

	strTime := "1153612428"

	intTime, err := strconv.Atoi(strTime)
	if err != nil {
		fmt.Println("can't parse time")
		return
	}

	tm := time.Unix(int64(intTime), 0)

	fmt.Println(tm)
}

func Test_NoErrorf(t *testing.T) {
	var err error
	err = errors.New("error message")

	b := assert.NoErrorf(t, err, "message %d", 12)
	fmt.Println("b:", b)
}
