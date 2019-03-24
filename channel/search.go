package channel

import (
	"fmt"
	"log"
	"net"
)

type SearchChannel struct {
	conn    net.Conn
	Verbose bool
}

func NewSearchChannel(address string) (*SearchChannel, error) {
	channel := &SearchChannel{}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %s", address, err.Error())
	}

	channel.conn = conn
	return channel, nil
}

func (c *SearchChannel) msg(format string, args ...string) string {
	msg := fmt.Sprintf(format, args)
	if c.Verbose {
		log.Println(msg)
	}
	return msg
}

func (c *SearchChannel) Start(password string) {
	fmt.Fprintf(c.conn, "START search %s", password)
}

func (c *SearchChannel) Quit() {
	fmt.Fprintf(c.conn, "QUIT")
	c.conn.Close()
}
