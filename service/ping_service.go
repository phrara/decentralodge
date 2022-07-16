package service

import (
	"context"
	"decentralodge/tool"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

func (s *Service) Ping(pn *tool.PeerNode) <-chan ping.Result {
	return s.pingService.Ping(context.Background(), pn.NodeInfo.ID)
}
