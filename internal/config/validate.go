package config

import (
	_ "embed"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema.json
var schemaJSON []byte

// schemaURI is a unique URI for the schema resource.
const schemaURI = "urn:aiub-companion:config"

// compiledSchema is the compiled JSON schema used for validation.
var compiledSchema *jsonschema.Schema

func init() {
	var schema any
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		panic(fmt.Errorf("unmarshal schema: %w", err))
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(schemaURI, schema); err != nil {
		panic(fmt.Errorf("add schema resource: %w", err))
	}

	var err error
	compiledSchema, err = c.Compile(schemaURI)
	if err != nil {
		panic(fmt.Errorf("compile schema: %w", err))
	}
}

func validate(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		if syntaxErr, ok := errors.AsType[*jsontext.SyntacticError](err); ok {
			return fmt.Errorf("invalid JSON syntax at offset %d: %w", syntaxErr.ByteOffset, err)
		}
		if semanticErr, ok := errors.AsType[*json.SemanticError](err); ok {
			return fmt.Errorf("invalid JSON for %s (expected %s): %w", semanticErr.JSONPointer, semanticErr.GoType, err)
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := compiledSchema.Validate(v); err != nil {
		if validationErr, ok := errors.AsType[*jsonschema.ValidationError](err); ok {
			b, _ := json.Marshal(validationErr)
			return fmt.Errorf("invalid config: %s", string(b))
		}
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
