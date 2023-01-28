package server

import (
	"fmt"
	"time"

	"github.com/GiulianoDecesares/commvent/primitives"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type Client struct {
	socket *websocket.Conn

	open bool

	writeWait time.Duration
	pongWait  time.Duration

	pingTicker *time.Ticker

	commandsBuffer chan primitives.Message
	stop           chan struct{}

	eventHandler        func(event *primitives.Message)
	disconnectedHandler func(disconnectionState error)
}

func NewClient(socket *websocket.Conn, config ClientConfig) *Client {
	client := &Client{
		socket:         socket,
		open:           true,
		writeWait:      config.WriteWait,
		pongWait:       config.PongWait,
		pingTicker:     time.NewTicker(config.PingPeriod),
		commandsBuffer: make(chan primitives.Message, config.BufferSize),
		stop:           make(chan struct{}, 1),
		eventHandler:   nil,
	}

	client.socket.SetReadLimit(config.MaxMessageSize)

	client.socket.SetPongHandler(func(string) error {
		client.trace("Pong received")
		client.updateReadDeadline() // Allow read while receiving pong from peer
		return nil
	})

	go client.ping()

	go client.receive()
	go client.send()

	return client
}

func (client *Client) HandleEvents(handler func(event *primitives.Message)) {
	if handler != nil {
		client.eventHandler = handler
	}
}

func (client *Client) HandleDisconnect(handler func(disconnectionState error)) {
	if handler != nil {
		client.disconnectedHandler = handler
	}
}

func (client *Client) SendCommand(command *primitives.Message) {
	if command != nil && client.open {
		client.commandsBuffer <- *command
	} else {
		client.error(fmt.Sprintf("Error while sending %s", command.ToString()))
	}
}

func (client *Client) close() {
	if client.open {
		client.debug("Closing")

		client.open = false

		client.socket.WriteMessage(websocket.CloseMessage, []byte{})

		client.pingTicker.Stop()

		client.trace("Stoping all coroutines safely")
		client.stop <- struct{}{}

		close(client.commandsBuffer)
		close(client.stop)

		result := client.socket.Close()

		if client.disconnectedHandler != nil {
			client.disconnectedHandler(result)
		} else {
			client.error("Null disconnected handler")
		}
	}
}

func (client *Client) receive() {
	client.debug("Starting receiving coroutine")
	client.updateReadDeadline()

	for {
		select {
		case <-client.stop:
			client.trace("Receiving coroutine finished by stop mechanism")
			return

		default:
			event := &primitives.Message{}

			if err := client.socket.ReadJSON(&event); err == nil {
				if client.eventHandler != nil {
					client.debug(fmt.Sprintf("Received %s", event.ToString()))
					client.eventHandler(event)
				} else {
					client.error(fmt.Sprintf("Null event handler while receiving event %s", event.ToString()))
				}
			} else {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					client.error(fmt.Sprintf("Unexpected error while receiving: %v", err))
				}

				client.close()
				client.trace("Receiving coroutine finished by close")
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
			client.trace("Sending coroutine finished by stop mechanism")
			return

		case command := <-client.commandsBuffer:
			client.updateWriteDeadline()

			if err := client.socket.WriteJSON(command); err != nil {
				client.error(fmt.Sprintf("Error while writing JSON: %v", err))
				return
			} else {
				client.debug(fmt.Sprintf("Sent %s", command.ToString()))
			}
		}
	}
}

func (client *Client) ping() {
	client.debug("Starting ping coroutine")

	for {
		select {
		case <-client.stop:
			client.trace("Ping coroutine finished by stop mechanism")
			return

		case <-client.pingTicker.C:
			client.trace("Sending ping")
			client.updateWriteDeadline()

			if err := client.sendPing(); err != nil {
				client.error("Error while sending ping")
				return
			}
		}
	}
}

func (client *Client) getRemoteAddress() string {
	return client.socket.RemoteAddr().String()
}

func (client *Client) getLocalAddress() string {
	return client.socket.LocalAddr().String()
}

func (client *Client) updateWriteDeadline() {
	client.socket.SetWriteDeadline(time.Now().Add(client.writeWait))
}

func (client *Client) updateReadDeadline() {
	client.socket.SetReadDeadline(time.Now().Add(client.pongWait))
}

func (client *Client) sendPing() error {
	return client.socket.WriteMessage(websocket.PingMessage, nil)
}

func (client *Client) trace(message string) {
	log.Tracef("[Client %s] %s", client.getRemoteAddress(), message)
}

func (client *Client) debug(message string) {
	log.Debugf("[Client %s] %s", client.getRemoteAddress(), message)
}

func (client *Client) info(message string) {
	log.Infof("[Client %s] %s", client.getRemoteAddress(), message)
}

func (client *Client) error(message string) {
	log.Errorf("[Client %s] %s", client.getRemoteAddress(), message)
}
