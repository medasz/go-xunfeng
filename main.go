package main

import (
	"log"

	"go-xunfeng/nascan"
	"go-xunfeng/web"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	go nascan.Run()
	web.Web()
}
