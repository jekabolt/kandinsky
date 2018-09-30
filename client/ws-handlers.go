package client

import (
	"bytes"
	"image/png"
	"os"

	"github.com/graarh/golang-socketio"
	"github.com/jekabolt/kandinsky/store"
	uuid "github.com/satori/go.uuid"
)

func (pool *SocketIOConnectedPool) onConnection(c *gosocketio.Channel) {
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
}

func (pool *SocketIOConnectedPool) onDisconnect(c *gosocketio.Channel) {
	pool.log.Infof("Disconnected %s", c.Id())
	pool.users.Delete(c.Id())
	pool.closeChByConnID.Delete(c.Id())
}

func (pool *SocketIOConnectedPool) uploadImage(c *gosocketio.Channel, userImg store.Image) {
	img, err := png.Decode(bytes.NewReader(userImg.Data))
	// img, str, err := image.Decode()
	f, err := os.Create(userImg.Filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
	//TODO: add image to queue for process

}
