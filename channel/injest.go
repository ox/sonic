package channel

import (
	"fmt"
	"log"
	"strings"
)

type InjestChannel struct {
	client  *Client
	Verbose bool
}

func NewInjestChannel(address string) (*InjestChannel, error) {
	client, err := NewClient(address)
	if err != nil {
		return nil, err
	}

	channel := &InjestChannel{
		client: client,
	}

	client.Connect()

	resp := <-client.Responses
	if !strings.Contains(resp, "CONNECTED") {
		client.Disconnect()
		return nil, fmt.Errorf("Could not connect to server: %s", resp)
	}

	return channel, nil
}

// Send accepts a format string, logs it if the client is Verbose, and sends it
// to the Sonic server. It waits for a response and returns it.
func (c *InjestChannel) Send(format string, args ...interface{}) (string, error) {
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

func (c *InjestChannel) Start(password string) error {
	resp, err := c.Send("START ingest %s", password)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(resp, "STARTED ingest") {
		return fmt.Errorf("Error starting: %s", resp)
	}

	return nil
}

func (c *InjestChannel) Quit() {
	c.Send("QUIT")
	c.client.Disconnect()
}

func (c *InjestChannel) Ping() error {
	resp, err := c.Send("PING")
	if err != nil {
		return err
	}

	if resp != "PONG" {
		return fmt.Errorf("Expected PONG, got: %s", resp)
	}

	return nil
}

func (c *InjestChannel) Push(collection, bucket, object, text string) error {
	resp, err := c.Send("PUSH %s %s %s \"%s\"", collection, bucket, object, text)
	if err != nil {
		return err
	}

	if resp != "OK" {
		return fmt.Errorf("Error pushing: %s", resp)
	}

	return nil
}