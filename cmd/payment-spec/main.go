// Command payment-spec builds the Payments integration entrypoint without editing the shared root.
package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("api/openapi.yaml")
	if err != nil {
		panic(err)
	}
	var doc map[string]any
	if err = yaml.Unmarshal(data, &doc); err != nil {
		panic(err)
	}
	doc["servers"] = []any{map[string]any{"url": "/v1", "description": "This server uses the v1 API prefix."}}
	paths := doc["paths"].(map[string]any)
	for path, ref := range map[string]string{
		"/payment-intents":              "./payments/payment_intents.yaml",
		"/payment-intents/{id}":         "./payments/payment_intents.yaml#/x-retrieve",
		"/payment-intents/{id}/confirm": "./payments/payment_intents.yaml#/x-confirm",
		"/payment-intents/{id}/capture": "./payments/payment_intents.yaml#/x-capture",
		"/payment-intents/{id}/cancel":  "./payments/payment_intents.yaml#/x-cancel",
		"/refunds":                      "./payments/refunds.yaml",
		"/refunds/{id}":                 "./payments/refunds.yaml#/x-retrieve",
	} {
		paths[path] = map[string]any{"$ref": ref}
	}
	// The generated entrypoint is written in dist; source modules stay in api.
	var refs func(any)
	refs = func(value any) {
		switch v := value.(type) {
		case map[string]any:
			for k, x := range v {
				if k == "$ref" {
					if ref, ok := x.(string); ok && strings.HasPrefix(ref, "./") {
						v[k] = "../api/" + strings.TrimPrefix(ref, "./")
					}
				} else {
					refs(x)
				}
			}
		case []any:
			for _, x := range v {
				refs(x)
			}
		}
	}
	refs(doc)
	out, err := yaml.Marshal(doc)
	if err != nil {
		panic(err)
	}
	if _, err = os.Stdout.Write(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
