package server

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second // Time allowed to write a message to the peer.

	commandsBufferSize = 256

	pongWait   = 5 * time.Second     // Time allowed to read the next pong message from the peer.
	pingPeriod = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.

	maxMessageSize = 1024 // Maximum message size allowed from peer.
)

type Client struct {
	socket *websocket.Conn

	open bool

	pingTicker *time.Ticker

	commandsBuffer chan Command
	stop           chan struct{}

	eventHandler        func(event *Event)
	disconnectedHandler func(disconnectionState error)
}

func NewClient(socket *websocket.Conn) *Client {
	client := &Client{
		socket:         socket,
		open:           true,
		pingTicker:     time.NewTicker(pingPeriod),
		commandsBuffer: make(chan Command, commandsBufferSize),
		stop:           make(chan struct{}, 1),
		eventHandler:   nil,
	}

	client.socket.SetReadLimit(maxMessageSize)

	client.socket.SetPongHandler(func(string) error {
		client.debug("Pong received")
		client.updateReadDeadline() // Read while receiving pong from peer
		return nil
	})

	go client.ping()

	go client.receive()
	go client.send()

	return client
}

func (client *Client) HandleEvents(handler func(event *Event)) {
	if handler != nil {
		client.eventHandler = handler
	}
}

func (client *Client) HandleDisconnect(handler func(disconnectionState error)) {
	if handler != nil {
		client.disconnectedHandler = handler
	}
}

func (client *Client) SendCommand(command *Command) {
	if command != nil && client.open {
		client.commandsBuffer <- *command
	}
}

func (client *Client) close() {
	if client.open {
		client.debug("Closing")

		client.open = false

		client.socket.WriteMessage(websocket.CloseMessage, []byte{})

		client.pingTicker.Stop()
		client.stop <- struct{}{} // Stop all coroutines safely

		close(client.commandsBuffer)
		close(client.stop)

		result := client.socket.Close()

		if client.disconnectedHandler != nil {
			client.disconnectedHandler(result)
		}
	}
}

func (client *Client) receive() {
	client.debug("Starting receiving coroutine")
	client.updateReadDeadline()

	for {
		select {
		case <-client.stop:
			client.debug("Receiving coroutine finished by stop mechanism")
			return

		default:
			event := &Event{}

			if err := client.socket.ReadJSON(&event); err == nil {
				client.eventHandler(event)
			} else {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					client.error(fmt.Sprintf("Unexpected error while receiving: %v", err))
				}

				client.close()
				client.debug("Receiving coroutine finished by close")
				return
			}
		}
	}
}

func (client *Client) send() {
	client.debug("Starting sending coroutine")

	for {
		select {
		case <-client.stop:
			client.debug("Sending coroutine finished by stop mechanism")
			return

		case command := <-client.commandsBuffer:
			client.updateWriteDeadline()

			if err := client.socket.WriteJSON(command); err != nil {
				client.error(fmt.Sprintf("Error while writing JSON: %v", err))
				return
			}
		}
	}
}

func (client *Client) ping() {
	client.debug("Starting ping coroutine")

	for {
		select {
		case <-client.stop:
			client.debug("Ping coroutine finished by stop mechanism")
			return

		case <-client.pingTicker.C:
			client.debug("Sending ping")
			client.updateWriteDeadline()

			if err := client.sendPing(); err != nil {
				client.error("Error while sending ping")
				return
			}
		}
	}
}

func (client *Client) updateWriteDeadline() {
	client.socket.SetWriteDeadline(time.Now().Add(writeWait))
}

func (client *Client) updateReadDeadline() {
	client.socket.SetReadDeadline(time.Now().Add(pongWait))
}

func (client *Client) sendPing() error {
	return client.socket.WriteMessage(websocket.PingMessage, nil)
}

func (client *Client) debug(message string) {
	log.Debugf("[Client %s] %s", client.socket.LocalAddr().String(), message)
}

func (client *Client) error(message string) {
	log.Errorf("[Client %s] %s", client.socket.LocalAddr().String(), message)
}
