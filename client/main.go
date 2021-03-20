// client.go
package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hias222/soundcontrol/client/util"
	"github.com/hias222/soundcontrol/client/webmodels"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

var (
	gitCommit  string
	versionTag string
	buildType  string

	verbose bool
)

type Soundcontrol struct {
	socket      *webmodels.MessageSocket
	stopChannel chan bool
}

func NewSoundcontrol(verbose bool) (*Soundcontrol, error) {
	newSocket, err := webmodels.NewWebsocket()

	s := &Soundcontrol{
		socket:      newSocket,
		stopChannel: make(chan bool),
	}

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
		// return nil, fmt.Errorf("create new SerialIO: %w", err)
	}

	return s, nil
}

func (s *Soundcontrol) Initialize() error {
	log.Println("Initializing")

	/*
		// load the config for the first time
		if err := d.config.Load(); err != nil {
			d.logger.Errorw("Failed to load config during initialization", "error", err)
			return fmt.Errorf("load config during init: %w", err)
		}

		// initialize the session map
		if err := d.sessions.initialize(); err != nil {
			d.logger.Errorw("Failed to initialize session map", "error", err)
			return fmt.Errorf("init session map: %w", err)
		}
	*/

	s.setupInterruptHandler()
	s.run()

	return nil
}

func (s *Soundcontrol) setupInterruptHandler() {
	interruptChannel := util.SetupCloseHandler()

	go func() {
		signal := <-interruptChannel
		log.Println("Interrupted", "signal", signal)
		s.signalStop()
	}()
}

func (s *Soundcontrol) run() {
	log.Println("Run loop starting")

	// connect to the arduino for the first time
	go func() {
		if err := s.socket.Start(); err != nil {
			log.Fatal("Failed to start first-time socket connection", "error", err)

		}
	}()

	// wait until stopped (gracefully)
	<-s.stopChannel
	log.Println("Stop channel signaled, terminating")

	if err := s.stop(); err != nil {
		log.Fatal("Failed to stop sound ", "error", err)
		os.Exit(1)
	} else {
		// exit with 0
		os.Exit(0)
	}
}

func (s *Soundcontrol) stop() error {
	log.Println("Stopping")

	s.socket.Stop()

	return nil
}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}

func (s *Soundcontrol) signalStop() {
	log.Println("Signalling stop channel")
	s.stopChannel <- true
}

func main() {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	verbose = true

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	// create the souncontrol instance
	s, err := NewSoundcontrol(verbose)

	if err != nil {
		log.Fatal("Failed to create Sound Control object", "error", err)
	}

	log.Println(s)

	// onwards, to glory
	if err = s.Initialize(); err != nil {
		log.Fatal("Failed to initialize sound ", "error", err)
	}

	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		/*
			case <-time.After(time.Duration(1) * time.Millisecond * 1000):
				// Send an echo packet every second
				err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from GolangDocs!"))
				if err != nil {
					log.Println("Error during writing to websocket:", err)
					return
				}
		*/
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case <-done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}

	}
}
