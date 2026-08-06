package usecase

import (
	helper "example/internal/helper"
	outputPort "example/internal/output/port"
)

type AbstractUsecase struct {
	outputPort.UserRepository
	*helper.AesHelper
}

func NewAbstractUsecase(oUserRepository outputPort.UserRepository, oAesHelper *helper.AesHelper) *AbstractUsecase {
	return &AbstractUsecase{
		UserRepository: oUserRepository,
		AesHelper:      oAesHelper,
	}
}
