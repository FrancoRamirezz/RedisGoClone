package redisclone

import (
	"bufio"
	"fmt"
	"io"

	"net"
	"sync"
)

// tcp server architecture will contain important components
type TCPServer struct {
	listner net.Listener      // Handles any of the incoming tcp connections. We are mirroing redis tcp connection
	clients map[net.Conn]bool //track the active client connections for connection managment
	mu      sync.Mutex        // ensures thread safe operations
	que     chan Event        // server as our event queue for proccessig. since were using mac think of this before we use keque
}

type Event struct {
	conn     net.Conn   // shows the client connection that were getting
	command  Command    // the redis command  like get, set, ping
	response chan error // response enables asychronous comm between goroutine
}

// command per each client request
type Command struct {
	cmd  string   // for reach command
	args []string // [Set, Get]
}
type Config struct {
	IP   net.IP
	Port int
}

// create the tcp server
func Server(c Config, address string) (*TCPServer, error) {
	if c.IP == nil {
		c.IP = net.ParseIP("0.0.0.0")
	}
	if c.Port == 0 {
		c.Port = 6379
	}
	// the conn will return the TCP pointer

	listen, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error occured", err)
	}
	// make sure to insitiate the tcp server
	return &TCPServer{
		listner: listen,
		clients: make(map[net.Conn]bool),
		que:     make(chan Event, 1000),
	}, nil

}

// event loop focuses how the event que works
// check for the event que proccess them events sequentially, handling command execution, managing error responses
func (server *TCPServer) eventLoop() {
	for event := range server.que {
		server.proccessCommand(event)

	}
}

// create each proccess command for each event loop
func (server *TCPServer) proccessCommand(e Event) *Event {

	return &Event{}
}

// now that we need to lock the connection for the event que we lock the tcpserver structure
// then we add our connection to the map after to handle thread safety

func (s *TCPServer) HandleConnection(conn net.Conn) {
	// make sure to lock thread safety
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
	}()
	// create the infinite loop to iterate throught
	re := bufio.NewReader(conn)
	for {
		// now reading the Redis command
		cmds, err := readCommands(re)
		if err != nil && err != io.EOF {
			fmt.Print("Error reading command", err)
		}
		// we will make a buffer to handle any of the
		res := make(chan error)
		// make the chan error for this
		s.que <- Event{
			conn:     conn,
			command:  cmds,
			response: res,
		}
		if err := <-res; err != nil {
			fmt.Print("Error Processing Comand", err)
			return
		}

	}

}

// some request could take longer than most and that it could block the other request during I/O operations
// ensures the redis commands are processed
func (server *TCPServer) Start() error {
	// call the go routine on the event loop
	go server.eventLoop()
	for {
		conn, err := server.listner.Accept()
		if err != nil {
			return err
		}
		// call a go routines for each conenction
		go server.HandleConnection(conn)
	}
}

// overview for future reference
// the server starts and initializes the event loop, then the client connects via tcp
//our connection handler creates a new event, que in the event loop, process them commands sequentially
//client end and handles any error in channels
