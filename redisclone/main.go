package redisclone

import (
	aofgo "Pokemonscraper/redisclone/AOF"
	Datastructure "Pokemonscraper/redisclone/DataStructure"
	"Pokemonscraper/redisclone/resp"

	"fmt"
	"io"
	"log"
	"net"
)

// TCP accept loop, one goroutine per client, the read→execute→respond cycle
//
//	type Server struct{
//		addr string
//		proto string
//	}
//
// will instiante the server here.. SO, we can create other files
//
//	func CreateServer(addr string, proto string)*Server{
//		return &Server{addr, proto}
//	}
const (
	AppendOnlyFile = "appendonly.aof"
)

// this follow a basic concurrent web server
func main() {
	//Step one activate the TTL in the background
	db := Datastructure.NewStore()
	fmt.Print("exipiration background", db)

	//Step two: append the log file
	aof, err := aofgo.NewAOFPath(AppendOnlyFile)
	if err != nil {
		log.Fatal("corrupt file", aof, err)
	}
	defer aof.Close()
	// listen takes two arguments or parameters a network string, tcp, udp, address string 6379
	ListenConn, err := net.Listen("tcp", "6379")
	if err != nil {
		log.Fatal(err)
	}
	defer ListenConn.Close()
	// we can use a panic panic(err)
	fmt.Print("Server started on port 6379")
	// here we close the listner once the function exits
	// here we will start the connection
	// this is the same as while True
	for {
		conn, err := ListenConn.Accept()
		if err != nil {
			continue
		}
		// here we will conncurenctly call
		// the connection
		go handleConnection(conn, db, aof)
	}

}

// then we call the concurrent function
// we will add the amount of bytes we will process
func handleConnection(conn net.Conn, db *Datastructure.Store, a *aofgo.AOF) {
	defer func() {
		log.Println("Closing connection", conn)
		conn.Close()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Print("recovering from error")
		}
	}()
	// the goal is to handle the conditions for the resp writer and reader
	// Note the resp writer is in its own struct just like the resp reader
	reader := resp.NewReader(conn)
	writer := resp.NewRespWriter(conn)
	for {
		// read all the commands form a single client connection
		// Note we will use a gouroutines to return the connect
		v, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("client disconnected", err, conn)
		}
		// the redis command line will get a response array
		// example: $5\r\nFranco\r\n
		var args []resp.RespValue
		if v.Type == resp.RespArray {
			args = v.Array
		} else {
			args = []resp.RespValue{v}
		}
		if len(args) == 0 {
			continue
		}

		// check for any issues and respond
		if err := writer.Write(resp.RespValue{}); err != nil {
			log.Printf("write error:%v", err)
			return
		}

	}

}
