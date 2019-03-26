package channel

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
	Responses chan string
	reading   []byte
	conn      net.Conn
}

func NewClient(address string) (*Client, error) {
	client := &Client{
		reading: make([]byte, 0),
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %s", address, err.Error())
	}

	client.conn = conn
	client.Responses = make(chan string)
	return client, err
}

func (c *Client) Connect(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.conn.Close()
				return
			default:
				data := make([]byte, 256)
				_, err := c.conn.Read(data)
				if err != nil && err != io.EOF {
					log.Println("Error reading from connection:", err.Error())
					c.conn.Close()
					return
				}

				// EOF means connection has already been closed.
				if err != nil && err == io.EOF {
					return
				}

				c.reading = append(c.reading, data[:]...)
			}
		}
	}()
}

func (c *Client) ParseMessages(ctx context.Context) {
	messages := bytes.Split(c.reading, []byte("\n"))
	for i, message := range messages {
		if i == len(messages)-1 {
			continue
		}

		trimmedBytes := bytes.Trim(message, "\x00")
		c.Responses <- strings.TrimSpace(string(trimmedBytes))
	}

	// keep the last read bytes in the buffer in case it's only part of the next response
	c.reading = messages[len(messages)-1]
}

func (c Client) Send(msg string) {
	fmt.Fprintf(c.conn, msg)
}
