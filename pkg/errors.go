package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) internalServerError(c *gin.Context, format string, args ...any) {
	s.logger.Errorf(format, args...)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
	})
}

func (s *Server) badRequest(c *gin.Context, warnMessage string, userMessage string) {
	s.logger.Warnf(warnMessage)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": userMessage,
	})
}

func (s *Server) notFound(c *gin.Context, message string, args ...any) {
	s.logger.Warnf(message, args...)
	c.JSON(http.StatusNotFound, gin.H{
		"error": "not found",
	})
}

func (s *Server) paramCantBeEmpty(c *gin.Context, param string) {
	s.badRequest(c,
		fmt.Sprintf("key \"%s\" is empty", param),
		fmt.Sprintf("key \"%s\" cannot be empty", param),
	)
}
