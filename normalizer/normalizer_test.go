package normalizer

import (
	// "fmt"
	"encoding/csv"
	"os"
	"log"
	"testing"
)

var TEST_FILE = "test_data.csv"

func TestNormalizer(t *testing.T) {
	data := openFile()
	for _, i := range data {
		temp := NormalizeNumber(i[0])
		if(temp != i[1]) {
			t.Errorf("Nomralization failed : \"%s\" != \"%s\"\n", temp, i[1])
		}
	}

}

func openFile() [][]string {
	f, err := os.Open(TEST_FILE)
	HandleErr(err)

	defer f.Close()
	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	HandleErr(err)

	return data
}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}