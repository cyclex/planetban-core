package postgre

import (
	"context"

	"github.com/cyclex/planet-ban/domain/repository"
	"gorm.io/gorm"
)

type postgreRepo struct {
	DB      *gorm.DB
	c       context.Context
	UrlHost string
}

func NewPostgreRepository(c context.Context, db *gorm.DB, urlHost string) repository.ModelRepository {
	return &postgreRepo{
		DB:      db,
		c:       c,
		UrlHost: urlHost,
	}
}
