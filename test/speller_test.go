package speller

import (
	"bytes"
	"encoding/json"
	"go-homework/internal/checker"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpeller(t *testing.T) {
	ch := checker.New()

	data := readTestData("data/data.json")
	want := readTestData("data/want.json")

	tests := []struct {
		texts     []string
		corrected []string
	}{
		{
			texts:     data,
			corrected: want,
		},
	}

	for _, tt := range tests {
		corrected, _ := ch.FixSpellsInBatchMode(tt.texts)
		assert.Equal(t, corrected, tt.corrected)
	}
}

func readTestData(fileName string) []string {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)

	var data map[string][]string

	dec := json.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	return data["texts"]
}
