package tftp

import (
	"log"
	"net"
)

func (h *Server) TftpSocketProcess(packet []byte, addr net.UDPAddr) {
	rrq, err := parseRRQ(packet)
	if err != nil {
		log.Printf("Rejected packet from %s: %v", addr.String(), err)
		return
	}

	log.Printf("Accepted RRQ from %s: filename='%s'", addr.String(), rrq.Filename)
	err = replyPort(addr, rrq)
	if err != nil {
		log.Printf("Error sending response to %s: %v", addr.String(), err)
		return
	}
}
