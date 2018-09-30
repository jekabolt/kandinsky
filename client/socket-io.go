package client

import (
	"fmt"
	"net/http"

	"github.com/graarh/golang-socketio"

	"github.com/graarh/golang-socketio/transport"
	"github.com/jekabolt/kandinsky/store"
	_ "github.com/jekabolt/slflog"
)

func (c *SocketIOConnectedPool) SetSocketIOHandlers(address string) error {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	pool, err := InitConnectedPool(server, address)
	if err != nil {
		return fmt.Errorf("connection pool initialization: %s", err.Error())
	}
	c = pool
	c.Server = server

	server.On(gosocketio.OnConnection, c.onConnection)

	server.On(gosocketio.OnDisconnection, c.onDisconnect)

	// Upload and process image
	server.On(store.UploadImage, c.uploadImage)

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	pool.log.Infof("Starting socketIO server on %s address", address)
	go func() {
		pool.log.Panicf("%s", http.ListenAndServe(address, serveMux))
	}()
	return nil
}
