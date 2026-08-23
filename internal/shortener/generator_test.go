package shortener

import (
	"strings"
	"testing"
)

func TestRandomCodeGenerator_Length(t *testing.T) {
	generator := RandomCodeGenerator{}
	code, err := generator.Generate(6)

	if err != nil {
		t.Fatalf("generate code: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected length 6, got %d", len(code))
	}
}

func TestRandomCodeGenerator_AllowedSymbols(t *testing.T) {
	generator := RandomCodeGenerator{}
	code, err := generator.Generate(6)

	if err != nil {
		t.Fatalf("generate code: %v", err)
	}

	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for _, symbol := range code {
		if !strings.ContainsRune(alphabet, symbol) {
			t.Fatalf("unexpected symbol %q", symbol)
		}
	}
}

func TestRandomCodeGenerator_IncorrectLength(t *testing.T) {
	generator := RandomCodeGenerator{}
	_, err := generator.Generate(0)

	if err == nil {
		t.Fatal("expected error")
	}

}
