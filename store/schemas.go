package store

const (
	// ws routes
	UploadImage = "upload:image"
)

type Image struct {
	Filename string
	Data     []byte
}
