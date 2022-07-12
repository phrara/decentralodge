package main

import (
	"decentralodge/core"
	"fmt"
)

func main() {

	node, err := core.GenerateNode()
	if err != nil {
		return
	}
	fmt.Println(node.NodeAddr)

	node.ServiceHandlerInit()

	select {}

}
