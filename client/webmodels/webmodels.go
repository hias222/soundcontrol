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
		connected:   false,
	}

	log.Println("Created socket instance")

	return newSocket, nil
}

func (newSocket *MessageSocket) Start() error {

	if newSocket.connected {
		log.Println("Already connected, can't start another without closing first")
		return errors.New("serial: connection already active")
	}

	log.Println("webmodel start ....")

	return nil

}

func (newSocket *MessageSocket) Stop() {
	if newSocket.connected {
		log.Println("Shutting down socket connection")
		newSocket.stopChannel <- true
	} else {
		log.Println("Not currently connected, nothing to stop")
	}
}
