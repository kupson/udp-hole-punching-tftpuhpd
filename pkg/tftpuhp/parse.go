package tftp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

const (
	opRRQ          = 1
	filenamePrefix = "v1_tftp_udp_"
)

var (
	filenameRegex = regexp.MustCompile(fmt.Sprintf(`^%s[0-9]+$`, filenamePrefix))

	minRRQLength = 4 + len(filenamePrefix)
)

type RRQ struct {
	Filename string
}

func parseRRQ(packet []byte) (*RRQ, error) {
	if len(packet) < minRRQLength {
		return nil, fmt.Errorf("packet too short: got %d bytes", len(packet))
	}

	opcode := binary.BigEndian.Uint16(packet[:2])
	if opcode != opRRQ {
		return nil, fmt.Errorf("not a RRQ")
	}

	remainder := packet[2:]
	parts := bytes.SplitN(remainder, []byte{0}, 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("malformed RRQ")
	}

	filename := string(parts[0])
	mode := strings.ToLower(string(parts[1]))

	if !filenameRegex.MatchString(filename) {
		return nil, fmt.Errorf("filename does not match required pattern: %s", filename)
	}

	if mode != "octet" {
		return nil, fmt.Errorf("unsupported mode")
	}

	return &RRQ{
		Filename: filename,
	}, nil
}
