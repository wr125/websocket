package main

import (
	"os"
	"os/signal"
	"syscall"
	"wr125/websocket/server"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := server.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}

}
