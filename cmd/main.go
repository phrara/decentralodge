package main

import (
	"decentralodge/core"
	"fmt"
	"time"
)

const LOGO = `
'        ___           ___           ___       ___       ___     
'       /\__\         /\  \         /\__\     /\__\     /\  \    
'      /:/  /        /::\  \       /:/  /    /:/  /    /::\  \   
'     /:/__/        /:/\:\  \     /:/  /    /:/  /    /:/\:\  \  
'    /::\  \ ___   /::\~\:\  \   /:/  /    /:/  /    /:/  \:\  \ 
'   /:/\:\  /\__\ /:/\:\ \:\__\ /:/__/    /:/__/    /:/__/ \:\__\
'   \/__\:\/:/  / \:\~\:\ \/__/ \:\  \    \:\  \    \:\  \ /:/  /
'        \::/  /   \:\ \:\__\    \:\  \    \:\  \    \:\  /:/  / 
'        /:/  /     \:\ \/__/     \:\  \    \:\  \    \:\/:/  /  
'       /:/  /       \:\__\        \:\__\    \:\__\    \::/  /   
'       \/__/         \/__/         \/__/     \/__/     \/__/    `

func main() {

	fmt.Printf("\u001B[1;35m%s\u001B[0m\n", LOGO)
	hnode, err := core.GenerateNode()
	if err != nil {
		return
	}
	fmt.Printf("\x1b[1;34mHost: %s\x1b[0m\n", hnode.NodeAddr.String())

	hnode.JoinNetwork()
	time.Sleep(time.Second * 2)
	fmt.Println("recent router table:\n", string(hnode.Router.RawData()))

	hnode.RouterDistributeOn(false, 5)
	time.Sleep(time.Second * 2)
	//hnode.Serv.SendFile(core.BootstrapNodes[0], "A test file content")

	select {}
}
