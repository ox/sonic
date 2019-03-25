package channel

import (
	"fmt"
	"log"
	"strings"
)

type SearchChannel struct {
	client  *Client
	Verbose bool
}

func NewSearchChannel(address string) (*SearchChannel, error) {
	client, err := NewClient(address)
	if err != nil {
		return nil, err
	}

	channel := &SearchChannel{
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

func (c *SearchChannel) Send(format string, args ...interface{}) (string, error) {
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

func (c *SearchChannel) Start(password string) error {
	resp, err := c.Send("START search %s", password)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(resp, "STARTED search") {
		return fmt.Errorf("Error starting: %s", resp)
	}

	return nil
}

func (c *SearchChannel) Quit() {
	c.Send("QUIT")
	c.client.Disconnect()
}

func (c *SearchChannel) Ping() error {
	resp, err := c.Send("PING")
	if err != nil {
		return err
	}

	if resp != "PONG" {
		return fmt.Errorf("Expected PONG, got: %s", resp)
	}

	return nil
}
