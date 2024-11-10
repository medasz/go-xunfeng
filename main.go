package main

import (
	"log"

	"go-xunfeng/nascan"
	"go-xunfeng/vulscan"
	"go-xunfeng/web"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	go nascan.Run()
	vulscan.Run()
	web.Web()
}
