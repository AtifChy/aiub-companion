//go:build ignore

package main

import (
	"aiub-companion/internal/settings"
	"encoding/json"
	"log"
	"os"

	"github.com/invopop/jsonschema"
)

func main() {
	schema := jsonschema.Reflect(&settings.Settings{})

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal schema: %v", err)
	}

	if err := os.WriteFile("schema.json", data, 0o644); err != nil {
		log.Fatalf("failed to write schema: %v", err)
	}
}
