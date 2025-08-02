package pubbet

import (
	"context"
	"github.com/misshanya/pubbet-sdk-go/internal"
	pb "github.com/misshanya/pubbet/gen/go/pubbet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Writer interface {
	WriteMessage(ctx context.Context, topicName string, message []byte) error
}

type writerClient interface {
	PublishMessages(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[pb.PublishMessagesRequest, pb.PublishMessagesResponse], error)
}

type writer struct {
	client writerClient
	conn   conn
	stream grpc.ClientStreamingClient[pb.PublishMessagesRequest, pb.PublishMessagesResponse]
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

// WriteMessage returns the channel to write messages in the following topic
func (w *writer) WriteMessage(ctx context.Context, topicName string, message []byte) error {
	if w.stream == nil {
		stream, err := w.client.PublishMessages(ctx)
		if err != nil {
			return err
		}
		w.stream = stream
	}

	err := w.stream.Send(&pb.PublishMessagesRequest{
		TopicName: topicName,
		Message:   message,
	})
	if err != nil {
		return err
	}

	return nil
}

// Shutdown gracefully closes connection with Pubbet server
func (w *writer) Shutdown() error {
	return w.conn.Close()
}
