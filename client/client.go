package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/GiulianoDecesares/commvent/primitives"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	maxMessageSize   = 1024 // Maximum message size allowed from peer.
	eventsBufferSize = 256

	pongWait   = 30 * time.Second    // Time allowed to read the next pong message from the peer
	pingPeriod = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait
)

type Client struct {
	socket *websocket.Conn

	eventsBuffer   chan primitives.Event
	commandHandler func(commmand *primitives.Command)
}

func NewClient(commandHandler func(command *primitives.Command)) IClient {
	client := &Client{
		eventsBuffer:   make(chan primitives.Event, eventsBufferSize),
		commandHandler: commandHandler,
	}

	return client
}

func (client *Client) Begin(url url.URL) error {
	var result error = nil

	url.Scheme = "ws"

	if client.socket, result = client.connect(url); result == nil {
		client.info(fmt.Sprintf("Connection to %s established", url.Host))

		client.socket.SetPingHandler(func(appData string) error {
			client.trace("Ping received. Sending pong")
			return client.socket.WriteMessage(websocket.PongMessage, nil)
		})

		go client.receive()
		go client.send()
	} else {
		client.error(fmt.Sprintf("Error trying to connect to %s: %s", url.Host, result.Error()))
	}

	return result
}

func (client *Client) Stop() error {
	client.info("Closing")

	close(client.eventsBuffer)
	return client.socket.Close()
}

func (client *Client) SendEvent(event *primitives.Event) {
	if event != nil {
		client.eventsBuffer <- *event
	}
}

func (client *Client) receive() {
	defer client.socket.Close()

	client.socket.SetReadLimit(maxMessageSize)

	for {
		command := &primitives.Command{}

		if err := client.socket.ReadJSON(&command); err == nil {
			if client.commandHandler != nil {
				client.debug(fmt.Sprintf("Received %s", command.ToString()))
				client.commandHandler(command)
			} else {
				client.error(fmt.Sprintf("Null command handler while receiving event %s", command.ToString()))
			}
		} else {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.error(fmt.Sprintf("Unexpected error while receiving: %v", err))
			}

			break
		}
	}
}

func (client *Client) send() {
	for {
		event, ok := <-client.eventsBuffer

		if !ok { // Check if closed channel
			client.trace("Events buffer closed")
			client.socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := client.socket.WriteJSON(event); err != nil {
			client.error(fmt.Sprintf("Error while writing JSON: %v", err))
			return
		} else {
			client.debug(fmt.Sprintf("Sent %s", event.ToString()))
		}
	}
}

func (client *Client) connect(url url.URL) (*websocket.Conn, error) {
	socket, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	return socket, err
}

func (client *Client) getLocalAddress() string {
	return client.socket.LocalAddr().String()
}

func (client *Client) trace(message string) {
	log.Tracef("[Client %s] %s", client.getLocalAddress(), message)
}

func (client *Client) debug(message string) {
	log.Debugf("[Client %s] %s", client.getLocalAddress(), message)
}

func (client *Client) info(message string) {
	log.Infof("[Client %s] %s", client.getLocalAddress(), message)
}

func (client *Client) error(message string) {
	log.Errorf("[Client %s] %s", client.getLocalAddress(), message)
}
