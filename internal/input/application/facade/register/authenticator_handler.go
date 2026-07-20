package input_application_facade

import (
	"context"

	inputApplicationFacade "example/internal/input/application/facade"
	pbFacadeRegister "example/pb/facade/register"
)

type AuthenticatorHandler struct {
	pbFacadeRegister.UnimplementedAuthenticatorServer
	*inputApplicationFacade.AbstractHandler
}

func NewAuthenticatorHandler(oAbstractHandler *inputApplicationFacade.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *AuthenticatorHandler) SingUp(oContext context.Context, oRequest *pbFacadeRegister.OneRequest) (*pbFacadeRegister.OneResponse, error) {

	return &pbFacadeRegister.OneResponse{
		Name: "AA",
	}, nil
}
