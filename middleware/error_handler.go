package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Detail  string `json:"detail,omitempty"`
	Code    int    `json:"code"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			logrus.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"error":  err.Error(),
			}).Error("Request error")

			if appErr, ok := err.Err.(*AppError); ok {
				c.JSON(appErr.Code, ErrorResponse{
					Success: false,
					Error:   appErr.Message,
					Detail:  appErr.Detail,
					Code:    appErr.Code,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "Internal server error",
				Detail:  err.Error(),
				Code:    http.StatusInternalServerError,
			})
		}
	}
}

func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"panic":  r,
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Success: false,
					Error:   "Internal server error",
					Detail:  "An unexpected error occurred",
					Code:    http.StatusInternalServerError,
				})

				c.Abort()
			}
		}()
		c.Next()
	}
}

func BadRequestError(message string, detail string) *AppError {
	return NewAppError(http.StatusBadRequest, message, detail)
}

func NotFoundError(message string, detail string) *AppError {
	return NewAppError(http.StatusNotFound, message, detail)
}

func InternalServerError(message string, detail string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, detail)
}

func UnauthorizedError(message string, detail string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, detail)
}

func ForbiddenError(message string, detail string) *AppError {
	return NewAppError(http.StatusForbidden, message, detail)
}

func ConflictError(message string, detail string) *AppError {
	return NewAppError(http.StatusConflict, message, detail)
}

func UnprocessableEntityError(message string, detail string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, message, detail)
}
