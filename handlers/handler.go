package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mrrancy/logAssignment/models"
	"mrrancy/logAssignment/store"
	"net/http"
)

type LogHandlerResponse struct {
	Message string              `json:"message"`
	Data    []models.LogPayload `json:"data"`
}

type Controller struct {
	Log   *zap.Logger
	Cache *store.MyMem
}

func NewController(logger *zap.Logger, cache *store.MyMem) *Controller {
	return &Controller{
		Log:   logger,
		Cache: cache,
	}
}

func (h *Controller) Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (h *Controller) LogHandler(c *gin.Context) {
	var payload models.LogPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("Failed to parse JSON payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	h.Cache.Put(payload)

	response := LogHandlerResponse{
		Message: "Log payload received successfully",
		Data:    h.Cache.GetAll(),
	}

	c.JSON(http.StatusOK, response)
}
