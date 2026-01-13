package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Error: http.StatusText(code), Message: message})
}

func RespondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func GetPaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}

func CalculateTotalPages(total, limit int) int {
	if limit == 0 {
		return 0
	}
	pages := total / limit
	if total%limit != 0 {
		pages++
	}
	return pages
}

func NullStringValue(ns interface{ Valid() bool; String() string }) string {
	if ns.Valid() {
		return ns.String()
	}
	return ""
}
