package generator_handler

import (
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	pgsql "github.com/xinliangnote/go-gin-api/internal/repository/pgsql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"

	"go.uber.org/zap"
)

type handler struct {
	db     mysql.Repo
	logger *zap.Logger
	cache  redis.Repo
	pgdb   pgsql.Repo
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo, pgdb pgsql.Repo) *handler {
	return &handler{
		logger: logger,
		cache:  cache,
		db:     db,
		pgdb:   pgdb,
	}
}
