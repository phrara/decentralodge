package main

import (
	"decentralodge/core"
	"fmt"
	"time"
)

func main() {

	hnode, err := core.GenerateNode()
	if err != nil {
		return
	}
	fmt.Println(hnode.NodeAddr.String())

	hnode.JoinNetwork()
	time.Sleep(time.Second * 2)
	fmt.Println(string(hnode.Router.RawData()))

	hnode.RouterDistributeOn(false, 5)
	time.Sleep(time.Second * 2)
	hnode.Serv.SendFile(core.BootstrapNodes[0], "A test file content")

	select {}
}
