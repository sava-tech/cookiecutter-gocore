package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response structure for consistent API responses
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
}

// Meta contains pagination or additional metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	TotalCount int64 `json:"total_count,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a success response with pagination metadata
func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a created response (201)
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Accepted sends an accepted response (202)
func Accepted(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusAccepted, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a no content response (204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a bad request response (400)
func BadRequest(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success:   false,
		Message:   message,
		Error:     "bad_request",
		ErrorCode: "BAD_REQUEST",
		Errors:    err,
	})
}

// Unauthorized sends an unauthorized response (401)
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success:   false,
		Message:   message,
		Error:     "unauthorized",
		ErrorCode: "UNAUTHORIZED",
	})
}

// Forbidden sends a forbidden response (403)
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success:   false,
		Message:   message,
		Error:     "forbidden",
		ErrorCode: "FORBIDDEN",
	})
}

// NotFound sends a not found response (404)
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success:   false,
		Message:   message,
		Error:     "not_found",
		ErrorCode: "NOT_FOUND",
	})
}

// Conflict sends a conflict response (409)
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success:   false,
		Message:   message,
		Error:     "conflict",
		ErrorCode: "CONFLICT",
	})
}

// UnprocessableEntity sends an unprocessable entity response (422)
func UnprocessableEntity(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success:   false,
		Message:   message,
		Error:     "unprocessable_entity",
		ErrorCode: "VALIDATION_ERROR",
		Errors:    errors,
	})
}

// TooManyRequests sends a too many requests response (429)
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success:   false,
		Message:   message,
		Error:     "too_many_requests",
		ErrorCode: "RATE_LIMIT_EXCEEDED",
	})
}

// InternalError sends an internal server error response (500)
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:   false,
		Message:   message,
		Error:     "internal_server_error",
		ErrorCode: "INTERNAL_SERVER_ERROR",
	})
}

// ServiceUnavailable sends a service unavailable response (503)
func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, Response{
		Success:   false,
		Message:   message,
		Error:     "service_unavailable",
		ErrorCode: "SERVICE_UNAVAILABLE",
	})
}

// ValidationError is a helper for formatting validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors sends validation errors in a structured format
func ValidationErrors(c *gin.Context, errors []ValidationError) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success:   false,
		Message:   "Validation failed",
		Error:     "validation_error",
		ErrorCode: "VALIDATION_ERROR",
		Errors:    errors,
	})
}