package authentication

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	bootstrap "example/bootstrap"
	inputApplicationFacade "example/internal/input/application/facade"
	interceptorFacadeAdmin "example/internal/interceptor/facade/admin"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pbFacadeAdminAuthentication "example/pb/facade/admin/authentication"
	"example/pkg"
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

	sAuthorization, oErr := oSelf.AuthenticatorUsecase.SignIn(oRequest.Name, oRequest.Password, bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.SECRET)

	if oErr != nil {
		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {
			return nil, status.Error(codes.Aborted, oDefaultError.Error())
		}

		return nil, status.Error(codes.Internal, oErr.Error())
	}

	if oHolder, bOk := interceptorFacadeAdmin.AuthHolderFromContext(oContext); bOk {
		oHolder.Authorization = sAuthorization
	}

	if oErr := grpc.SetHeader(oContext, metadata.Pairs("A", sAuthorization)); oErr != nil {
		return nil, oErr
	}

	return &pbFacadeAdminAuthentication.OneResponse{}, nil
}
