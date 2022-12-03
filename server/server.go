package server

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	url              url.URL
	upgrader         websocket.Upgrader
	newClientHandler func(client *Client)
}

func NewServer(url url.URL) *Server {
	return &Server{
		url:      url,
		upgrader: websocket.Upgrader{},
	}
}

func (server *Server) Listen(onNewClient func(client *Client)) error {
	log.Debug("Commvent server listening...")

	server.newClientHandler = onNewClient // Catch handler
	return http.ListenAndServe(server.url.Host, http.HandlerFunc(server.onNewConnection))
}

func (server *Server) onNewConnection(writer http.ResponseWriter, request *http.Request) {
	ws, err := server.upgrader.Upgrade(writer, request, nil)

	log.Debugf("New client %s trying to connect", request.RequestURI)

	if err == nil {
		log.Debugf("New client %s connection upgraded to ws", request.RequestURI)

		client := NewClient(ws)
		server.newClientHandler(client)
	} else {
		log.Errorf("Error while trying to upgrade connection from client %s: %s", request.RequestURI, err.Error())
	}
}
