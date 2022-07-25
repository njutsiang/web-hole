package main

import (
	"hole/app"
	"hole/command"
)

func main() {
	app.InitConfig()
	command.Run()
}
