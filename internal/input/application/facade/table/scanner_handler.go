package facade

import (
	"context"

	inputApplicationFacade "example/internal/input/application/facade"
	pbFacadeTable "example/pb/facade/table"
)

type ScannerHandler struct {
	pbFacadeTable.UnimplementedScannerServer
	*inputApplicationFacade.AbstractHandler
}

func NewScannerHandler(oAbstractHandler *inputApplicationFacade.AbstractHandler) *ScannerHandler {
	return &ScannerHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *ScannerHandler) AddUser(ctx context.Context, req *pbFacadeTable.OneRequest) (*pbFacadeTable.OneResponse, error) {

	return &pbFacadeTable.OneResponse{
		Name: "AA",
	}, nil
}
