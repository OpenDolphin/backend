package server

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/denysvitali/social/backend/pkg/models/api"
	pg_model "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"time"
)

type Server struct {
	logger *logrus.Logger
	e      *gin.Engine
	pgDB   *gorm.DB

	// isDemo defines whether the server is running in demo mode: when this mode is enabled, the DB is
	// pre-filled with demo data.
	isDemo bool
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
	DemoMode    bool

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

	pgdb.Logger = logger.Default.LogMode(logger.Info)

	s := Server{
		e:      gin.New(),
		pgDB:   pgdb,
		logger: config.Logger,
	}

	if config.DemoMode {
		s.isDemo = true
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

	if s.isDemo {
		s.logger.Info("Filling DB with demo data")
		err := s.addDemoData()
		if err != nil {
			s.logger.Fatalf("unable to add demo data: %v", err)
		}
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	s.e.Use(cors.New(corsConfig))
	s.e.Use(gin.Logger())

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

func (s *Server) apiV1GetPosts(c *gin.Context) {
	var posts []pg_model.Post

	tx := s.pgDB.
		Model(&posts).
		Preload("Author").
		Joins("INNER JOIN user_likes ON user_likes.post_id = posts.id").
		Group("posts.id").
		Select("posts.*, COUNT(user_likes.post_id) AS likes").
		Limit(50).
		Order("posts.id DESC").
		Scan(&posts)
	if tx.Error != nil {
		s.internalServerError(c, "unable to fetch posts: %v", tx.Error)
		return
	}

	var apiPosts []api.Post
	var postsResponse api.PostsResponse

	authorsMap := map[uint64]bool{}

	for _, p := range posts {
		var ulidBytes [16]byte
		copy(ulidBytes[:], p.ID[:16])
		pUlid := ulid.ULID(ulidBytes)
		apiPosts = append(apiPosts, api.Post{
			ID:        pUlid.String(),
			Content:   p.Content,
			Likes:     p.Likes,
			Author:    p.AuthorID,
			CreatedAt: time.Unix(int64(pUlid.Time()/1000), 0),
		})
		authorsMap[p.AuthorID] = true
	}

	// Fetch Authors
	var authorIds []uint64
	for k := range authorsMap {
		authorIds = append(authorIds, k)
	}

	var apiUsers []api.User
	var authors []pg_model.User
	tx = s.pgDB.Where("id IN ?", authorIds).Find(&authors)
	if tx.Error != nil {
		s.internalServerError(c, "unable to find authors: %v", tx.Error)
		return
	}

	for _, u := range authors {
		apiUsers = append(apiUsers, api.User{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			Username:    u.Username,
			Verified:    u.Verified,
		})
	}

	postsResponse.Posts = apiPosts
	postsResponse.Users = apiUsers

	c.JSON(http.StatusOK, postsResponse)
}

func (s *Server) addDemoData() error {

	tx := s.pgDB.Raw("TRUNCATE users CASCADE").Scan(nil)
	if tx.Error != nil {
		return tx.Error
	}

	err := s.createDemoUsers()
	if err != nil {
		return err
	}

	err = s.createDemoPosts()
	if err != nil {
		return err
	}
	return nil
}
