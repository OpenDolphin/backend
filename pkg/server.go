package server

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	pg_model "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	logger *logrus.Logger
	e      *gin.Engine
	pgDB   *gorm.DB
}

type ArangoConfig struct {
	Endpoints []string
	Username  string
	Password  string

	Database string
}

type Config struct {
	Arango      ArangoConfig
	PostgresDSN string

	Logger *logrus.Logger
}

func New(config Config) (*Server, error) {
	if config.Logger == nil {
		config.Logger = logrus.New()
		config.Logger.Warnf("nil logger passed in config, creating a new logger")
	}

	pgdb, err := setupPostgres(config)
	if err != nil {
		return nil, fmt.Errorf("unable to set-up PostgreSQL")
	}

	s := Server{
		e:      gin.New(),
		pgDB:   pgdb,
		logger: config.Logger,
	}

	s.init()

	return &s, nil
}

func setupPostgres(config Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(config.PostgresDSN), &gorm.Config{})
}

func setupArango(config Config) (driver.Client, driver.Database, error) {
	conn, err := arangohttp.NewConnection(
		arangohttp.ConnectionConfig{
			Endpoints: config.Arango.Endpoints,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create arango HTTP connection: %v", err)
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
		return nil, nil, fmt.Errorf("unable to create ArangoDB client: %v", err)
	}

	db, err := c.Database(context.TODO(), config.Arango.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get ArangoDB database: %v", err)
	}
	return c, db, nil
}

func (s *Server) Listen(addr ...string) error {
	return s.e.Run(addr...)
}

func (s *Server) init() {
	// init db
	s.initPostgreSQL()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	s.e.Use(cors.New(corsConfig))

	g := s.e.Group("/api/v1")
	s.initAPIv1(g)
}

func (s *Server) initPostgreSQL() {
	for _, v := range []any{
		&pg_model.User{},
		&pg_model.ProfilePicture{},
		&pg_model.Post{},
		&pg_model.Tag{},
	} {
		err := s.pgDB.AutoMigrate(v)
		if err != nil {
			s.logger.Fatalf("unable to automigrate %t: %v", v, err)
		}
	}
}
