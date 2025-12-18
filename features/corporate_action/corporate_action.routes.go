package corporate_action

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *CorporateActionHandler) {
	corporateAction := router.Group("/corporate-action")
	{
		corporateAction.POST("", handler.CreateCorporateAction)
		corporateAction.POST("/:id/process", handler.ProcessCorporateAction)
		corporateAction.GET("", handler.GetAllCorporateActions)
	}
}
