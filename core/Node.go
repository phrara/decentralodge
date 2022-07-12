package core

import (
	"decentralodge/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	// p2p节点
	Host host.Host
	// 节点信息
	NodeInfo peer.AddrInfo
	// p2p节点标识
	NodeAddr multiaddr.Multiaddr
}

func GenerateNode() (*Node, error) {

	// 读取配置
	c := &config.Config{}
	c.Load()
	node := new(Node)

	// 获取节点
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(c.AddrString()),
		libp2p.Identity(c.PrvKey),
	)
	if err != nil {
		return nil, err
	}
	node.Host = h

	// 获取节点信息
	node.NodeInfo = peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}

	// 获取节点标识
	addrs, err := peer.AddrInfoToP2pAddrs(&node.NodeInfo)
	if err != nil {
		return nil, err
	}
	node.NodeAddr = addrs[0]

	return node, nil
}

func (n *Node) Close() error {
	err := n.Host.Close()
	return err
}
