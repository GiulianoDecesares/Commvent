package server

import (
	"fmt"
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
		url: url,
		upgrader: websocket.Upgrader{
			CheckOrigin: checkOrigin,
		},
	}
}

func (server *Server) Listen(onNewClient func(client *Client)) error {
	if onNewClient == nil {
		server.warning("Null client handler")
	}

	server.info("Listening...")

	server.newClientHandler = onNewClient
	return http.ListenAndServe(server.url.Host, http.HandlerFunc(server.onNewConnection))
}

func (server *Server) onNewConnection(writer http.ResponseWriter, request *http.Request) {
	ws, err := server.upgrader.Upgrade(writer, request, nil)

	server.debug(fmt.Sprintf("Connection attempt from: %s", request.RemoteAddr))

	if err == nil {
		server.info(fmt.Sprintf("Client %s connected", request.RemoteAddr))

		if server.newClientHandler != nil {
			server.newClientHandler(NewClient(ws))
		} else {
			server.warning(fmt.Sprintf("No handler to manage %s connection", request.RemoteAddr))
		}
	} else {
		server.error(fmt.Sprintf("Error while trying to upgrade connection from client %s: %s", request.RemoteAddr, err.Error()))
	}
}

func (server *Server) debug(message string) {
	log.Debugf("[Commvent Server] %s", message)
}

func (server *Server) info(message string) {
	log.Infof("[Commvent Server] %s", message)
}

func (server *Server) warning(message string) {
	log.Warnf("[Commvent Server] %s", message)
}

func (server *Server) error(message string) {
	log.Errorf("[Commvent Server] %s", message)
}

func checkOrigin(request *http.Request) bool {
	return true // TODO :: Fix this possible security issue
}
