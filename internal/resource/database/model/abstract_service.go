package service

import (
	"landan-backend-grpc/pkg"
)

func NewAbstractService(oDatabaseFactory *pkg.DatabaseFactory) *AbstractService {
	return &AbstractService{
		DatabaseFactory: oDatabaseFactory, // 因為controller 的 abstract_controller 目錄不同，小寫屬性不能共享
	}
}

type AbstractService struct {
	DatabaseFactory *pkg.DatabaseFactory
}
