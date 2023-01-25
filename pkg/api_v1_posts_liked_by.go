package server

import (
	"fmt"
	pgmodel "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) apiV1PostLikedBy(c *gin.Context) {
	postId, err := parsePostId(c)
	if err != nil {
		s.badRequest(c, fmt.Sprintf("unable to parse post id: %v", err), "invalid post id")
		return
	}
	var users []pgmodel.User
	tx := s.pgDB.
		Model(pgmodel.User{}).
		Joins("JOIN user_likes ON user_likes.user_id = users.id").
		Limit(50).
		Where("user_likes.post_id = ?", postId).
		Find(&users)
	if tx.Error != nil {
		s.internalServerError(c, "unable to find likes by post: %v", tx.Error)
		return
	}

	c.JSON(http.StatusOK, users)
}
