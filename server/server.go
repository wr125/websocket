package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	listener         net.Listener
	db               memoryDB
	connections      map[int]net.Conn
	connCloseTimeout time.Duration
	quit             chan struct{}
	exited           chan struct{}
}

func NewServer() *Server {
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("creating listener failed", err.Error())
	}
	srv := &Server{
		listener:         l,
		quit:             make(chan struct{}),
		exited:           make(chan struct{}),
		db:               newDB(),
		connections:      map[int]net.Conn{},
		connCloseTimeout: 10 * time.Second,
	}
	go srv.serve()
	return srv
}

func (srv *Server) serve() {
	//var handlers sync.WaitGroup
	var id int
	fmt.Println("listening for clients")
	for {
		select {
		case <-srv.quit:
			fmt.Println("shutting server down")
			err := srv.listener.Close()
			if err != nil {
				fmt.Println("unable to close listener", err.Error())
			}
			if len(srv.connections) > 0 {
				srv.warnConnections(srv.connCloseTimeout)
				<-time.After(srv.connCloseTimeout)
				srv.closeConnections()
			}
			//handlers.Wait()
			close(srv.exited)
			return
		default:
			tcpListener := srv.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				fmt.Println("listener deadline failed to be set correctly", err.Error())
			}

			conn, err := tcpListener.Accept()
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err != nil {
				fmt.Println("connection accept failed:", err.Error())
			}

			write(conn, "Welcome to in-memoryDB server")
			//handlers.Add(1)
			srv.connections[id] = conn
			go func(connID int) {
				fmt.Println("current client id:", connID, "joined")
				srv.handleConn(conn)
				//handlers.Done()
				delete(srv.connections, connID)
				fmt.Println("current client id:", connID, "left")
			}(id)
			id++
		}
	}
}

func write(c net.Conn, s string) {
	_, err := fmt.Fprintf(c, "%s\n-> ", s)
	if err != nil {
		log.Fatal(err)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		l := strings.ToLower(strings.TrimSpace(scanner.Text()))
		values := strings.Split(l, " ")

		switch {
		case len(values) == 3 && values[0] == "set":
			srv.db.set(values[1], values[2])
			write(conn, "OK")
		case len(values) == 2 && values[0] == "get":
			val, found := srv.db.get(values[1])
			if !found {
				write(conn, fmt.Sprintf("key %s not found", values[1]))
			} else {
				write(conn, val)
			}
		case len(values) == 2 && values[0] == "delete":
			srv.db.delete(values[1])
			write(conn, "OK")
		case len(values) == 1 && values[0] == "exit":
			if err := conn.Close(); err != nil {
				fmt.Println("could not close connection", err.Error())
			}
		default:
			write(conn, fmt.Sprintf("UNKNOWN: %s", l))
		}
	}
}

func (srv *Server) warnConnections(timeout time.Duration) {
	for _, conn := range srv.connections {
		write(conn, fmt.Sprintf("host wants to shut down the server in: %s", srv.connCloseTimeout.String()))
	}
}

func (srv *Server) closeConnections() {
	fmt.Println("closing all connections")
	for id, conn := range srv.connections {
		err := conn.Close()
		if err != nil {
			fmt.Println("could not close connection with id:", id, err.Error())
		}
	}
}

func (srv *Server) Stop() {
	fmt.Println("stopping the database server")
	close(srv.quit)
	<-srv.exited
	fmt.Println("saving in-memory records to file")
	srv.db.save()
	fmt.Println("database server successfully stopped")
}
