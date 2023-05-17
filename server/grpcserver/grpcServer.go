package grpcserver

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	log                *log.Logger
	tpcServer          *grpc.Server
	serverPort         string
	isServerRegistered bool
}

func NewGrpcServer(l *log.Logger, t *grpc.Server, port string) *GrpcServer {
	return &GrpcServer{log: l, tpcServer: t, serverPort: port, isServerRegistered: false}
}

func (gs *GrpcServer) Run() {

	if !gs.isServerRegistered {
		gs.log.Println("Before run grpc server needs to call 'RegisterGrpcServer'.")
		os.Exit(1)
	}

	tcpListener, err := net.Listen("tcp", gs.serverPort)
	if err != nil {
		log.Fatal("Unable to listen tpc", err)
	}
	gs.log.Println("running on port " + gs.serverPort)

	err = gs.tpcServer.Serve(tcpListener)

	if err != nil {
		gs.log.Printf("Server error listen: %s\n", err)
		os.Exit(1)
	}
}

// sd - service description is provided by proto file auto generated
// srv - server is the Struct that implemt the interface generate by the proto file
// for example AuthServer interface{} authGrpcServeImplementation struct{}
// we need to pass de authGrpcServeImplementation reference
func (gs *GrpcServer) RegisterGrpcServer(sd grpc.ServiceDesc, srv interface{}) {
	gs.isServerRegistered = true
	gs.tpcServer.RegisterService(&sd, srv)
}

func (tv *GrpcServer) Shutdown() {
	tv.log.Println("Stopping server")
	tv.tpcServer.GracefulStop()
}
