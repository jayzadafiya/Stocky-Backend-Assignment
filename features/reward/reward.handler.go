package reward

import (
	"net/http"
	"strconv"

	"stocky-backend/middleware"

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
		c.Error(middleware.BadRequestError("Invalid request body", err.Error()))
		return
	}

	reward, err := h.service.CreateReward(req)
	if err != nil {
		logrus.Errorf("Error creating reward: %v", err)
		c.Error(middleware.InternalServerError("Failed to create reward", err.Error()))
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
		c.Error(middleware.InternalServerError("Failed to retrieve rewards", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *RewardHandler) GetRewardsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.Error(middleware.BadRequestError("Invalid user ID", err.Error()))
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
		c.Error(middleware.InternalServerError("Failed to retrieve user rewards", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *RewardHandler) AdjustReward(c *gin.Context) {
	var req AdjustRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.BadRequestError("Invalid request body", err.Error()))
		return
	}

	adjustment, err := h.service.AdjustReward(req)
	if err != nil {
		logrus.Errorf("Error adjusting reward: %v", err)
		c.Error(middleware.InternalServerError("Failed to adjust reward", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reward adjusted successfully",
		"data":    adjustment,
	})
}
