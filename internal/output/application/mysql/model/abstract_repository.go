package output_application_mysql

import (
	"gorm.io/gorm"
)

type AbstractRepository struct {
	db *gorm.DB
}

func NewAbstractRepository(oDb *gorm.DB) *AbstractRepository {
	return &AbstractRepository{
		db: oDb,
	}
}
