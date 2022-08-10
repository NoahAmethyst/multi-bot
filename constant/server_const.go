package constant

import "time"

const (
	DefaultGRPCPort = 9090

	GrpcServer        = "GRPC_SERVER"
	ControllerService = "CONTROLLER_SERVICE"
	TStoreService     = "TSTORE_SERVICE"

	GrpcListenPort = "GRPC_LISTEN_PORT"

	TimeOut = time.Second * 60 * 5
)
