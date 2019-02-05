package codegen

import (
	"testing"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
)

func TestTypeUnionAsInput(t *testing.T) {
	err := generate("inputunion", `testdata/unioninput.graphqls`)

	require.EqualError(t, err, "unable to build object definition: Query.addBookmark: cannot use Bookmarkable as argument b because UNION is not a valid input type")
}

func TestTypeInInput(t *testing.T) {
	err := generate("typeinput", `testdata/typeinput.graphqls`)

	require.EqualError(t, err, "unable to build input definition: BookmarkableInput.item: cannot use Item because OBJECT is not a valid input type")
}

func generate(name string, schemaFilename string) error {
	_, err := NewSchema(&config.Config{
		SchemaFilename: config.SchemaFilenames{schemaFilename},
		Exec:           config.PackageConfig{Filename: "gen/" + name + "/exec.go"},
		Model:          config.PackageConfig{Filename: "gen/" + name + "/model.go"},
	})

	return err
}