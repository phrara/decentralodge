package service

import (
	"bufio"
	"context"
	"decentralodge/tool"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
)

func RecvFileHandler(s network.Stream) {
	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	serv.router.AddNode(pn)
	fmt.Println("Get a file from", pn.String())

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	file, _ := rw.ReadString('>')
	defer s.Close()
	fmt.Println(tool.UnwrapFile(file))

}

func (s *Service) SendFile(pn *tool.PeerNode, file string) bool {
	s.Host.Peerstore().AddAddrs(pn.ID(), pn.NodeInfo.Addrs, peerstore.PermanentAddrTTL)

	if r := <-s.Ping(pn); r.Error == nil {
		stream, err := s.Host.NewStream(context.Background(), pn.ID(), FT)
		if err != nil {
			fmt.Println(err)
			return false
		}
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		rw.WriteString(tool.WrapFile(file))
		rw.Flush()
		fmt.Println("send a file successfully")
		return true
	} else {
		return false
	}
}
