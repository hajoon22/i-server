package client

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type client struct {
	EchoID      int
	IsRewritten bool

	Addr     string
	LastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = map[string]*client{}
)

func newClient(id string, echoID int, addr string) {
	mu.Lock()
	defer mu.Unlock()

	var (
		c  *client
		ok bool
	)

	c, ok = clients[id]
	if !ok {
		clients[id] = &client{
			EchoID:      echoID,
			IsRewritten: echoID != DefualtEchoID,
			Addr:        addr,
			LastSeen:    time.Now(),
		}

		c = clients[id]
		log.Printf("new client: %s (%s)\r\n", c.Addr, id)
	}

	c.LastSeen = time.Now()
}

func Maintain() {
	for {
		now := time.Now()

		mu.Lock()
		for id, c := range clients {
			if now.Sub(c.LastSeen) > ClientTimeout {
				delete(clients, id)
			}
		}
		mu.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func Listen() error {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("listen packet error")
	}
	defer conn.Close()

	buf := make([]byte, 1500)
	for {
		n, peer, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		packet, err := icmp.ParseMessage(1, buf[:n])
		if err != nil {
			continue
		}

		if packet.Type == ipv4.ICMPTypeEcho {
			echo, ok := packet.Body.(*icmp.Echo)
			if !ok || len(echo.Data) <= 0 {
				continue
			}

			if echo.Seq == KeepaliveEchoSeq {
				go newClient(string(echo.Data), echo.ID, peer.String())
			}
		}
	}
}
