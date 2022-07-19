package service

import (
	"bufio"
	"decentralodge/tool"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
)

func RecvFile(s network.Stream) {
	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	serv.router.AddNode(pn)
	fmt.Println("Get a file from", pn.String())

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	content, _ := rw.ReadString('\n')
	defer s.Close()
	fmt.Println(content)

}

func (s *Service) SendFile() {

}
