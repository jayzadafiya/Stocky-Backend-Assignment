package reward

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *RewardHandler) {
	
	rewards := router.Group("/reward")
	{
		rewards.POST("", handler.CreateReward)
		rewards.GET("", handler.GetAllRewards)
		rewards.GET("/user/:userId", handler.GetRewardsByUserID)
	}
}
