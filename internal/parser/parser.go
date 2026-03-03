package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
	"github.com/dablon/runflow/internal/models"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(content string) (*models.ParsedRunbook, error) {
	var rb models.ParsedRunbook
	err := yaml.Unmarshal([]byte(content), &rb)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &rb, nil
}

func (p *Parser) Validate(content string) error {
	_, err := p.Parse(content)
	return err
}
