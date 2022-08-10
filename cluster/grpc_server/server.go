package grpc_server

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"multi-bot/cluster/bot_service"
	"multi-bot/constant"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/utils/log"
	"net"
	"runtime/debug"
	"strconv"
)

var (
	customFunc grpc_recovery.RecoveryHandlerFuncContext
)

var GrpcServer *grpc.Server

func Start(grpcPort string) {

	if len(grpcPort) == 0 {
		grpcPort = strconv.Itoa(constant.DefaultGRPCPort)
	}

	grpcAddr := fmt.Sprintf("0.0.0.0:%s", grpcPort)
	lis, err := net.Listen("tcp", grpcAddr)

	if len(grpcPort) > 0 {
		grpcPort = fmt.Sprintf("0.0.0.0:%s", grpcPort)
	}

	if err != nil {
		log.Error().Fields(map[string]interface{}{"action": "grpc listener error", "error": err.Error()}).Send()

	}
	log.Info().Msgf("bot service start at address %s", grpcAddr)

	// Define customfunc to handle panic
	//Todo panic error
	customFunc = func(ctx context.Context, p interface{}) error {
		log.Error().Msgf("[PANIC] %s\n\n%s", p, string(debug.Stack()))
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(customFunc),
	}

	// Create a server. Recovery handlers should typically be last in the chain_info so that other middleware
	// (e.g. logging) can operate on the recovered state instead of being directly affected by any panic

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
			//otgrpc.OpenTracingServerInterceptor(thisTracer),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
			//grpc_opentracing.StreamServerInterceptor(topts...),
		),
	)

	//register bot server
	alliance_bot_pb.RegisterBotServiceServer(grpcServer, bot_service.Server{})
	GrpcServer = grpcServer

	reflection.Register(grpcServer)

	if thisServer := ProvideHTTP(grpcAddr, grpcServer); thisServer != nil && len(grpcAddr) > 0 {
		err = thisServer.Serve(lis)
	} else {
		err = grpcServer.Serve(lis)
	}

	if err != nil {
		log.Error().Err(err)
	}
}
