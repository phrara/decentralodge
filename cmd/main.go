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
	fmt.Println(hnode.NodeAddr)

	hnode.JoinNetwork()
	time.Sleep(time.Second * 2)
	fmt.Println(hnode.Router.RawData())

	select {}
}
