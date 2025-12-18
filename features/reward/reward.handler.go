package reward

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RewardHandler struct {
	service *RewardService
}

func NewRewardHandler(service *RewardService) *RewardHandler {
	return &RewardHandler{service: service}
}

func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reward, err := h.service.CreateReward(req)
	if err != nil {
		logrus.Errorf("Error creating reward: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create reward: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reward created successfully with ledger entries",
		"data":    reward,
	})
}

func (h *RewardHandler) GetAllRewards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rewards, err := h.service.GetAllRewards(page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting rewards: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rewards"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *RewardHandler) GetRewardsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rewards, err := h.service.GetRewardsByUserID(userID, page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting user rewards: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user rewards"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}
