package server

import (
	"errors"
	pgmodel "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) apiV1GetUsers(c *gin.Context) {
	var users []pgmodel.User
	tx := s.pgDB.Limit(50).Find(&users)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "post not found")
			return
		}
		s.internalServerError(c, "unable to get users: %v", tx.Error)
		return
	}

	c.JSON(http.StatusOK, users)
}
