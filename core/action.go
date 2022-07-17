package core

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"time"
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

func (n *HostNode) routerDistribute(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 25)
	for t := range ticker.C {
		fmt.Println("Router Distribution Start: ", t.Format("2006-01-02 15:04:05"))
		errNum := n.Serv.RouterDistribute()

		// If the errNum >= 33% of the sum of nodes,
		// we regard this as the bad situation of network, then try again after 8 sec;
		// if the errNum >= 75% of the sum of nodes,
		// we regard this as the fatal error of network, stop the ticker
		if errNum >= n.Router.Sum()/4*3 {
			fmt.Println("fatal Network error")
			ticker.Stop()
			return
		} else if errNum >= n.Router.Sum()/3 {
			fmt.Println("Bad Network Situation")
			time.Sleep(time.Second * 10)
			n.Serv.RouterDistribute()
		}
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		default:
		}
	}
}

var cancel context.CancelFunc

func (n *HostNode) RouterDistributeOn() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancel = cancelFunc
	go n.routerDistribute(ctx)
}

func (n *HostNode) RouterDistributeOff() {
	cancel()
}
