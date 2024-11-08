package main

import (
	"go-xunfeng/vulscan"
	"log"

	"go-xunfeng/nascan"
	"go-xunfeng/web"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	go nascan.Run()
	vulscan.Run()
	web.Web()
}
