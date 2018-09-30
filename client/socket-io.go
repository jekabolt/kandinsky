package client

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"sync"

	"github.com/graarh/golang-socketio"

	"github.com/graarh/golang-socketio/transport"
	"github.com/jekabolt/kandinsky/store"
	"github.com/jekabolt/slf"
	_ "github.com/jekabolt/slflog"
	"github.com/satori/go.uuid"
)

func (c *SocketIOConnectedPool) SetSocketIOHandlers(address string) error {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	pool, err := InitConnectedPool(server, address)
	if err != nil {
		return fmt.Errorf("connection pool initialization: %s", err.Error())
	}
	c.Server = server
	c.closeChByConnID = &sync.Map{}
	c.users = &sync.Map{}
	c.log = slf.WithContext("socket-pool")

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		userid, err := uuid.NewV4()
		if err != nil {
			pool.log.Errorf("uuid.NewV4: %s", err.Error())
			return
		}
		user, err := newSocketIOUser(c, userid.String())
		if err != nil {
			pool.log.Errorf("get socketio headers: %s", err.Error())
			return
		}

		pool.users.LoadOrStore(c.Id(), user)
		pool.closeChByConnID.Store(c.Id(), user.closeCh)
		pool.log.Infof("Conncted %v", c.Id())
	})

	// Upload image
	server.On(uploadImage, func(c *gosocketio.Channel, userImg store.Image) {
		img, err := png.Decode(bytes.NewReader(userImg.Data))
		// img, str, err := image.Decode()
		f, err := os.Create(userImg.Filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		png.Encode(f, img)
		//save the imgByte to file
		// out, err := os.Create("./QRImg.png")
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		pool.log.Infof("Disconnected %s", c.Id())
		pool.users.Delete(c.Id())
		pool.closeChByConnID.Delete(c.Id())
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	pool.log.Infof("Starting socketIO server on %s address", address)
	go func() {
		pool.log.Panicf("%s", http.ListenAndServe(address, serveMux))
	}()
	return nil
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

func InitConnectedPool(server *gosocketio.Server, address string) (*SocketIOConnectedPool, error) {
	pool := &SocketIOConnectedPool{
		users:           &sync.Map{},
		closeChByConnID: &sync.Map{},
		log:             slf.WithContext("socket-pool"),
	}
	pool.log.Info("InitConnectedPool")

	return pool, nil
}

type SocketIOConnectedPool struct {
	users *sync.Map //  conn.Id()  to socket user

	closeChByConnID *sync.Map
	// closeChByConnID map[string]chan string   // when connection was finished, send close signal to his goroutine

	Server *gosocketio.Server
	log    slf.StructuredLogger
}

type SocketIOUser struct {
	userID  string
	closeCh *chan struct{}
	log     slf.StructuredLogger
	conn    *gosocketio.Channel
}
