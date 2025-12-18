package main

import (
	"os"
	"os/signal"
	"syscall"

	"stocky-backend/config"
	"stocky-backend/features/corporate_action"
	"stocky-backend/features/reward"
	"stocky-backend/features/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	config.InitLogger()

	db, err := config.ConnectDatabase()
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release"
	}
	gin.SetMode(ginMode)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Stocky Backend Assignment API is running",
		})
	})

	api := router.Group("/api")
	{
		rewardService := reward.NewRewardService(db)
		rewardHandler := reward.NewRewardHandler(rewardService)
		reward.RegisterRoutes(api, rewardHandler)

		corporateActionService := corporate_action.NewCorporateActionService(db)
		corporateActionHandler := corporate_action.NewCorporateActionHandler(corporateActionService)
		corporate_action.RegisterRoutes(api, corporateActionHandler)

		userService := user.NewUserService(db)
		userHandler := user.NewUserHandler(userService)
		user.RegisterRoutes(api, userHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		logrus.Infof("Server starting on port %s", port)
		if err := router.Run(":" + port); err != nil {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")
}
