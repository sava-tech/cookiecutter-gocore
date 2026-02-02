package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// StandardResponse represents a standard API response
type StandardResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field string `json:"field" example:"field"`
	Error string `json:"error" example:"must be a valid field"`
}

// ErrorResponse returns a standardized error response
func ErrorResponse(err error) gin.H {
	return gin.H{
		"status": false,
		"error":  err.Error(),
	}
}

func ErrorResponseMessage(err string) gin.H {
	return gin.H{
		"status": false,
		"error":  err,
	}
}

// ValidationErrorResponse returns a standardized validation error response
func ValidationErrorResponse(errors []ValidationError) gin.H {
	return gin.H{
		"success": false,
		"errors":  errors,
	}
}

// SuccessResponse returns a standardized success response
func SuccessResponse(message string, data interface{}) gin.H {
	return gin.H{
		"status":  true,
		"message": message,
		"data":    data,
	}
}

// CreatedResponse returns a standardized created response
func CreatedResponse(message string, data interface{}) gin.H {
	return gin.H{
		"status":  true,
		"message": message,
		"data":    data,
	}
}

func ResponseError(message string, err error) gin.H {
	return gin.H{
		"status": false,
		"error":  err,
	}
}

// DeletedResponse returns a standardized deleted response
func DeletedResponse(message string) gin.H {
	return gin.H{
		"status":  true,
		"message": message,
	}
}

// Helper functions for common HTTP responses
func RespondWithError(ctx *gin.Context, code int, err error) {
	if code == http.StatusInternalServerError {
		//TODO: create a logger for handling logs so we can push it externally too
		log.Printf("Internal Server Error: %v", err)
		ctx.JSON(code, ErrorResponse(errors.New("something went wrong, please try again")))
		return

	}
	if code == http.StatusBadRequest {
		var errors []ValidationError

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationError := range validationErrors {
				errors = append(errors, ValidationError{
					Field: validationError.Field(),
					Error: validationError.Error(),
				})
			}
			ctx.JSON(code, ValidationErrorResponse(errors))
			return
		}

	}

	ctx.JSON(code, ErrorResponse(err))

}

func RespondWithSuccess(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, SuccessResponse(message, data))
}

func RespondWithCreated(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusCreated, CreatedResponse(message, data))
}

func RespondWithDeleted(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, DeletedResponse(message))
}
