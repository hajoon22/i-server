package client

import (
	"fmt"
	"log"
	"math/rand/v2"
	"server/config"
	"server/protocol"
	"strings"
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

const (
	DefualtEchoID    = 2222
	MessageEchoSeq   = 1111
	KeepaliveEchoSeq = 2222
	ClientTimeout    = 15 * time.Second
)

var (
	mu      sync.Mutex
	clients = map[string]*client{}
)

func Fetch() map[string]client {
	mu.Lock()
	defer mu.Unlock()

	results := map[string]client{}
	for id, c := range clients {
		results[id] = *c
	}

	return results
}

func (c *client) sendMessage(cfg *config.Config, message string) {
	dst := strings.Split(c.Addr, ".")
	if len(dst) != 4 {
		return
	}

	var refAddr string
	for {
		refAddr = cfg.Servers[rand.IntN(len(cfg.Servers))]
		if refAddr == cfg.ServerAddr {
			continue
		}

		break
	}

	ref := strings.Split(refAddr, ".")
	if len(ref) != 4 {
		return
	}

	inner, err := protocol.BuildECHO([4]string(dst), [4]string(ref), c.EchoID, MessageEchoSeq, message)
	if err != nil {
		return
	}

	if c.IsRewritten {
		protocol.SendUnreach(c.Addr, inner)
	} else {
		protocol.SendIPIP(cfg.RelayAddr, inner)
	}
}

func NewMessage(cfg *config.Config, message string) int {
	mu.Lock()
	defer mu.Unlock()

	counter := 0
	for _, c := range clients {
		counter++
		c.sendMessage(cfg, message)
	}

	return counter
}

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
