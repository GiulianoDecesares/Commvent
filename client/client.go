package client

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize   = 1024 // Maximum message size allowed from peer.
	eventsBufferSize = 256

	pongWait   = 30 * time.Second    // Time allowed to read the next pong message from the peer
	pingPeriod = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait
)

type Client struct {
	socket *websocket.Conn

	eventsBuffer   chan Event
	commandHandler func(commmand *Command)
}

func NewClient(commandHandler func(command *Command)) *Client {
	client := &Client{
		eventsBuffer:   make(chan Event, eventsBufferSize),
		commandHandler: commandHandler,
	}

	return client
}

func (client *Client) Begin(url url.URL) error {
	var result error = nil

	url.Scheme = "ws"

	if client.socket, result = client.connect(url); result == nil {
		client.socket.SetPingHandler(func(appData string) error {
			// fmt.Println("Ping received. Sending pong")
			return client.socket.WriteMessage(websocket.PongMessage, nil)
		})

		go client.receive()
		go client.send()
	}

	return result
}

func (client *Client) Stop() error {
	fmt.Println("Stopping client")

	close(client.eventsBuffer)
	return client.socket.Close()
}

func (client *Client) SendEvent(event *Event) {
	if event != nil {
		client.eventsBuffer <- *event
	}
}

func (client *Client) receive() {
	defer client.socket.Close()

	client.socket.SetReadLimit(maxMessageSize)

	for {
		command := &Command{}

		if err := client.socket.ReadJSON(&command); err == nil {
			client.commandHandler(command)
		} else {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}
	}
}

func (client *Client) send() {
	for {
		event, ok := <-client.eventsBuffer
		// client.socket.SetWriteDeadline(time.Now().Add(writeWait))

		if !ok { // Check if closed channel
			fmt.Println("Closing event buffer")
			client.socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		client.socket.WriteJSON(event)
	}
}

func (client *Client) connect(url url.URL) (*websocket.Conn, error) {
	socket, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	return socket, err
}
