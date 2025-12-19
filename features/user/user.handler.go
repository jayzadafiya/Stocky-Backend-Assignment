package user

import (
	"net/http"
	"strconv"

	"stocky-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, err := h.service.GetAllUsers(page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting users: %v", err)
		c.Error(middleware.InternalServerError("Failed to retrieve users", err.Error()))
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(middleware.BadRequestError("Invalid user ID", err.Error()))
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		logrus.Errorf("Error getting user: %v", err)
		c.Error(middleware.NotFoundError("User not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) GetTodayStockRewards(c *gin.Context) {
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

	rewards, err := h.service.GetTodayStockRewards(userID, page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting today's stock rewards: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve today's stock rewards"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *UserHandler) GetHistoricalINRValues(c *gin.Context) {
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

	historicalValues, err := h.service.GetHistoricalINRValues(userID, page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting historical INR values: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve historical INR values"})
		return
	}

	c.JSON(http.StatusOK, historicalValues)
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.service.GetUserStats(userID)
	if err != nil {
		logrus.Errorf("Error getting user stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *UserHandler) GetUserPortfolio(c *gin.Context) {
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

	portfolio, err := h.service.GetUserPortfolio(userID, page, pageSize)
	if err != nil {
		logrus.Errorf("Error getting user portfolio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}
