package core

import (
	"context"
	"decentralodge/config"
	"decentralodge/router"
	"decentralodge/service"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

type HostNode struct {
	// p2p节点
	Host host.Host
	// 节点信息
	NodeInfo *peer.AddrInfo
	// p2p节点标识
	NodeAddr multiaddr.Multiaddr
	// context
	Ctx context.Context
	// 相关协议服务
	Serv *service.Service
	// 路由表
	Router *router.Router
}

func GenerateNode() (*HostNode, error) {

	// 读取配置
	c := &config.Config{}
	c.Load()
	node := new(HostNode)
	node.Ctx = context.Background()

	// 获取节点
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(c.AddrString()),
		libp2p.Identity(c.PrvKey),
		libp2p.Ping(false),
	)
	if err != nil {
		return nil, err
	}
	node.Host = h

	// 获取节点信息
	node.NodeInfo = &peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
		// 获取节点标识
	}

	addrs, err := peer.AddrInfoToP2pAddrs(node.NodeInfo)
	if err != nil {
		return nil, err
	}

	node.NodeAddr = addrs[0]

	// 初始化路由表
	node.Router = router.InitRouterTable(node.NodeAddr.String())
	// 初始化协议服务
	node.Serv = service.NewService(node.Host, node.Router).ServiceHandlerRegister()

	return node, nil
}

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
