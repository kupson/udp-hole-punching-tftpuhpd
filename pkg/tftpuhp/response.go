package tftp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	dataOpcode   = 3
	dataBlockNum = 1
)

func replyPort(addr net.UDPAddr, _ *RRQ) error {
	var response bytes.Buffer

	binary.Write(&response, binary.BigEndian, uint16(dataOpcode))
	binary.Write(&response, binary.BigEndian, uint16(dataBlockNum))
	response.WriteString(fmt.Sprintf("%s:%d\n", addr.IP.String(), addr.Port))

	conn, err := net.ListenUDP("udp", nil) // bind to any available port
	if err != nil {
		return fmt.Errorf("failed to create response socket: %v", err)
	}
	defer conn.Close()

	_, err = conn.WriteToUDP(response.Bytes(), &addr)
	if err != nil {
		return fmt.Errorf("failed to send DATA packet: %v", err)
	}

	return nil
}
