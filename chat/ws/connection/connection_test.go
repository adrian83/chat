package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestU(t *testing.T) {
	config := &websocket.Config{}
	rwc := &MyReadWriteCloser{
		read:  make([]string, 0),
		wrote: make([]string, 0),
	}

	_, err := websocket.NewClient(config, rwc)

	assert.NoError(t, err, "")
}

type MyReadWriteCloser struct {
	read  []string
	wrote []string
}

func (rwc *MyReadWriteCloser) Read(p []byte) (int, error) {
	rwc.read = append(rwc.read, string(p))
	return len(p), nil
}

func (rwc *MyReadWriteCloser) Write(p []byte) (int, error) {
	rwc.wrote = append(rwc.wrote, string(p))
	return len(p), nil
}

func (rwc *MyReadWriteCloser) Close() error {
	return nil
}
