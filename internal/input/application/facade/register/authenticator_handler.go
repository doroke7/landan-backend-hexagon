package facade

import (
	"context"

	inputFacade "example/internal/input/application/facade"
	pb "example/pb/facade/register"
)

type AuthenticatorHandler struct {
	pb.UnimplementedAuthenticatorServer
	*inputFacade.AbstractHandler
}

func NewAuthenticatorHandler(oAbstractHandler *inputFacade.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *AuthenticatorHandler) SingUp(oContext context.Context, oRequest *pb.OneRequest) (*pb.OneResponse, error) {

	return &pb.OneResponse{
		Name: "AA",
	}, nil
}
