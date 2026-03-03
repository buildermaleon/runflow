package executor_test

import (
	"testing"
	"github.com/dablon/runflow/internal/executor"
	"github.com/dablon/runflow/internal/models"
)

func TestExecutor_ExecuteStep(t *testing.T) {
	e := executor.New()
	
	step := models.Step{
		Name:    "echo",
		Command: "echo hello",
		Timeout: 30,
	}
	
	result := e.ExecuteStep(step, map[string]string{})
	
	if !result.Success {
		t.Errorf("ExecuteStep() failed: %s", result.Error)
	}
}

func TestExecutor_Variables(t *testing.T) {
	e := executor.New()
	
	step := models.Step{
		Name:    "echo-var",
		Command: "echo {{NAME}}",
		Timeout: 30,
	}
	
	result := e.ExecuteStep(step, map[string]string{"NAME": "world"})
	
	if !result.Success {
		t.Errorf("ExecuteStep() failed: %s", result.Error)
	}
}
