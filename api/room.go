package api

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type room struct {
	clients map[*client]bool

	join chan *client

	leave chan *client

	forward chan []byte

	server *Server
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) setServer(s *Server) {
	r.server = s
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			log.Printf("Client joined. Total clients: %d", len(r.clients))
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
			log.Printf("Client left. Total clients: %d", len(r.clients))
		case msg := <-r.forward:
			log.Printf("Messagge received: %s", string(msg))
			for client := range r.clients {
				select {
				case client.receive <- msg:
					log.Printf("Message sent")
				default:
					log.Printf("Cannot send the message, removing the client")
					delete(r.clients, client)
					close(client.receive)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (r *room) ServerHTTP(w http.ResponseWriter, req *http.Request) {

	username := req.URL.Query().Get("username")
	groupId := req.URL.Query().Get("group_id")

	roomIdConverted, _ := strconv.ParseInt(groupId, 10, 64)

	if username == "" {
		log.Println("Username not given")
		http.Error(w, "Username is requested", http.StatusBadRequest)
		return
	}

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Error during the upgrade:", err)
		return
	}
	log.Println("Websocket established")

	client := &client{
		socket:   socket,
		receive:  make(chan []byte, messageBufferSize),
		room:     r,
		username: username,
		roomId:   roomIdConverted,
		server:   r.server,
	}
	r.join <- client
	log.Println("Client added join")
	defer func() {
		r.leave <- client
		log.Println("Client added leave")
	}()
	go client.read()
	client.write()
}
