package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func checksum(b []byte) uint16 {
	var sum uint32
	for i := 0; i < len(b)-1; i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	return ^uint16(sum)
}

func buildIPHdr(src, dst [4]string, protocol byte, size int) []byte {
	hdr := make([]byte, 20)
	hdr[0] = 0x45
	binary.BigEndian.PutUint16(hdr[2:], uint16(20+size))
	hdr[8] = 64
	hdr[9] = protocol
	for i := 0; i < 4; i++ {
		s, _ := strconv.Atoi(src[i])
		d, _ := strconv.Atoi(dst[i])

		hdr[12+i] = byte(s)
		hdr[16+i] = byte(d)
	}

	binary.BigEndian.PutUint16(hdr[10:], checksum(hdr))
	return hdr
}

func buildEchoRequest(id, seq int, message string) ([]byte, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  seq,
			Data: []byte(message),
		},
	}

	buf, err := msg.Marshal(nil)
	if err != nil {
		return nil, fmt.Errorf("msg marshal error: %w", err)
	}

	return buf, nil
}

func BuildECHO(src, dst [4]string, id, seq int, message string) ([]byte, error) {
	icmph, err := buildEchoRequest(id, seq, message)
	if err != nil {
		return nil, fmt.Errorf("build echo request error: %w", err)
	}

	iph := buildIPHdr(src, dst, 1, len(icmph))
	return append(iph, icmph...), nil
}

// i am poor...
func buildUnreach(data []byte) ([]byte, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeDestinationUnreachable,
		Code: 3,
		Body: &icmp.DstUnreach{
			Data: data,
		},
	}
	return msg.Marshal(nil)
}

func SendIPIP(dst string, inner []byte) error {
	conn, err := net.Dial("ip4:4", dst)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(inner)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

func SendUnreach(dst string, data []byte) error {
	icmph, err := buildUnreach(data)
	if err != nil {
		return fmt.Errorf("build unreach error: %w", err)
	}

	conn, err := net.Dial("ip4:icmp", dst)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(icmph)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}
