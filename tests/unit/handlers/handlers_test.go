package handlers_test

import (
	"testing"
	"github.com/dablon/runflow/internal/handlers"
	"github.com/dablon/runflow/internal/parser"
)

func TestNewHandler(t *testing.T) {
	p := parser.New()
	h := handlers.New(p)
	
	if h == nil {
		t.Fatal("Handler should not be nil")
	}
	
	if h.parser == nil {
		t.Error("Parser should be set")
	}
}

func TestHandler_RunbookCRUD(t *testing.T) {
	p := parser.New()
	h := handlers.New(p)
	
	// Would test CreateRunbook, ListRunbooks, etc.
	// This is a placeholder for full handler tests
}
