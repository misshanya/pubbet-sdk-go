package pubbet

import (
	"context"
	"github.com/misshanya/pubbet-sdk-go/internal"
	pb "github.com/misshanya/pubbet/gen/go/pubbet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
)

type Reader interface {
	ListenMessages(ctx context.Context, topicName string) (<-chan []byte, error)
}

type readerClient interface {
	ListenMessages(ctx context.Context, in *pb.ListenTopicRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[pb.TopicMessages], error)
}

type conn interface {
	Close() error
}

type reader struct {
	client readerClient
	conn   conn
}

// NewReader creates a new instance of Pubbet reader
func NewReader(pubbetAddr string, creds credentials.TransportCredentials) (Reader, error) {
	reader := &reader{}

	client, conn, err := internal.NewClient(pubbetAddr, creds)
	if err != nil {
		return nil, err
	}
	reader.client = client
	reader.conn = conn

	return reader, nil
}

// ListenMessages returns a channel with the messages
func (r *reader) ListenMessages(ctx context.Context, topicName string) (<-chan []byte, error) {
	ch := make(chan []byte)

	stream, err := r.client.ListenMessages(ctx, &pb.ListenTopicRequest{TopicName: topicName})
	if err != nil {
		return nil, err
	}

	go func() {
		resp, err := stream.Recv()
		if err == io.EOF {
			close(ch)
			return
		}

		msg := resp.GetMessage()
		ch <- msg
	}()

	return ch, nil
}

// Shutdown gracefully closes connection with Pubbet server
func (r *reader) Shutdown() error {
	return r.conn.Close()
}
