package main

import (
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/command"
)

func main() {
	app.InitConfig()
	app.InitLog()
	command.Run()
}
