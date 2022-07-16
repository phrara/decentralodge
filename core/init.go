package core

import (
	"decentralodge/tool"
)

// BootstrapNodes 引导节点
var BootstrapNodes []*tool.PeerNode

func init() {

	BootstrapNodes = getBootstrapNodes()

}

func getBootstrapNodes() []*tool.PeerNode {
	bsn := make([]*tool.PeerNode, 0)
	pn := tool.ParsePeerNode("/ip4/127.0.0.1/tcp/8083/p2p/QmPjfJ4pPUyScEG68jC44tjY7NQ7fdjFszDM3i69YTJBMH")
	bsn = append(bsn, pn)
	return bsn
}
