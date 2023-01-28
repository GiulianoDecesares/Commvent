package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config   ServerConfig
	upgrader websocket.Upgrader
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		config: config,
		upgrader: websocket.Upgrader{
			CheckOrigin: checkOrigin,
		},
	}
}

func (server *Server) HandleHttp(endpoint string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(endpoint, handler)
}

func (server *Server) Handle(endpoint string, handler func(client *Client)) {
	http.HandleFunc(endpoint, func(writer http.ResponseWriter, request *http.Request) {
		server.debug(fmt.Sprintf("Connection attempt from: %s", request.RemoteAddr))

		if ws, err := server.upgrader.Upgrade(writer, request, nil); err == nil {
			server.info(fmt.Sprintf("Client %s connected", request.RemoteAddr))

			if handler != nil {
				handler(NewClient(ws, server.config.ClientConfig))
			} else {
				server.warning(fmt.Sprintf("No handler to manage %s connection", request.RemoteAddr))
			}
		} else {
			server.error(fmt.Sprintf("Error while trying to upgrade connection from client %s: %s", request.RemoteAddr, err.Error()))
		}
	})
}

func (server *Server) Listen() error {
	server.info("Listening...")
	return http.ListenAndServe(server.config.Server.Host+":"+server.config.Server.Port, nil)
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
