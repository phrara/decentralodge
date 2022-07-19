package core

import (
	"context"
	"fmt"
	"time"
)

func (n *HostNode) JoinNetwork() {
	n.Router.Clear()
	for _, bn := range BootstrapNodes {
		node := bn
		if node == nil {
			continue
		}
		go func() {
			res := <-n.Serv.Ping(node)
			if res.Error != nil {
				fmt.Println("can not connect Bootstrap Node")
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

func (n *HostNode) routerDistribute(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(time.Second * period)
	for t := range ticker.C {
		fmt.Println("Router Distribution Start: ", t.Format("2006-01-02 15:04:05"))
		errNum := n.Serv.RouterDistribute()

		// If the errNum > 33% of the sum of nodes,
		// we regard this as the bad situation of network, then try again after 8 sec;
		// if the errNum > 75% of the sum of nodes,
		// we regard this as the fatal error of network, stop the ticker
		if errNum > n.Router.Sum()/4*3 {
			fmt.Println("fatal Network error")
			fmt.Println("Please restart your server")
			ticker.Stop()
			return
		} else if errNum > n.Router.Sum()/3 {
			fmt.Println("Bad Network Situation")
			time.Sleep(time.Second * (period / 2))
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

func (n *HostNode) autoPeriod(nodeNum int) time.Duration {
	switch {
	case nodeNum >= 100:
		return 60
	case nodeNum >= 50:
		return 30
	case nodeNum >= 25:
		return 15
	default:
		return 10
	}
}

var cancel context.CancelFunc

// RouterDistributeOn
// The unit of argument `period` is second.
// The argument `period` will be affective if `auto` is true,
// and it'll be useless if `auto` is false.
func (n *HostNode) RouterDistributeOn(auto bool, period int) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancel = cancelFunc
	if auto {
		go n.routerDistribute(ctx, n.autoPeriod(n.Router.Sum()))
	} else {
		go n.routerDistribute(ctx, time.Duration(period))
	}
}

func (n *HostNode) RouterDistributeOff() {
	cancel()
}
