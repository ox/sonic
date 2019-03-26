package channel

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type IngestChannel struct {
	client     *Client
	cancelFunc context.CancelFunc
	Verbose    bool
}

func NewIngestChannel(address string) (*IngestChannel, error) {
	client, err := NewClient(address)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	channel := &IngestChannel{
		client:     client,
		cancelFunc: cancelFunc,
	}

	client.Connect(ctx)
	client.ParseMessages(ctx)

	resp := <-client.Responses
	if !strings.Contains(resp, "CONNECTED") {
		cancelFunc()
		return nil, fmt.Errorf("Could not connect to server: %s", resp)
	}

	return channel, nil
}

// Send accepts a format string, logs it if the client is Verbose, and sends it
// to the Sonic server. It waits for a response and returns it.
func (c *IngestChannel) Send(format string, args ...interface{}) (string, error) {
	Send := fmt.Sprintf(fmt.Sprintf("%s\n", format), args...)
	if c.Verbose {
		log.Print(Send)
	}

	c.client.Send(Send)
	resp := <-c.client.Responses

	if c.Verbose {
		log.Println(resp)
	}

	return resp, nil
}

func (c *IngestChannel) Start(password string) error {
	resp, err := c.Send("START ingest %s", password)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(resp, "STARTED ingest") {
		return fmt.Errorf("Error starting: %s", resp)
	}

	return nil
}

func (c *IngestChannel) Quit() {
	c.Send("QUIT")
	c.cancelFunc()
}

func (c *IngestChannel) Ping() error {
	resp, err := c.Send("PING")
	if err != nil {
		return err
	}

	if resp != "PONG" {
		return fmt.Errorf("Expected PONG, got: %s", resp)
	}

	return nil
}

func (c *IngestChannel) Push(collection, bucket, object, text string) error {
	resp, err := c.Send("PUSH %s %s %s \"%s\"", collection, bucket, object, text)
	if err != nil {
		return err
	}

	if resp != "OK" {
		return fmt.Errorf("Error pushing: %s", resp)
	}

	return nil
}
