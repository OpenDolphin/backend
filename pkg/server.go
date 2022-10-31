package server

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger *logrus.Logger
	e      *gin.Engine

	graphConn driver.Client
	graphDb   driver.Database
}

type ArangoConfig struct {
	Endpoints []string
	Username  string
	Password  string

	Database string
}

type Config struct {
	Arango ArangoConfig
	Logger *logrus.Logger
}

func New(config Config) (*Server, error) {
	conn, err := arangohttp.NewConnection(
		arangohttp.ConnectionConfig{
			Endpoints: config.Arango.Endpoints,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create arango HTTP connection: %v", err)
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
		Authentication: driver.BasicAuthentication(
			config.Arango.Username,
			config.Arango.Password,
		),
		SynchronizeEndpointsInterval: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create ArangoDB client: %v", err)
	}

	if config.Logger == nil {
		config.Logger = logrus.New()
		config.Logger.Warnf("nil logger passed in config, creating a new logger")
	}

	db, err := c.Database(context.TODO(), config.Arango.Database)
	if err != nil {
		return nil, fmt.Errorf("unable to get ArangoDB database: %v", err)
	}

	s := Server{
		e:         gin.New(),
		graphConn: c,
		graphDb:   db,
		logger:    config.Logger,
	}

	s.init()

	return &s, nil
}

func (s *Server) Listen(addr ...string) error {
	return s.e.Run(addr...)
}

func (s *Server) init() {
	g := s.e.Group("/api/v1")
	s.initAPIv1(g)
}
