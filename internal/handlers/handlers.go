package handlers

import (
	"net/http"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/dablon/runflow/internal/parser"
	"github.com/dablon/runflow/internal/models"
)

type Handler struct {
	parser     *parser.Parser
	runbooks   map[uint]*models.Runbook
	executions map[uint]*models.Execution
	providers  map[uint]*models.Provider
	nextID     uint
}

func New(p *parser.Parser) *Handler {
	return &Handler{
		parser:     p,
		runbooks:   make(map[uint]*models.Runbook),
		executions: make(map[uint]*models.Execution),
		providers:  make(map[uint]*models.Provider),
		nextID:     1,
	}
}

func (h *Handler) CreateRunbook(c *gin.Context) {
	var rb models.Runbook
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
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
	
	log.Printf("Updated runbook: %d", rb.ID)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse: " + err.Error()})
		return
	}
	
	exec := &models.Execution{
		ID:       h.nextID,
		RunbookID: rb.ID,
		Status:   "success",
	}
	h.nextID++
	h.executions[exec.ID] = exec
	
	log.Printf("Executed runbook: %s (ID: %d)", rb.Name, rb.ID)
	_ = parsed
	
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
		c.JSON(http.StatusOK, gin.H{"logs": exec.Output})
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
