package main

import (
	"fmt"
	"go-xunfeng/config"
	"go-xunfeng/web"
)

func main() {
	web.Web()
	fmt.Println(config.Cfg)
}
