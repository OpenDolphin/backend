package server

import (
	"errors"
	pgmodel "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) apiV1TagsByText(c *gin.Context) {
	text := c.Param("text")
	if text == "" {
		s.badRequest(c, "text is empty", "text cannot be empty")
		return
	}

	var p pgmodel.Tag
	tx := s.pgDB.
		First(&p, "text = ?", text)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "post not found")
			return
		}
		s.internalServerError(c, "unable to get tags with text %s: %v", text, tx.Error)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (s *Server) apiV1TagsGetPosts(c *gin.Context) {
	text := c.Param("text")
	if text == "" {
		s.badRequest(c, "text is empty", "text cannot be empty")
		return
	}

	var p []pgmodel.Post
	tx := s.pgDB.
		Raw(`
		SELECT * FROM posts WHERE posts.id IN (
		SELECT DISTINCT s1.post_id FROM 
		(
			SELECT posts.*, post_tags.* FROM "posts" 
		    INNER JOIN post_tags ON post_tags.post_id = posts.id
		) as s1 
		INNER JOIN tags on tags.id = s1.tag_id 
		WHERE tags.text = ?
		LIMIT ?
		)`, text, 100).Preload("Tags").Find(&p)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "posts not found")
			return
		}
		s.internalServerError(c, "unable to get tags with text %s: %v", text, tx.Error)
		return
	}

	c.JSON(http.StatusOK, p)
}
