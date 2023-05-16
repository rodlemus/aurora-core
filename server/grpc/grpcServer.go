package grpc

import (
	"context"
	"log"
	"net"
	"os"

	protostv "github.com/rodlemus/aurora-services/auth-service/protos/tokenvalidation"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	log        *log.Logger
	tpcServer  *grpc.Server
	serverPort string
	protostv.UnimplementedTokenValidationServer
}

func NewGrpcServer(l *log.Logger, t *grpc.Server, port string) *GrpcServer {
	return &GrpcServer{log: l, tpcServer: t, serverPort: port}
}


func (tv *GrpcServer) Run() {
	protostv.RegisterTokenValidationServer(tv.tpcServer, tv)
	tcpListener, err := net.Listen("tcp", tv.serverPort)
	if err != nil {
		log.Fatal("Unable to listen tpc", err)
	}
	tv.log.Println("running on port " + tv.serverPort)

	err = tv.tpcServer.Serve(tcpListener)

	if err != nil {
		tv.log.Printf("Server error listen: %s\n", err)
		os.Exit(1)
	}
}

func (gs *GrpcServer) registerServer(serverRegister grpc.ServiceRegistrar, sd grpc.ServiceDesc, srv interface{}) {
	serverRegister.RegisterService(sd, srv)
}

func (tv *GrpcServer) Shutdown() {
	tv.log.Println("Stopping server")
	tv.tpcServer.GracefulStop()
}
