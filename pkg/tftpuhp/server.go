package tftp

import (
	"context"
	"log"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"
)

const udpPacketSize = 1500

type Server struct {
	tftp_socket *net.UDPConn
}

func (h *Server) Listen(context context.Context, systemd bool) error {
	if systemd {
		// systemd starts passing FDs at 3
		fd := uintptr(3)
		file := os.NewFile(fd, "systemd-udp69")
		if file == nil {
			log.Fatalf("failed to create os.File from fd %d", fd)
		}

		pc, err := net.FilePacketConn(file)
		if err != nil {
			log.Fatalf("failed to convert FD to PacketConn: %v", err)
		}

		conn, ok := pc.(*net.UDPConn)
		if !ok {
			log.Fatalf("expected *net.UDPConn but got %T", pc)
		}

		h.tftp_socket = conn
	} else {
		socket, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(netip.MustParseAddrPort("0.0.0.0:69")))
		if err != nil {
			log.Fatalf("failed to listen on tftp port: %v", err)
		}

		h.tftp_socket = socket
	}

	var serverShutdown sync.WaitGroup

	go func() {
		<-context.Done()
        h.tftp_socket.Close()
	}()

	serverShutdown.Add(1)
	go func() {
		buffer := make([]byte, udpPacketSize)
		for {
			n, addr, err := h.tftp_socket.ReadFromUDP(buffer)
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					log.Print("Closing UDP socket.")
				} else {
					log.Printf("udp read error: %v", err)
				}
				serverShutdown.Done()
				return
			}
			packet := make([]byte, n)
			copy(packet, buffer[0:n])
			go h.TftpSocketProcess(packet, *addr)
		}
	}()

	serverShutdown.Wait()
	return nil
}
