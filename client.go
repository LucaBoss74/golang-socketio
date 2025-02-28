package gosocketio

import (
	"net"
	"fmt"
	"strings"
	"github.com/LucaBoss74/golang-socketio/transport"
	"strconv"
)

const (
	webSocketProtocol = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl       = "/socket.io/?EIO=3&transport=websocket"
)

/**
Socket.io client representation
*/
type Client struct {
	methods
	Channel
}

/**
Get ws/wss url by host and port - wrapper for GetUrlWithPath
 */
func GetUrl(host string, port int, secure bool) string {
	return GetUrlWithPath(host, port, secure, "")
}

/**
Get ws/wss url by host, port and path
*/
func GetUrlWithPath(host string, port int, secure bool, path string) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	if path != "" {
		workString := strings.ReplaceAll(socketioUrl, "/socket.io/", fmt.Sprintf("%s/", path))
		return prefix + net.JoinHostPort(host, strconv.Itoa(port)) + workString
	}
	return prefix + net.JoinHostPort(host, strconv.Itoa(port)) + socketioUrl
}

/**
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(url string, tr transport.Transport) (*Client, error) {
	c := &Client{}
	c.initChannel()
	c.initMethods()

	var err error
	c.conn, err = tr.Connect(url)
	if err != nil {
		return nil, err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	return c, nil
}

/**
Close client connection
*/
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
