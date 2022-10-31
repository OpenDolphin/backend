package server

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/denysvitali/social/backend/pkg/models"
	v1requests "github.com/denysvitali/social/backend/pkg/requests/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) initAPIv1(g *gin.RouterGroup) {
	g.GET("/users/by_username/:username", s.apiV1UserByUsername)
	g.GET("/users/:id", s.apiV1Users)
	g.POST("/users", s.apiV1CreateUser)
	g.POST("/users/:id/follows/:target_id", s.apiV1SetUserFollows)
}

func (s *Server) apiV1Users(c *gin.Context) {
	userKey := c.Param("id")
	if userKey == "" {
		s.paramCantBeEmpty(c, "id")
		return
	}

	ctx := context.TODO()
	coll, err := s.graphDb.Collection(ctx, UsersCollection)
	if err != nil {
		s.internalServerError(c, "unable to get collection: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var res map[string]any
	docMeta, err := coll.ReadDocument(ctx, userKey, &res)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			s.notFound(c, "requested user \"%s\" but not found", userKey)
			return
		}
		s.internalServerError(c, "unable to read document: %v", err)
		return
	}

	s.logger.Debugf("docMeta: %v", docMeta)
	c.JSON(http.StatusOK, res)
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

	ctx := context.TODO()
	g, err := s.graphDb.Graph(ctx, SocialNetworkGraph)
	if err != nil {
		s.logger.Errorf("unable to get graph \"%s\": %v", SocialNetworkGraph, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	coll, v, err := g.EdgeCollection(ctx, SocialNetworkRelations)
	if err != nil {
		s.logger.Errorf("unable to get edge collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	s.logger.Debugf("v=%v", v)

	ctx = driver.WithQueryCount(ctx)

	// Check if relation already exists (we don't want duplicates)
	cursor, err := s.graphDb.Query(ctx, `FOR v, e, p in 1..1 OUTBOUND @actor_id GRAPH "social_network"
FILTER e.label == "follows" AND e._from == @actor_id AND e._to == @target_id
RETURN v`,
		map[string]any{
			"actor_id":  fmt.Sprintf("users/%s", actorUserIdKey),
			"target_id": fmt.Sprintf("users/%s", targetUserIdKey),
		})

	if err != nil {
		s.logger.Errorf("unable to perform query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	err = cursor.Close()
	if err != nil {
		s.logger.Errorf("unable to close cursor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	if cursor.Count() >= 1 {
		s.logger.Warnf("trying to follow, but user is already being followed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already being followed",
		})
		return
	}

	docMeta, err := coll.CreateDocument(ctx, map[string]string{
		"_from": fmt.Sprintf("users/%s", actorUserIdKey),
		"_to":   fmt.Sprintf("users/%s", targetUserIdKey),
		"label": "follows",
	})
	if err != nil {
		s.logger.Errorf("unable to create document for follow relationship: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	s.logger.Debugf("docMeta=%v", docMeta)
	c.JSON(http.StatusOK, docMeta)
	return
}

func (s *Server) apiV1CreateUser(c *gin.Context) {
	ctx := context.TODO()

	var req v1requests.CreateUser
	err := c.BindJSON(&req)
	if err != nil {
		s.badRequest(c,
			fmt.Sprintf("unable to bind JSON: %v", err),
			"unable to parse JSON",
		)
		return
	}

	coll, err := s.graphDb.Collection(ctx, UsersCollection)
	if err != nil {
		s.internalServerError(c, "unable to get users collection: %v", err)
		return
	}

	meta, err := coll.CreateDocument(ctx, models.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if driver.IsConflict(err) {
			s.badRequest(c,
				fmt.Sprintf("username conflict: %v", err),
				"an user with this username already exists",
			)
			return
		}
		s.internalServerError(c, "unable to create document: %v", err)
		return
	}

	c.JSON(http.StatusOK, meta)
}