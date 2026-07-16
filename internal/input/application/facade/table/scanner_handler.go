package facade

import (
	"context"

	inputFacade "example/internal/input/application/facade"
	pb "example/pb/facade/table"
)

type ScannerHandler struct {
	pb.UnimplementedScannerServer
	*inputFacade.AbstractHandler
}

func NewScannerHandler(oAbstractHandler *inputFacade.AbstractHandler) *ScannerHandler {
	return &ScannerHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *ScannerHandler) AddUser(ctx context.Context, req *pb.OneRequest) (*pb.OneResponse, error) {

	return &pb.OneResponse{
		Name: "AA",
	}, nil
}
