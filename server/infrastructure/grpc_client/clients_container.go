package grpcclient

import (
	"fmt"
	"log"

	"github.com/kurnhyalcantara/koer-tax-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClients struct {
	AuthServiceClient    AuthService
	AccountServiceClient AccountService
}

func InitGrpcClients(certFile string, cfg *config.Config) *ServiceClients {
	dialOptions, err := SetDialOptions(certFile)
	if err != nil {
		log.Fatal(err)
	}

	authConn, err := NewClientConn(cfg.AuthService, "Auth Service", dialOptions...)
	if err != nil {
		log.Fatal(err)
	}

	accountConn, err := NewClientConn(cfg.AccountService, "Account Service", dialOptions...)
	if err != nil {
		log.Fatal(err)
	}

	return &ServiceClients{
		AuthServiceClient:    NewAuthGrpcClient(authConn),
		AccountServiceClient: NewAccountGrpcClient(accountConn),
	}

}

func NewClientConn(address string, clientName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if address == "" {
		return nil, fmt.Errorf("%s address is empty", clientName)
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", clientName, err)
	}

	log.Printf("[service - connection] %s %s, on %s", clientName, conn.GetState().String(), address)

	return conn, nil
}

func SetDialOptions(certFile string) ([]grpc.DialOption, error) {
	var err error
	var creds credentials.TransportCredentials
	if certFile != "" {
		creds, err = credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return nil, fmt.Errorf("failed connect client tls with cert file: %v", err)
		}
	} else {
		creds = insecure.NewCredentials()
	}

	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(creds),
	)

	return opts, nil
}
