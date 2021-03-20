package webmodels

import (
	"errors"
	"fmt"
	"log"
)

type MessageSocket struct {
	serverURL string

	stopChannel chan bool
	connected   bool

	lastKnownNumSliders        int
	currentSliderPercentValues []float32
}

func AllUsers() {
	fmt.Println("All Users")
}

func NewWebsocket() (*MessageSocket, error) {

	newSocket := &MessageSocket{
		serverURL:   "ws://localhost:8080" + "/ws",
		stopChannel: make(chan bool),
		connected:   true,
	}

	log.Println("Created socket instance")

	newSocket.start()

	return newSocket, nil
}

func (newSocket *MessageSocket) start() error {

	if newSocket.connected {
		log.Println("Already connected, can't start another without closing first")
		return errors.New("serial: connection already active")
	}

	return nil

}
