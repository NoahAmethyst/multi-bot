package grpc_server

import (
	"context"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/utils/log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

// ProvideHTTP convert grpc into http
// so that grpc server can support http request
func ProvideHTTP(endpoint string, grpcServer *grpc.Server) *http.Server {

	//Create a new gwmux，It is a request multiplexer for grpc-gateway.
	//It matches the http request to the pattern and calls the corresponding handler.
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := alliance_bot_pb.RegisterBotServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, opts)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "register endpoint",
			"error":  err,
		}).Send()
		return nil
	}
	//Create a new mux, which is a request multiplexer for http
	mux := http.NewServeMux()

	//register gmunx
	mux.Handle("/", gwmux)
	log.Info().Fields(map[string]interface{}{
		"action":   "provide http transport to grpc server",
		"endPoint": endpoint,
	}).Send()
	return &http.Server{
		Addr:    endpoint,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}
}

// grpcHandlerFunc redirect different request into different handler
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
