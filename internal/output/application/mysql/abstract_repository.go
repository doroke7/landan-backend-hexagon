package mysql

import (
	"gorm.io/gorm"
)

type AbstractRepository struct {
	DB *gorm.DB
}

func NewAbstractRepository(oDb *gorm.DB) *AbstractRepository {
	return &AbstractRepository{
		DB: oDb,
	}
}
