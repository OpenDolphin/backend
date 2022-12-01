package server

import (
	"errors"
	pgmodel "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) apiV1PostsByAuthorUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		s.badRequest(c, "username is empty", "username cannot be empty")
		return
	}

	var p pgmodel.Post
	tx := s.pgDB.
		Preload("Author").
		Preload("Tags").
		First(&p, "username = ?", username)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "post not found")
			return
		}
		s.internalServerError(c, "unable to get post with username %s: %v", username, tx.Error)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (s *Server) apiV1PostsByAuthorId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.badRequest(c, "id is empty", "id cannot be empty")
		return
	}

	var p pgmodel.Post
	tx := s.pgDB.
		Preload("Author").
		Preload("Tags").
		First(&p, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "post not found")
			return
		}
		s.internalServerError(c, "unable to get post with id %s: %v", id, tx.Error)
		return
	}

	c.JSON(http.StatusOK, p)
}
