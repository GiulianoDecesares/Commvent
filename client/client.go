package client

import (
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

func NewClient(commandHandler func(command *primitives.Command)) *Client {
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
		client.socket.SetPingHandler(func(appData string) error {
			log.Debug("Ping received. Sending pong")
			return client.socket.WriteMessage(websocket.PongMessage, nil)
		})

		go client.receive()
		go client.send()
	}

	return result
}

func (client *Client) Stop() error {
	log.Debug("Stopping client")

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
			client.commandHandler(command)
		} else {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("Unexpected error while receiving: %s", err.Error())
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
			log.Debug("Closing event buffer")
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
