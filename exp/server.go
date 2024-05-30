package exp

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	Host string
	Peer map[string]*websocket.Conn
}

func NewServer(host string) *Server {
	return &Server{
		Host: host,
		Peer: make(map[string]*websocket.Conn),
	}
}

func (s *Server) msgHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed: ", err)
		return
	}
	defer c.Close()

	for {
		// read a message from the wesocket connection
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read message failed: ", err)
			break
		}
		log.Printf("recv[%s]: %s", mt, message)
		// TODO: do anything that switch according to message types
	}
}

func (s *Server) Listen() {
	http.HandleFunc("/", s.msgHandler)
	err := http.ListenAndServe(s.Host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("Listening on ws://" + s.Host)
}
