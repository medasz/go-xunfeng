package vulscan

import (
	"github.com/medasz/kunpeng"
)

func Run() {
	println("开始...", len(kunpeng.Greeter.GetPlugins()))
}
