package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *UserHandler) {
	users := router.Group("/users")
	{
		users.GET("", handler.GetAllUsers)
		users.GET("/:id", handler.GetUserByID)
	}

	router.GET("/today-stocks/:userId", handler.GetTodayStockRewards)
	router.GET("/historical-inr/:userId", handler.GetHistoricalINRValues)
	router.GET("/stats/:userId", handler.GetUserStats)
	router.GET("/portfolio/:userId", handler.GetUserPortfolio)
}
