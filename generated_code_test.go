package main

import (
	"bytes"
	"go/format"
	"os"
	"testing"
)

func TestGeneratedModelsAreFormatted(t *testing.T) {
	source, err := os.ReadFile("azurecaf/models_generated.go")
	if err != nil {
		t.Fatalf("read generated models: %v", err)
	}

	formatted, err := format.Source(source)
	if err != nil {
		t.Fatalf("parse generated models: %v", err)
	}
	if !bytes.Equal(source, formatted) {
		t.Fatal("azurecaf/models_generated.go is not gofmt-clean; run go generate")
	}
}
