package tools

import (
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

func Ping(ip string) error {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	defer conn.Close()
	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Hello, are you there!"),
		},
	}
	wb, err := msg.Marshal(nil)
	if err != nil {
		return err
	}
	dst, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return err
	}
	start := time.Now()
	if _, err := conn.WriteTo(wb, dst); err != nil {
		log.Println(err)
	}
	// 接收 ICMP 报文
	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	// 解析 ICMP 报文
	msg, err = icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return err
	}
	// 打印结果
	switch msg.Type {
	case ipv4.ICMPTypeEchoReply:
		echoReply, ok := msg.Body.(*icmp.Echo)
		if !ok {
			return errors.New("invalid ICMP Echo Reply message")
		}
		if peer.String() == ip && echoReply.ID == os.Getpid()&0xffff && echoReply.Seq == 1 {
			fmt.Printf("reply from %s: seq=%d time=%v\n", dst.String(), msg.Body.(*icmp.Echo).Seq, duration)
			return nil
		}
	default:
		return errors.New(fmt.Sprintf("unexpected ICMP message type: %v\n", msg.Type))
	}
	return nil
}
