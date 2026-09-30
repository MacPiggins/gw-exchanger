package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/MacPiggins/gw-exchanger/internal/logging"
	"github.com/MacPiggins/gw-exchanger/internal/service"
	pb "github.com/MacPiggins/gw-proto/go/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedExchangeServiceServer
	exchangeService *service.ExchangeService
}

func New(exchangeService *service.ExchangeService) *GrpcServer {
	return &GrpcServer{exchangeService: exchangeService}
}

func (s *GrpcServer) GetExchangeRate(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	logging.AppendCtx(ctx, slog.Group(
		"Grpc",
		slog.String("method", "GetExchangeRate"),
		slog.Group(
			"values",
			slog.String("from", req.From),
			slog.String("to", req.To),
		),
	))
	// some request data checks?
	slog.InfoContext(ctx, "grpc request")
	rate, err := s.exchangeService.GetExchangeRate(ctx, req.From, req.To)
	if err != nil {
		slog.ErrorContext(ctx, "error while processing request", slog.Any("error", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &pb.ExchangeRateResponse{From: req.From, To: req.To, Rate: float32(rate)}, nil
}

func (s *GrpcServer) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	logging.AppendCtx(ctx, slog.Group(
		"Grpc",
		slog.String("method", "GetExchangeRates"),
		slog.String(
			"values", "empty",
		),
	))
	slog.InfoContext(ctx, "grpc request")
	rates, err := s.exchangeService.GetExchangeRates(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error while processing request", slog.Any("error", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &pb.ExchangeRatesResponse{Rates: rates}, nil
}

func (s *GrpcServer) Run(lis net.Listener) error {
	slog.Info("starting grpc server", slog.String("address", lis.Addr().String()))
	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, s)
	err := grpcServer.Serve(lis)
	if err != nil {
		slog.Error("error while running server", slog.Any("error", err))
	}
	return err
}
