package internal_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/construct/v2/internal"
)

func TestGenerate(t *testing.T) {
	m := myTypeStructMapping()
	var buf bytes.Buffer
	outputFilename, err := internal.Generate(m, "repository", "mappings.go", &buf)
	if err != nil {
		t.Fatalf("error generating code: %v", err)
	}

	const expectedOutputFilename = "mappings_mytype_gen.go"
	if outputFilename != expectedOutputFilename {
		t.Errorf("expected output filename to be %s, but got %s", expectedOutputFilename, outputFilename)
	}

	fixtureOut, err := os.ReadFile("./fixtures/repository/" + expectedOutputFilename)
	if err != nil {
		t.Fatalf("error reading fixture file: %v", err)
	}

	if buf.String() != string(fixtureOut) {
		assert.Equal(t, string(fixtureOut), buf.String())
	}
}
