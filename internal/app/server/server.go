package server

import (
	"context"

	"github.com/dwikalam/calcgorpc/internal/app/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedCalculatorServer
}

func New() *Server {
	return &Server{}
}

func (s *Server) Add(
	ctx context.Context,
	in *pb.CalculationRequest,
) (
	*pb.CalculationResponse,
	error,
) {
	return &pb.CalculationResponse{
		Result: in.A + in.B,
	}, nil
}

func (s *Server) Substract(
	ctx context.Context,
	in *pb.CalculationRequest,
) (
	*pb.CalculationResponse,
	error,
) {
	return &pb.CalculationResponse{
		Result: in.A - in.B,
	}, nil
}

func (s *Server) Multiply(
	ctx context.Context,
	in *pb.CalculationRequest,
) (
	*pb.CalculationResponse,
	error,
) {
	return &pb.CalculationResponse{
		Result: in.A * in.B,
	}, nil
}

func (s *Server) Divide(
	ctx context.Context,
	in *pb.CalculationRequest,
) (
	*pb.CalculationResponse,
	error,
) {
	if in.B == 0 {
		return nil,
			status.Error(
				codes.InvalidArgument,
				"cannot divide by zero",
			)
	}

	return &pb.CalculationResponse{
			Result: in.A / in.B,
		},
		nil
}
