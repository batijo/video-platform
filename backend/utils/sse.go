package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Broker struct {
	Clients        map[chan string]uint
	NewClients     chan Client
	DefunctClients chan chan string
	Messages       chan Message
}

type Client struct {
	ClientChannel chan string
	ClientId      uint
}

type Message struct {
	Msg      string
	ClientId uint
}

var (
	B       *Broker
	started = false
)

func NewSseServer() {
	B = &Broker{
		Clients:        make(map[chan string]uint),
		NewClients:     make(chan Client),
		DefunctClients: make(chan (chan string)),
		Messages:       make(chan Message),
	}
	B.start()
}

func (b *Broker) start() {
	started = true
	go func() {
		for {
			select {
			case newCl := <-b.NewClients:
				// start sending client messages
				b.Clients[newCl.ClientChannel] = newCl.ClientId

			case defCl := <-b.DefunctClients:
				// stop sending client messages
				delete(b.Clients, defCl)
				close(defCl)

			case msg := <-b.Messages:
				// there is a new message to send to all clients
				for clChan, id := range b.Clients {
					if msg.ClientId == id {
						clChan <- msg.Msg
					} else if msg.ClientId == 0 {
						clChan <- msg.Msg
					}
				}
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Channel for broker to send messages
	messageChan := make(chan string)

	// Add this client to the map of those that should receive updates
	token := mux.Vars(r)["token"]
	userID, _, err := GetUserIDFromToken(token)
	if err != nil {
		http.Error(w, fmt.Sprint("Error geting user ID for SSE connection: ", err), http.StatusInternalServerError)
		return
	}
	client := Client{
		ClientChannel: messageChan,
		ClientId:      userID,
	}
	b.NewClients <- client

	// Listen to the closing of the http connection
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify

		// Remove this client from the map of attached clients when `EventHandler` exits
		b.DefunctClients <- messageChan
		//log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Don't close the connection, instead loop endlessly
	for {
		msg, open := <-messageChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)

		// Flush data as soon as it has been updated
		f.Flush()
	}
}

func UpdateAllUsersMessage(msg string) {
	if !started {
		log.Println("SSE server has not been started")
		return
	}
	message := Message{
		msg,
		0,
	}

	B.Messages <- message
}

func UpdateUserMessage(msg string, ClientID uint) {
	if !started {
		log.Println("SSE server has not been started")
		return
	}
	message := Message{
		fmt.Sprintf("<%v> %v", time.Now().Format("15:04:05"), msg),
		ClientID,
	}

	B.Messages <- message
}
