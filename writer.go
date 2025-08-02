package pubbet

import (
	"context"
	"errors"
	"github.com/misshanya/pubbet-sdk-go/internal"
	pb "github.com/misshanya/pubbet/gen/go/pubbet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Writer interface {
	WriteMessages(ctx context.Context, messages []*Message) error
}

type writerClient interface {
	PublishMessages(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[pb.PublishMessagesRequest, pb.PublishMessagesResponse], error)
}

type writer struct {
	client writerClient
	conn   conn
}

// NewWriter creates a new instance of Pubbet writer
func NewWriter(pubbetAddr string, creds credentials.TransportCredentials) (Writer, error) {
	writer := &writer{}

	client, conn, err := internal.NewClient(pubbetAddr, creds)
	if err != nil {
		return nil, err
	}
	writer.client = client
	writer.conn = conn

	return writer, nil
}

// WriteMessages sends messages to Pubbet server
func (w *writer) WriteMessages(ctx context.Context, messages []*Message) error {
	stream, err := w.client.PublishMessages(ctx)
	if err != nil {
		return err
	}

	var allErrors error

	for _, msg := range messages {
		err := stream.Send(&pb.PublishMessagesRequest{
			TopicName: msg.TopicName,
			Message:   msg.Data,
		})
		if err != nil {
			allErrors = errors.Join(allErrors, err)
		}
	}

	return allErrors
}

// Shutdown gracefully closes connection with Pubbet server
func (w *writer) Shutdown() error {
	return w.conn.Close()
}
