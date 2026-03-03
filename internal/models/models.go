package models

import (
	"time"
)

type Runbook struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Content     string    `json:"content" gorm:"type:text"`
	Environment string    `json:"environment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Execution struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RunbookID   uint      `json:"runbook_id"`
	Status      string    `json:"status"` // pending, running, success, failure
	Output      string    `json:"output" gorm:"type:text"`
	Error       string    `json:"error"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type Provider struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"uniqueIndex;not null"`
	Type       string    `json:"type"` // aws, azure, gcp, kubernetes, etc.
	Config     string    `json:"config" gorm:"type:text"` // encrypted JSON
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ParsedRunbook struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Environment string            `yaml:"environment"`
	Variables   map[string]string `yaml:"variables"`
	Steps       []Step            `yaml:"steps"`
	OnFailure  []Step            `yaml:"on_failure"`
	OnSuccess  []Step           `yaml:"on_success"`
}

type Step struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Provider    string            `yaml:"provider"`
	Env         map[string]string `yaml:"env"`
	Timeout     int               `yaml:"timeout"`
	Retries     int               `yaml:"retries"`
	RetryDelay  int               `yaml:"retry_delay"`
}

type ExecutionResult struct {
	Success bool
	Output  string
	Error   string
}
