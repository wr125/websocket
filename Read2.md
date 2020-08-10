package main

import (
	"os"
	"os/signal"
	"syscall"
	"wr125/websocket/server"
)

// won't work here, because of security reasons on the platform
// feel free to download the code and play with it on your own

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := server.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
}
