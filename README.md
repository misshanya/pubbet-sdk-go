# Pubbet SDK for Go

A Golang library that provides a simple way to interact with [Pubbet](https://github.com/misshanya/pubbet).

## Usage

```shell
go get github.com/misshanya/pubbet-sdk-go
```

### Producer code example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/misshanya/pubbet-sdk-go"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	writer, err := pubbet.NewWriter("localhost:5000", insecure.NewCredentials())
	if err != nil {
		fmt.Println("failed to init writer:", err)
		os.Exit(1)
	}
	defer writer.Shutdown()

	ctx := context.Background()

	messages := []*pubbet.Message{
		{
			TopicName: "best-topic-ever",
			Data:      []byte("Hello world from first message!"),
		},
		{
			TopicName: "best-topic-ever",
			Data:      []byte("Second message"),
		},
	}

	err = writer.WriteMessages(ctx, messages)
	if err != nil {
		fmt.Println("failed to write message:", err)
		os.Exit(1)
	}
}
```

### Consumer code example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/misshanya/pubbet-sdk-go"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	reader, err := pubbet.NewReader("localhost:5000", insecure.NewCredentials())
	if err != nil {
		fmt.Println("failed to init reader:", err)
		os.Exit(1)
	}
	defer reader.Shutdown()

	ctx := context.Background()

	ch, err := reader.ListenMessages(ctx, "best-topic-ever")
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}

	for v := range ch {
		fmt.Println("received msg:", string(v))
	}
}
```

> [!NOTE]
> This consumer will keep listening to new messages until the ctx is cancelled or connection is closed