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

	f := internal.StartFile("repository")
	err := internal.GenerateMapping(f, m, "repository")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = f.Render(&buf)
	require.NoError(t, err)

	fixtureOut, err := os.ReadFile("./fixtures/repository/mappings_mytype_gen.go")
	require.NoError(t, err, "error reading fixture file: %v")

	assert.Equal(t, string(fixtureOut), buf.String())
}

func TestSplitOutputFilename(t *testing.T) {
	m := myTypeStructMapping()

	outputFilename := internal.SplitOutputFilename(m, "mappings.go")
	assert.Equal(t, "mappings_mytype_gen.go", outputFilename)
}

func TestCombinedOutputFilename(t *testing.T) {
	outputFilename := internal.CombinedOutputFilename("mappings.go")
	assert.Equal(t, "mappings_gen.go", outputFilename)
}

func TestGenerateSamePackage(t *testing.T) {
	m := myTypeStructMapping()

	f := internal.StartFile("fixtures")

	err := internal.GenerateMapping(f, m, "fixtures")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = f.Render(&buf)
	require.NoError(t, err)

	fixtureOut, err := os.ReadFile("./fixtures/fixture_mytype_gen.go")
	require.NoError(t, err, "error reading fixture file: %v")

	assert.Equal(t, string(fixtureOut), buf.String())
}
