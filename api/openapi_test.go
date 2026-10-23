package api_test

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestOpenAPISpecValidation(t *testing.T) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromFile("openapi.yaml")
	if err != nil {
		t.Fatalf("failed to load openapi.yaml: %v", err)
	}

	ctx := context.Background()
	if err := doc.Validate(ctx); err != nil {
		t.Fatalf("openapi.yaml validation failed: %v", err)
	}

	if doc.Info.Title == "" {
		t.Error("expected non-empty info.title")
	}

	// Verify required tags
	expectedTags := map[string]bool{
		"System":                     false,
		"Authentication & Identity": false,
		"Payments & Core Engine":     false,
		"Webhooks":                   false,
	}

	for _, tag := range doc.Tags {
		if _, ok := expectedTags[tag.Name]; ok {
			expectedTags[tag.Name] = true
		}
	}

	for tag, found := range expectedTags {
		if !found {
			t.Errorf("expected tag %q not found in openapi.yaml", tag)
		}
	}
}
