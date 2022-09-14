package nascan

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

type Nascan struct {
	timeout time.Duration
	data    []byte
}

func NewNascan() (Nascan, error) {
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(time.Now().Format("2006-01-02 15:04:05")),
		},
	}
	wb, err := msg.Marshal(nil)
	return Nascan{
		timeout: time.Second * 3,
		data:    wb,
	}, err
}

func (n Nascan) mPing(ips []string) (data []string) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Println(err)
		return data
	}
	defer conn.Close()
	ch := make(chan string, 1)
	go func() {
		for _, ip := range ips {
			dst, err := net.ResolveIPAddr("ip", ip)
			if err != nil {
				log.Println(err)
				continue
			}
			if _, err := conn.WriteTo(n.data, dst); err != nil {
				log.Println(err)
			}
			ch <- ip
		}
		close(ch)
	}()
	for v := range ch {
		// 接收 ICMP 报文
		reply := make([]byte, 1500)
		err = conn.SetReadDeadline(time.Now().Add(n.timeout))
		if err != nil {
			log.Println(err)
			continue
		}
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			log.Println(v, err)
			continue
		}
		// 解析 ICMP 报文
		msg, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			log.Println(err)
			continue
		}
		if msg.Type == ipv4.ICMPTypeEchoReply {
			_, ok := msg.Body.(*icmp.Echo)
			if ok {
				log.Println(peer.String(), "active")
				data = append(data, peer.String())
			}
		}
	}
	return data
}
