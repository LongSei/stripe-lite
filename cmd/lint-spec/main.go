package main

import (
	"context"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	specPath := "api/openapi.yaml"
	if len(os.Args) > 1 {
		specPath = os.Args[1]
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading OpenAPI spec at %s: %v\n", specPath, err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := doc.Validate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "OpenAPI validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OpenAPI spec %s is valid (title: %q, version: %q, paths: %d)\n",
		specPath, doc.Info.Title, doc.Info.Version, doc.Paths.Len())
}
