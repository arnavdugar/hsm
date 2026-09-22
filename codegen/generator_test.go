package main_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/arnavdugar/hsm/codegen/golang"
	"github.com/arnavdugar/hsm/codegen/mermaid"
	"github.com/arnavdugar/hsm/codegen/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The example tests execute checked-in generated code. Keep those artifacts in
// sync with the generator so regressions cannot hide behind stale output.
func TestGeneratedExamples(t *testing.T) {
	paths, err := filepath.Glob("../example/*/machine.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, paths)
	for _, path := range paths {
		directory := filepath.Dir(path)
		t.Run(filepath.Base(directory), func(t *testing.T) {
			input, err := os.ReadFile(path)
			require.NoError(t, err)
			machine, err := parser.Parse(bytes.NewReader(input))
			require.NoError(t, err)

			generated, err := golang.Create(machine).Render()
			require.NoError(t, err)
			assertGenerated(t, filepath.Join(directory, "machine.go"), generated, "\npackage ")

			if machine.Codegen.Mermaid.Enabled {
				generated, err := mermaid.Create(machine).Render()
				require.NoError(t, err)
				assertGenerated(t, filepath.Join(directory, machine.Codegen.Mermaid.Filename), generated, "# State Machine\n")
			}
		})
	}
}

func assertGenerated(t *testing.T, path string, generated []byte, marker string) {
	t.Helper()
	expected, err := os.ReadFile(path)
	require.NoError(t, err)
	// Regeneration commands contain os.Args, which differ in a test process.
	expectedStart := bytes.Index(expected, []byte(marker))
	generatedStart := bytes.Index(generated, []byte(marker))
	require.NotEqual(t, -1, expectedStart)
	require.NotEqual(t, -1, generatedStart)
	assert.Equal(t, string(expected[expectedStart:]), string(generated[generatedStart:]), "%s is stale; run go generate ./example/...", path)
}
