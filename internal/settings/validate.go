package settings

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema.json
var schemaJSON []byte

var compiledSchema *jsonschema.Schema

func init() {
	var schema any
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		panic(fmt.Errorf("failed to unmarshal schema: %w", err))
	}

	// Use a unique URI for the schema resource. This URI is only used internally by the jsonschema compiler.
	const schemaURI = "urn:aiub-companion:settings"

	c := jsonschema.NewCompiler()
	if err := c.AddResource(schemaURI, schema); err != nil {
		panic(fmt.Errorf("failed to add schema resource: %w", err))
	}

	var err error
	compiledSchema, err = c.Compile(schemaURI)
	if err != nil {
		panic(fmt.Errorf("failed to compile schema: %w", err))
	}
}

func validate(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if err := compiledSchema.Validate(v); err != nil {
		if validationErr, ok := errors.AsType[*jsonschema.ValidationError](err); ok {
			return fmt.Errorf("invalid settings: %v", validationErr.DetailedOutput())
		}
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
