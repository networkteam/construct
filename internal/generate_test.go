package internal_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/construct/v2/internal"
)

func TestGenerate(t *testing.T) {
	m := myTypeStructMapping()

	var buf bytes.Buffer
	outputFilename, err := internal.Generate(m, "repository", "mappings.go", &buf)
	require.NoErrorf(t, err, "error generating code: %v")

	const expectedOutputFilename = "mappings_mytype_gen.go"
	require.Equal(t, expectedOutputFilename, outputFilename, "expected output filename to be %s, but got %s", expectedOutputFilename, outputFilename)

	fixtureOut, err := os.ReadFile("./fixtures/repository/" + expectedOutputFilename)
	require.NoError(t, err, "error reading fixture file: %v")

	assert.Equal(t, string(fixtureOut), buf.String())
}

func TestGenerateSamePackage(t *testing.T) {
	m := myTypeStructMapping()

	var buf bytes.Buffer
	outputFilename, err := internal.Generate(m, "fixtures", "fixture.go", &buf)
	require.NoErrorf(t, err, "error generating code: %v")

	const expectedOutputFilename = "fixture_mytype_gen.go"
	require.Equal(t, expectedOutputFilename, outputFilename, "expected output filename to be %s, but got %s", expectedOutputFilename, outputFilename)

	fixtureOut, err := os.ReadFile("./fixtures/other/" + expectedOutputFilename)
	require.NoError(t, err, "error reading fixture file: %v")

	assert.Equal(t, string(fixtureOut), buf.String())
}
