package client

import (
	"fmt"
	"sync"

	"github.com/graarh/golang-socketio"

	"github.com/jekabolt/slf"
	_ "github.com/jekabolt/slflog"
)

type SocketIOConnectedPool struct {
	users           *sync.Map //  conn.Id()  to socket user
	closeChByConnID *sync.Map // when connection was finished, send close signal to his goroutine
	Server          *gosocketio.Server
	log             slf.StructuredLogger
}

type SocketIOUser struct {
	userID  string
	closeCh *chan struct{}
	log     slf.StructuredLogger
	conn    *gosocketio.Channel
}

func InitConnectedPool(server *gosocketio.Server, address string) (*SocketIOConnectedPool, error) {
	pool := &SocketIOConnectedPool{
		users:           &sync.Map{},
		closeChByConnID: &sync.Map{},
		log:             slf.WithContext("socket-pool"),
	}
	pool.log.Info("InitConnectedPool")

	return pool, nil
}

func newSocketIOUser(c *gosocketio.Channel, userid string) (*SocketIOUser, error) {
	closeCh := make(chan struct{}, 0)
	return &SocketIOUser{
		userID:  userid,
		log:     slf.WithContext(fmt.Sprintf("socket-conn %v", userid)),
		closeCh: &closeCh,
		conn:    c,
	}, nil
}
