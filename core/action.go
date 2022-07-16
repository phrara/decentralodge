package core

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peerstore"
)

func (n *HostNode) JoinNetwork() {
	n.Router.Clear()
	for _, bn := range BootstrapNodes {
		node := bn
		n.Host.Peerstore().AddAddrs(node.NodeInfo.ID, node.NodeInfo.Addrs, peerstore.PermanentAddrTTL)
		go func() {
			res := <-n.Serv.Ping(node)
			if res.Error != nil {
				return
			} else {
				// ping 通了
				// 发出加入申请
				if b := n.Serv.JoinApply(node); b {
					fmt.Println("Join Network Successfully")
					n.Router.AddNode(node)
				}
			}
		}()
	}
}

func (n *HostNode) Close() error {
	err := n.Host.Close()
	n.Router.Clear()
	return err
}
