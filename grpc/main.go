package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "test/grpc/test"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		server()
	}()

	time.Sleep(1 * time.Second)

	client()
}

func zaplogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := cfg.Build()
	defer logger.Sync()

	return logger
}

type pbServer struct {
	pb.UnimplementedGreeterServer
	logger *zap.Logger
}

func (s *pbServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func server() {
	logger := zaplogger()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &pbServer{logger: logger})
	logger.Info("Server is running on port :50052")
	if err := s.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}

func client() {
	logger := zaplogger()

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Failed to connect", zap.Error(err))
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	numOfRequests := 10000

	start := time.Now()
	for i := 0; i < numOfRequests; i++ {
		_, err := c.SayHello(context.Background(), &pb.HelloRequest{
			Name: "World",
		})
		if err != nil {
			logger.Fatal("Failed to call", zap.Error(err))
		}
	}

	end := time.Now()

	logger.Info(fmt.Sprintf("GRPC time taken for %d Unary calls", numOfRequests), zap.Duration("duration", end.Sub(start)))
}
