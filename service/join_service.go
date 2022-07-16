package service

import (
	"bufio"
	"context"
	"decentralodge/tool"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"strings"
)

func JoinApplyHandler(s network.Stream) {
	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	fmt.Println("receive a join application from", pn.NodeAddr)

	// 节点加入路由表
	serv.router.AddNode(pn)
	fmt.Println("recent router table:\n", serv.router.RawData())

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	rw.WriteString("acc\n")
	rw.Flush()

}

func (s *Service) JoinApply(bootstrapNode *tool.PeerNode) bool {
	stream, err := s.Host.NewStream(context.Background(), bootstrapNode.NodeInfo.ID, JOIN)
	if err != nil {
		fmt.Println(err)
		return false
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	ack, _ := rw.ReadString('\n')

	defer stream.Close()

	if ack == "" {
		return false
	}
	if ack != "\n" {
		if strings.Contains(ack, "acc") {
			return true
		} else {
			return false
		}
	}
	return false
}
