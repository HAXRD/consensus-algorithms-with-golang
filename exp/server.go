package exp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var upgrader = websocket.Upgrader{}
var mutex = &sync.Mutex{}

type Server struct {
	Host   string                     `json:"host"`
	Port   int                        `json:"port"`
	WsPort int                        `json:"ws_port"`
	Peers  map[string]*websocket.Conn `json:"peers"`
}

func NewServer(host string, wsPort int) *Server {
	return &Server{
		Host:   host,
		Port:   wsPort + 10000,
		WsPort: wsPort,
		Peers:  make(map[string]*websocket.Conn),
	}
}

func (s *Server) httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	out, err1 := json.Marshal(s)
	if err1 != nil {
		log.Println("json marshal error:", err1)
	}
	_, err := fmt.Fprintf(w, string(out))
	if err != nil {
		log.Println("httpHandler:", err)
	}
}

func (s *Server) msgHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed: ", err)
		return
	}
	defer conn.Close()

	for {
		// read a message from the websocket connection
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message failed: ", err)
			break
		}
		log.Printf("recv[%v]: %s\n", mt, message)
		// TODO: do anything that switch according to message types
	}
}

func (s *Server) connectPeers(addresses []string) {
	for _, address := range addresses {
		url := fmt.Sprintf("ws://%s/ws", address)
		log.Printf("Trying connecting to peer %s\n", url)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Error connecting to peer %s: %v\n", url, err)
			continue
		}
		log.Printf("Connected to peer %s\n", url)
		//defer conn.Close()

		mutex.Lock()
		s.Peers[address] = conn
		mutex.Unlock()
	}
}

func (s *Server) broadcastMessage(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for address, conn := range s.Peers {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error broadcasting message to peer %s: %v", address, err)
			conn.Close()
			delete(s.Peers, address)
		}
	}
}

func (s *Server) Listen(peerAddresses []string) {
	// start a http server
	http.HandleFunc("/server-info", s.httpHandler)
	go func() {
		url := s.Host + ":" + strconv.Itoa(s.Port)
		log.Printf("Http server started at %s\n", url)
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatalf("Http server error: %v\n", err)
		}
	}()

	// start websocket server
	http.HandleFunc("/ws", s.msgHandler)
	go func() {
		url := s.Host + ":" + strconv.Itoa(s.WsPort)
		log.Printf("WebSocket server started at %s\n", url)
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatalf("WebSocket server error: %v\n", err)
		}
	}()

	// connect to peers
	if len(peerAddresses) > 0 {
		s.connectPeers(peerAddresses)
	}
}
