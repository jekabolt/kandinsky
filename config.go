package kandinsky

import (
	"github.com/jekabolt/ARGraffti-back/store"
)

// Configuration is a struct with all service options
type Configuration struct {
	Database    store.Conf
	RestAddress string
	WsAddress   string
}
