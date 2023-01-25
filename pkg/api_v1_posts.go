package server

import (
	"errors"
	"fmt"
	"github.com/denysvitali/social/backend/pkg/models/api"
	pgmodel "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"net/http"
)

func parsePostId(c *gin.Context) (*ulid.ULID, error) {
	postId := c.Param("id")
	u, err := ulid.Parse(postId)
	if err != nil {
		return nil, fmt.Errorf("invalid ULID")
	}
	return &u, nil
}

func (s *Server) apiV1GetSinglePost(c *gin.Context) {
	postId, err := parsePostId(c)
	if err != nil {
		s.badRequest(c, fmt.Sprintf("unable to parse post id: %v", err), "invalid post id")
		return
	}

	var post pgmodel.Post
	tx := s.pgDB.
		Preload("Author").
		Model(pgmodel.Post{}).
		Where("posts.id=?", postId).
		Find(&post)

	if tx.Error != nil {
		s.internalServerError(c, "unable to fetch posts: %v", tx.Error)
		return
	}

	if tx.RowsAffected == 0 {
		s.notFound(c, "post not found")
		return
	}

	// Get Author
	a := s.getAuthor(post)
	postsResponse := api.PostsResponse{
		Posts: []api.Post{getApiPost(post)},
		Users: []api.User{a},
	}

	c.JSON(http.StatusOK, postsResponse)
}

func (s *Server) apiV1PostsByAuthorUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		s.badRequest(c, "username is empty", "username cannot be empty")
		return
	}

	var posts []pgmodel.Post
	tx := s.pgDB.
		Model(&pgmodel.Post{}).
		Joins("JOIN users ON posts.author_id = users.id").
		Where("users.username = ?", username).
		Find(&posts)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "post not found")
			return
		}
		s.internalServerError(c, "unable to get post with username %s: %v", username, tx.Error)
		return
	}

	var apiPosts []api.Post
	for _, v := range posts {
		apiPosts = append(apiPosts, getApiPost(v))
	}

	c.JSON(http.StatusOK, apiPosts)
}

func (s *Server) apiV1PostsByAuthorId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.badRequest(c, "id is empty", "id cannot be empty")
		return
	}

	var posts []pgmodel.Post
	tx := s.pgDB.
		Preload("Author").
		Preload("Tags").
		Find(&posts, "author.id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "post not found")
			return
		}
		s.internalServerError(c, "unable to get post with id %s: %v", id, tx.Error)
		return
	}

	var apiPosts []api.Post

	for _, v := range posts {
		apiPosts = append(apiPosts, getApiPost(v))
	}

	c.JSON(http.StatusOK, posts)
}
