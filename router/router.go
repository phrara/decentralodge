package router

import (
	"container/list"
	"decentralodge/tool"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var rt *Router

var mutex sync.Mutex

func init() {

	rt = &Router{
		table:    nil,
		HostNode: nil,
	}

}

type Table map[int]*list.List

type Router struct {
	table    Table
	HostNode *tool.PeerNode
}

func InitRouterTable(hn string) *Router {
	n := tool.ParsePeerNode(hn)
	rt.table = make(map[int]*list.List)
	rt.table[0] = list.New()
	rt.table[0].PushBack(n)
	rt.HostNode = n
	return rt
}

func (r *Router) AddNode(pn *tool.PeerNode) {
	hnode := r.HostNode
	dist := tool.GetDistByXor(hnode.NodeInfo.ID, pn.NodeInfo.ID)

	mutex.Lock()
	if r.table[dist] == nil {
		r.table[dist] = list.New()
		r.table[dist].PushBack(pn)
	} else {
		if b := r.ContainsAt(dist, pn); !b {
			r.table[dist].PushBack(pn)
		}
	}
	mutex.Unlock()
	return
}

func (r *Router) DelNode(pn *tool.PeerNode) {
	hnode := r.HostNode
	dist := tool.GetDistByXor(hnode.NodeInfo.ID, pn.NodeInfo.ID)

	mutex.Lock()
	if r.table[dist] == nil {
		return
	} else {
		for e := r.table[dist].Front(); e != nil; e = e.Next() {
			if e.Value.(*tool.PeerNode).NodeInfo.ID == pn.NodeInfo.ID {
				r.table[dist].Remove(e)
			}
		}
	}
	mutex.Unlock()
}

func (r *Router) Contains(pn *tool.PeerNode) bool {
	hnode := r.HostNode
	dist := tool.GetDistByXor(hnode.NodeInfo.ID, pn.NodeInfo.ID)
	if dist == 0 {
		return true
	}
	if r.table[dist] == nil {
		return false
	} else {
		for e := r.table[dist].Front(); e != nil; e = e.Next() {
			if e.Value.(*tool.PeerNode).NodeInfo.ID == pn.NodeInfo.ID {
				return true
			}
		}
	}
	return false
}

func (r *Router) ContainsAt(dist int, pn *tool.PeerNode) bool {
	if dist == 0 {
		return true
	}
	if r.table[dist] == nil {
		return false
	} else {
		for e := r.table[dist].Front(); e != nil; e = e.Next() {
			if e.Value.(*tool.PeerNode).NodeInfo.ID == pn.NodeInfo.ID {
				return true
			}
		}
	}
	return false
}

func (r *Router) Cap() int {
	size := 0
	for i := 0; i <= 256; i++ {
		size = size + r.table[i].Len()
	}
	return size
}

func (r *Router) Clear() {
	for i := 1; i <= 256; i++ {
		r.table[i] = nil
	}
}

func (r *Router) Update(table Table) {
	for i := 1; i <= 256; i++ {
		if table[i] != nil {
			for e := table[i].Front(); e != nil; e = e.Next() {
				r.AddNode(e.Value.(*tool.PeerNode))
			}
		}
	}
}

/*
	Sample:
	1:/ip4/127.0.0.1/tcp/2300/p2p/Qmx;/ip4/127.0.0.1/tcp/2301/p2p/Qmx;||2:/ip4/127.0.0.1/tcp/2302/p2p/Qmx;||
*/

func (r *Router) RawData() string {
	data := strings.Builder{}
	for i := 1; i <= 256; i++ {
		if r.table[i] != nil {
			data.WriteString(fmt.Sprintf("%d:", i))
			for e := r.table[i].Front(); e != nil; e = e.Next() {
				addrStr := e.Value.(*tool.PeerNode).NodeAddr.String()
				data.WriteString(addrStr + ";")
			}
			data.WriteString("||")
		}
	}
	return data.String()
}

func (r *Router) ParseData(raw string) Table {
	// The distances of addresses in every row are the same
	table := make(Table)
	distList := strings.Split(raw, "||")
	for _, str := range distList {
		if str == "" {
			continue
		}
		row := strings.Split(str, ":")
		addrs := strings.Split(row[1], ";")
		dist, _ := strconv.ParseInt(row[0], 10, 64)
		addrList := list.New()
		for _, addr := range addrs {
			if addr == "" {
				continue
			}
			addrList.PushBack(tool.ParsePeerNode(addr))
		}
		table[int(dist)] = addrList
	}
	return table
}
