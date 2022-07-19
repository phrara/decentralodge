package service

import (
	"bufio"
	"context"
	"decentralodge/tool"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
)

// Router Distribution
// Router table will be distributed periodically
// When peers get the distributed router tables, they use then to update their own router tables
// This automatically renew the router info of the decentralized network

func RouterDistributeHandler(s network.Stream) {
	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	serv.router.AddNode(pn)
	fmt.Println("Get a distributed router table from", pn.String())

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	router, _ := rw.ReadString('\n')
	defer s.Close()
	if router == "" || router == "\n" {
		return
	} else {
		//Parse the raw data and use it to update the local router table
		fmt.Println("remote router info is: \n", router)
		data := serv.router.ParseData(router)
		serv.router.Update(data)
	}
}

func (s *Service) RouterDistribute() (errNum int) {

	var wg sync.WaitGroup

	localRouter := serv.router.RawData()
	nodes := serv.router.AllNodes()
	for e := nodes.Front(); e != nil; e = e.Next() {
		pn := e.Value.(*tool.PeerNode)

		wg.Add(1)
		go func(p *tool.PeerNode) {

			defer wg.Done()

			var err error = nil
			var stream network.Stream = nil
			if b := <-s.Ping(p); b.Error == nil {
				s.Host.Peerstore().AddAddrs(p.ID(), p.NodeInfo.Addrs, peerstore.PermanentAddrTTL)
				stream, err = s.Host.NewStream(context.Background(), p.ID(), RD)
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
				rw.WriteString(localRouter + "\n")
				rw.Flush()
			} else if b.Error != nil || err != nil {
				// Record the number of errors to evaluate the situation of network
				fmt.Println(b.Error, err)
				errNum = errNum + 1
				return
			}

		}(pn)
	}

	wg.Wait()
	return errNum
}
