package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type ChatServer struct {
	mu      sync.Mutex
	clients map[int]chan string
	nextID  int
}

type JoinArgs struct{}
type JoinReply struct {
	ID int
}

type MessageArgs struct {
	ID      int
	Message string
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[int]chan string),
		nextID:  1,
	}
}

// Client joins and gets assigned an ID + channel
func (cs *ChatServer) Join(_ JoinArgs, reply *JoinReply) error {
	cs.mu.Lock()
	id := cs.nextID
	cs.nextID++
	ch := make(chan string, 10)
	cs.clients[id] = ch
	cs.mu.Unlock()

	*reply = JoinReply{ID: id}

	// Broadcast join message
	cs.broadcast(id, fmt.Sprintf("User %d joined", id))
	return nil
}

// Client sends a message
func (cs *ChatServer) SendMessage(args MessageArgs, _ *struct{}) error {
	cs.broadcast(args.ID, fmt.Sprintf("User %d: %s", args.ID, args.Message))
	return nil
}

// Broadcast to all except sender
func (cs *ChatServer) broadcast(senderID int, msg string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for id, ch := range cs.clients {
		if id != senderID {
			select {
			case ch <- msg:
			default:
				// drop if channel full
			}
		}
	}
}

// Stream messages to client
func (cs *ChatServer) Stream(id int, reply *[]string) error {
	cs.mu.Lock()
	ch, ok := cs.clients[id]
	cs.mu.Unlock()
	if !ok {
		return fmt.Errorf("client %d not found", id)
	}

	// Collect all messages currently in channel
	var msgs []string
	for {
		select {
		case m := <-ch:
			msgs = append(msgs, m)
		default:
			*reply = msgs
			return nil
		}
	}
}

func main() {
	server := NewChatServer()
	rpc.Register(server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Server running on port 1234...")
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}
