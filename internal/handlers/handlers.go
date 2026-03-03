package handlers

import (
	"net/http"
	"strconv"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dablon/runflow/internal/parser"
	"github.com/dablon/runflow/internal/executor"
	"github.com/dablon/runflow/internal/models"
)

type Handler struct {
	parser    *parser.Parser
	exec     *executor.Executor
	runbooks   map[uint]*models.Runbook
	executions map[uint]*models.Execution
	providers  map[uint]*models.Provider
	nextID     uint
}

func New(p *parser.Parser, e *executor.Executor) *Handler {
	return &Handler{
		parser:    p,
		exec:     e,
		runbooks:   make(map[uint]*models.Runbook),
		executions: make(map[uint]*models.Execution),
		providers:  make(map[uint]*models.Provider),
		nextID:     1,
	}
}

func (h *Handler) CreateRunbook(c *gin.Context) {
	var rb models.Runbook
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	
	if err := h.parser.Validate(rb.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YAML: " + err.Error()})
		return
	}
	
	rb.ID = h.nextID
	h.nextID++
	h.runbooks[rb.ID] = &rb
	
	log.Printf("Created runbook: %s (ID: %d)", rb.Name, rb.ID)
	c.JSON(http.StatusCreated, rb)
}

func (h *Handler) ListRunbooks(c *gin.Context) {
	var list []models.Runbook
	for _, rb := range h.runbooks {
		list = append(list, *rb)
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetRunbook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	if rb, ok := h.runbooks[uint(id)]; ok {
		c.JSON(http.StatusOK, rb)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Runbook not found"})
}

func (h *Handler) UpdateRunbook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	var rb models.Runbook
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	rb.ID = uint(id)
	h.runbooks[rb.ID] = &rb
	
	log.Printf("Updated runbook: %d", id)
	c.JSON(http.StatusOK, rb)
}

func (h *Handler) DeleteRunbook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	if _, ok := h.runbooks[uint(id)]; ok {
		delete(h.runbooks, uint(id))
		log.Printf("Deleted runbook: %d", id)
		c.JSON(http.StatusOK, gin.H{"deleted": true})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Runbook not found"})
}

func (h *Handler) ExecuteRunbook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	rb, ok := h.runbooks[uint(id)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Runbook not found"})
		return
	}
	
	parsed, err := h.parser.Parse(rb.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parse error: " + err.Error()})
		return
	}
	
	// Create execution record
	exec := &models.Execution{
		ID:        h.nextID,
		RunbookID: rb.ID,
		Status:    "running",
		StartedAt: time.Now(),
	}
	h.nextID++
	h.executions[exec.ID] = exec
	
	// Execute steps REAL (not mocked!)
	go func() {
		var output string
		success := true
		
		// Execute each step
		for _, step := range parsed.Steps {
			log.Printf("Executing step: %s", step.Name)
			
			// Set defaults
			if step.Timeout == 0 {
				step.Timeout = 300 // 5 min default
			}
			
			result := h.exec.ExecuteStep(step, parsed.Variables)
			
			output += "=== " + step.Name + " ===\n"
			output += result.Output + "\n"
			
			if !result.Success {
				output += "ERROR: " + result.Error + "\n"
				success = false
				
				// Run on_failure steps
				for _, failStep := range parsed.OnFailure {
					output += "=== ON FAILURE: " + failStep.Name + " ===\n"
					failResult := h.exec.ExecuteStep(failStep, parsed.Variables)
					output += failResult.Output + "\n"
				}
				break
			}
		}
		
		if success {
			// Run on_success steps
			for _, succStep := range parsed.OnSuccess {
				output += "=== ON SUCCESS: " + succStep.Name + " ===\n"
				succResult := h.exec.ExecuteStep(succStep, parsed.Variables)
				output += succResult.Output + "\n"
			}
		}
		
		// Update execution
		now := time.Now()
		exec.Status = "success"
		if !success {
			exec.Status = "failure"
		}
		exec.Output = output
		exec.FinishedAt = &now
		
		log.Printf("Execution %d completed: %s", exec.ID, exec.Status)
	}()
	
	c.JSON(http.StatusAccepted, exec)
}

func (h *Handler) GetExecution(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	if exec, ok := h.executions[uint(id)]; ok {
		c.JSON(http.StatusOK, exec)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
}

func (h *Handler) GetExecutionLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	if exec, ok := h.executions[uint(id)]; ok {
		c.JSON(http.StatusOK, gin.H{"status": exec.Status, "logs": exec.Output})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
}

func (h *Handler) CreateProvider(c *gin.Context) {
	var p models.Provider
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	p.ID = h.nextID
	h.nextID++
	h.providers[p.ID] = &p
	
	log.Printf("Created provider: %s", p.Name)
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) ListProviders(c *gin.Context) {
	var list []models.Provider
	for _, p := range h.providers {
		list = append(list, *p)
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	if _, ok := h.providers[uint(id)]; ok {
		delete(h.providers, uint(id))
		c.JSON(http.StatusOK, gin.H{"deleted": true})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
}
