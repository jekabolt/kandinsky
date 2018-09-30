package kandinsky

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jekabolt/kandinsky/client"
	"github.com/jekabolt/slf"
)

var (
	log = slf.WithContext("kandinsky")
)

// Multy is a main struct of service
type Kandinsky struct {
	config *Configuration
	// userStore  store.UserStore
	// restClient *client.RestClient
	WsServer *client.SocketIOConnectedPool
	route    *gin.Engine
}

// Init initializes Multy instance
func Init(conf *Configuration) (*Kandinsky, error) {
	kandinsky := &Kandinsky{
		config:   conf,
		WsServer: &client.SocketIOConnectedPool{},
	}

	err := kandinsky.WsServer.SetSocketIOHandlers(conf.WsAddress)
	if err != nil {
		return nil, fmt.Errorf("SetSocketIOHandlers: %s on port %s", err.Error(), conf.WsAddress)
	}
	fmt.Println("kandinsky.WsServer ", kandinsky.WsServer)

	return kandinsky, nil
}

// Run runs service
func (kandinsky *Kandinsky) Run() error {
	log.Info("Running server")
	kandinsky.route.Run(kandinsky.config.RestAddress)
	return nil
}
