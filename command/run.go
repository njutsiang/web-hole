package command

import (
	"github.com/njutsiang/web-hole/command/controllers"
	"os"
)

// 执行命令
func Run() bool {
	if len(os.Args) < 2 {
		return false
	}
	command, ok := GetRouters()[os.Args[1]]
	if ok {
		command()
		return true
	} else {
		return false
	}
}

// 定义命令路由
func GetRouters() map[string]func() {
	return map[string]func(){
		"StartFrontend": controllers.StartFrontend,
		"StartProxy": controllers.StartProxy,
		"StartBackend": controllers.StartBackend,
	}
}
