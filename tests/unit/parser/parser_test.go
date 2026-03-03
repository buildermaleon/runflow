package parser_test

import (
	"testing"
	"github.com/dablon/runflow/internal/parser"
)

func TestParser_Parse(t *testing.T) {
	p := parser.New()
	
	content := `
name: Test Runbook
version: "1.0"
variables:
  APP_NAME: testapp
steps:
  - name: Echo
    command: echo hello
`
	
	rb, err := p.Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	
	if rb.Name != "Test Runbook" {
		t.Errorf("Name = %s, want Test Runbook", rb.Name)
	}
}

func TestParser_Validate(t *testing.T) {
	p := parser.New()
	
	err := p.Validate("name: test")
	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestParser_InvalidYAML(t *testing.T) {
	p := parser.New()
	
	_, err := p.Parse("name: [unclosed")
	if err == nil {
		t.Error("Should error on invalid YAML")
	}
}
