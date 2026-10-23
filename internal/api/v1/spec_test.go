package v1_test

import (
	"testing"

	v1 "github.com/thep1ckaxe91/stripe-lite/internal/api/v1"
)

func TestRawSpecData(t *testing.T) {
	data, err := v1.RawSpecData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("data preview: %s", string(data[:200]))
}
