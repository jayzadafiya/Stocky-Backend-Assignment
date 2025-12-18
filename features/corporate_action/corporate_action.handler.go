package corporate_action

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CorporateActionHandler struct {
	service *CorporateActionService
}

func NewCorporateActionHandler(service *CorporateActionService) *CorporateActionHandler {
	return &CorporateActionHandler{service: service}
}

func (h *CorporateActionHandler) CreateCorporateAction(c *gin.Context) {
	var req CreateCorporateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ActionType == ActionStockSplit && req.SplitRatio <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "split_ratio is required and must be greater than 0 for stock split"})
		return
	}

	if req.ActionType == ActionMerger {
		if req.MergerToSymbol == "" || req.MergerRatio <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "merger_to_symbol and merger_ratio are required for merger"})
			return
		}
	}

	action, err := h.service.CreateCorporateAction(req)
	if err != nil {
		logrus.Errorf("Error creating corporate action: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Corporate action created successfully",
		"data":    action,
	})
}

func (h *CorporateActionHandler) ProcessCorporateAction(c *gin.Context) {
	actionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
		return
	}

	if err := h.service.ProcessCorporateAction(actionID); err != nil {
		logrus.Errorf("Error processing corporate action: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Corporate action processed successfully"})
}

func (h *CorporateActionHandler) GetAllCorporateActions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	response, err := h.service.GetAllCorporateActions(page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting corporate actions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve corporate actions"})
		return
	}

	c.JSON(http.StatusOK, response)
}
