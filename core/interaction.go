package core

import (
	"context"
	"decentralodge/service"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"log"
)

const (
	CHAT = "/chat/1.0.0"
)

func (n *Node) ServiceHandlerInit() {
	n.Host.SetStreamHandler(CHAT, service.ChatHandler)

}

func (n *Node) Chat(p string) {
	maddr, err := multiaddr.NewMultiaddr(p)
	if err != nil {
		log.Println(err)
		return
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
		return
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	n.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	stream, err := n.Host.NewStream(context.Background(), info.ID, CHAT)
	if err != nil {
		log.Println(err)
		return
	}
	service.ChatWithPeer(stream)
}
