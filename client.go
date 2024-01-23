package ariesocketio

import (
	"net"
	"strconv"

	"github.com/arielitovsky/ariesocketio/protocol"
	"github.com/arielitovsky/ariesocketio/websocket"
)

const (
	webSocketProtocol       = "ws://"
	webSocketSecureProtocol = "wss://"
	socketIOUrl             = "/socket.io/?transport=websocket"
)

type Option func(c *Client)

/*
*  Set the data handler that should be used when using a CONNECT command
* the handler should return a data struct that will be used
 */
func WithConnectData(handler ConnectDataHandler) Option {
	return func(c *Client) {
		c.methods.SetConnectDataHandler(handler)
	}
}

type Client struct {
	methods
	Channel
}

func GetUrl(host string, port int, secure bool) string {
	var prefix string

	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}

	return prefix + net.JoinHostPort(host, strconv.Itoa(port)) + socketIOUrl
}

func Dial(url string, tr websocket.Transport, options ...Option) (*Client, error) {

	c := &Client{}
	c.initChannel()

	for _, opt := range options {
		opt(c)
	}

	var err error

	if tr.Protocol == protocol.Protocol3 {
		url = url + "&EIO=3"
	} else if tr.Protocol == protocol.Protocol4 {
		url = url + "&EIO=4"
	} else {
		url = url + "&EIO=4"
	}

	c.conn, err = tr.Connect(url)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)

	return c, nil
}

func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
