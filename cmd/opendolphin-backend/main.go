package main

import (
	"github.com/alexflint/go-arg"
	"github.com/denysvitali/social/backend/pkg"
	"github.com/sirupsen/logrus"
)

var args struct {
	Debug *bool `arg:"-D"`

	ArangoEndpoints []string `arg:"--arango-endpoints,env:ARANGO_ENDPOINTS"`
	ArangoUsername  string   `arg:"--arango-username,env:ARANGO_USERNAME"`
	ArangoPassword  string   `arg:"--arango-password,env:ARANGO_PASSWORD"`
	ArangoDatabase  string   `arg:"--arango-database,env:ARANGO_DATABASE"`
	PostgresDSN     string   `arg:"--postgres-dsn,env:DATABASE_URL"`

	IsDemo bool `arg:"env:DEMO_MODE" default:"false"`

	ListenAddr string `arg:"--listen-addr,env:LISTEN_ADDR"`
}

var logger = logrus.New()

func main() {
	arg.MustParse(&args)

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	s, err := server.New(server.Config{
		PostgresDSN: args.PostgresDSN,
		Logger:      logger,
		DemoMode:    args.IsDemo,
	})

	if err != nil {
		logger.Fatalf("unable to create server: %v", err)
	}

	err = s.Listen(args.ListenAddr)
	if err != nil {
		logger.Fatalf("unable to listen: %v", err)
	}
}
