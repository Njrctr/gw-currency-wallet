package handlers

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, errorResponse{Message: message})
}
