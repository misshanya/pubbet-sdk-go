package internal

import (
	"fmt"
	pb "github.com/misshanya/pubbet/gen/go/pubbet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewClient(addr string, creds credentials.TransportCredentials) (pb.PubbetClient, *grpc.ClientConn, error) {
	grpcConn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(
			creds,
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init connection to the pubbet: %w", err)
	}

	grpcClient := pb.NewPubbetClient(grpcConn)
	return grpcClient, grpcConn, nil
}
