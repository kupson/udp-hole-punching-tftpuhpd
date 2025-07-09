package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/kupson/udp-hole-punching-tftpuhpd/pkg/tftpuhp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)

	listenPidStr := os.Getenv("LISTEN_PID")
	listenFdsStr := os.Getenv("LISTEN_FDS")

	if listenPidStr != "" && listenFdsStr != "" {
		pid, err := strconv.Atoi(listenPidStr)
		if err != nil || pid != os.Getpid() {
			log.Fatalf("Invalid or mismatched LISTEN_PID: got %v, expected %v", listenPidStr, os.Getpid())
		}

		numFds, err := strconv.Atoi(listenFdsStr)
		if err != nil || numFds != 1 {
			log.Fatalf("Invalid LISTEN_FDS: %v", listenFdsStr)
		}
	}

	udpServer := &tftp.Server{}

	var serverShutdown sync.WaitGroup
	ctx, cancelCtx := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Cancel context if interrupt signal received
	go func() {
		sig := <-sigChan
		log.Printf("Signal %s - gracefully shutting down.", sig)
		cancelCtx()
	}()

	// Udp server, context cancelled on error
	serverShutdown.Add(1)
	go func() {
		err := udpServer.Listen(ctx, (listenPidStr != "" && listenFdsStr != ""))
		if err != nil {
			log.Fatalf("error starting tftp server: %s", err)
		}
		cancelCtx()
		serverShutdown.Done()
	}()

	serverShutdown.Wait()
	log.Print("Shutdown done.")
}
