package postgre

import (
	"context"

	"github.com/cyclex/planet-ban/domain/repository"
	"gorm.io/gorm"
)

type postgreRepo struct {
	DB          *gorm.DB
	c           context.Context
	waitingTime int64
}

func NewPostgreRepository(c context.Context, db *gorm.DB) repository.ModelRepository {
	return &postgreRepo{
		DB: db,
		c:  c,
	}
}
