package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dablon/runflow/internal/models"
)

type Executor struct{}

func New() *Executor {
	return &Executor{}
}

func (e *Executor) ExecuteStep(step models.Step, vars map[string]string) *models.ExecutionResult {
	cmd := replaceVariables(step.Command, vars)
	
	env := os.Environ()
	for k, v := range step.Env {
		env = append(env, k+"="+replaceVariables(v, vars))
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(step.Timeout)*time.Second)
	defer cancel()
	
	cmdExec := exec.CommandContext(ctx, "bash", "-c", cmd)
	cmdExec.Env = env
	
	var stdout, stderr bytes.Buffer
	cmdExec.Stdout = &stdout
	cmdExec.Stderr = &stderr
	
	err := cmdExec.Run()
	
	output := stdout.String()
	if stderr.String() != "" {
		output += "\n" + stderr.String()
	}
	
	if err != nil {
		return &models.ExecutionResult{
			Success: false,
			Output:  output,
			Error:   err.Error(),
		}
	}
	
	return &models.ExecutionResult{
		Success: true,
		Output:  output,
	}
}

func replaceVariables(input string, vars map[string]string) string {
	result := input
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}

var _ = fmt.Println
