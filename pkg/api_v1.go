package server

import (
	"errors"
	"fmt"
	pg_model "github.com/denysvitali/social/backend/pkg/models/postgres"
	v1requests "github.com/denysvitali/social/backend/pkg/requests/v1"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) initAPIv1(g *gin.RouterGroup) {
	g.POST("/users", s.apiV1CreateUser)
	g.GET("/users/@:username", s.apiV1UserByUsername)
	g.GET("/users/@:username/profile_picture", s.apiV1ProfilePictureByUsername)
	g.GET("/users/:id", s.apiV1GetUserById)
	g.POST("/users/:id/follows/:target_id", s.apiV1SetUserFollows)

	g.GET("/posts/@:username", s.apiV1PostsByAuthorUsername)
	g.GET("/posts/:id", s.apiV1PostsByAuthorId)
	g.GET("/tags/:text", s.apiV1TagsByText)
	g.GET("/tags/:text/posts", s.apiV1TagsGetPosts)
}

func (s *Server) apiV1GetUserById(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (s *Server) apiV1UserByUsername(c *gin.Context) {
	usernameKey := c.Param("username")
	if usernameKey == "" {
		s.badRequest(c,
			"user provided an invalid parameter username",
			"invalid parameter username",
		)
		return
	}

	var user pg_model.User
	tx := s.pgDB.First(&user, "username = ?", usernameKey)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "user not found")
			return
		}
		s.internalServerError(c, "unable to get user by username: %v", tx.Error)
		return
	}

	if user.Deleted {
		s.notFound(c, "user doesn't exist anymore")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) apiV1ProfilePictureByUsername(c *gin.Context) {
	usernameKey := c.Param("username")
	if usernameKey == "" {
		s.badRequest(c,
			"user provided an invalid parameter username",
			"invalid parameter username",
		)
		return
	}

	var pp pg_model.ProfilePicture
	tx := s.pgDB.First(&pp).
		Joins("inner join users ON users.id = profile_picture.user_id").
		Where("users.username = ?", usernameKey).
		Order("profile_picture.last_updated DESC")
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			s.notFound(c, "not found")
			return
		}
		s.internalServerError(c, "unable to find profile picture: %v", tx.Error)
		return
	}

	c.JSON(http.StatusOK, pp)
}

var NotImplementedError = map[string]string{
	"error": "not implemented",
}

func (s *Server) apiV1SetUserFollows(c *gin.Context) {
	actorUserIdKey := c.Param("id")
	if actorUserIdKey == "" {
		s.paramCantBeEmpty(c, "id")
		return
	}

	targetUserIdKey := c.Param("target_id")
	if targetUserIdKey == "" {
		s.paramCantBeEmpty(c, "target_id")
		return
	}

	c.JSON(http.StatusNotImplemented, NotImplementedError)
	return
}

func (s *Server) apiV1CreateUser(c *gin.Context) {
	var req v1requests.CreateUser
	err := c.BindJSON(&req)
	if err != nil {
		s.badRequest(c,
			fmt.Sprintf("unable to bind JSON: %v", err),
			"unable to parse JSON",
		)
		return
	}

	c.JSON(http.StatusNotImplemented, NotImplementedError)
}
