package logging

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/rs/zerolog/log"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schemas
var schemaFiles embed.FS

// NewValidationWriter compiles the JSONSchema a producers a Writer which validates the provided JSON documents against
// the schema.
// It discards the data after validation.
func NewValidationWriter() io.Writer {
	compiler := jsonschema.NewCompiler()
	//nolint:wrapcheck
	err := fs.WalkDir(schemaFiles, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Panic().Str("path", path).Err(err).Msg("Failed to walk schemas")
		}

		if !dirEntry.Type().IsRegular() || filepath.Ext(dirEntry.Name()) != ".json" {
			return nil
		}

		reader, err := schemaFiles.Open(path)
		if err != nil {
			return err
		}

		return compiler.AddResource(filepath.Clean(path), reader)
	})
	if err != nil {
		log.Panic().Err(err).Msg("Failed to add all jsonschema resources")
	}

	schema, err := compiler.Compile("schemas/kube_transition_metrics.schema.json")
	if err != nil {
		log.Panic().Err(err).Msg("Failed to compile jsonschema")
	}

	return &validationWriter{schema: schema}
}

type validationWriter struct {
	schema *jsonschema.Schema
}

func (w *validationWriter) Write(data []byte) (int, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var document any

	err := decoder.Decode(&document)
	if err != nil {
		err = fmt.Errorf("failed to decode document for schema validation: %w", err)
	} else {
		err = w.schema.Validate(document)
		if err != nil {
			err = fmt.Errorf("document is not validated by schema: %w", err)
		}
	}

	return len(data), err
}
