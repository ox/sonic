package channel

import (
	"bytes"
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
	listening bool
}

func NewClient(address string) (*Client, error) {
	client := &Client{
		reading:   make([]byte, 0),
		listening: false,
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %s", address, err.Error())
	}

	client.conn = conn
	return client, err
}

func (c *Client) Connect() {
	c.Responses = make(chan string)
	c.listening = true

	go func() {
		for c.listening {
			data := make([]byte, 256)
			_, err := c.conn.Read(data)
			if err != nil && err != io.EOF {
				log.Println("Error reading from connection:", err.Error())
				c.Disconnect()
				return
			}

			// EOF means connection has already been closed. Parse the last messages
			if err != nil && err == io.EOF {
				c.ParseMessages()
				return
			}

			c.reading = append(c.reading, data[:]...)
			c.ParseMessages()
		}
	}()
}

func (c *Client) Disconnect() {
	close(c.Responses)
	c.listening = false
	c.conn.Close()
}

func (c *Client) ParseMessages() {
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

func (c *Client) Send(msg string) {
	fmt.Fprintf(c.conn, msg)
}
