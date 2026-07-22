package authentication

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pbFacadeAdminAuthentication "example/pb/facade/admin/authentication"
)

type AuthenticatorHandler struct {
	pbFacadeAdminAuthentication.UnimplementedAuthenticatorServer
	*inputApplicationFacade.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationFacade.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oContext context.Context, oRequest *pbFacadeAdminAuthentication.OneRequest) (*pbFacadeAdminAuthentication.OneResponse, error) {
	sAuthorization, err := oSelf.AuthenticatorUsecase.SignIn(oRequest.Name, oRequest.Password)
	if err != nil {
		return nil, err
	}

	if err := grpc.SetHeader(oContext, metadata.Pairs("authorization", sAuthorization)); err != nil {
		return nil, err
	}

	return &pbFacadeAdminAuthentication.OneResponse{}, nil
}
