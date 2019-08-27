package chord

import (
	"fmt"
	"net"
)

type Info struct {
	Addr net.TCPAddr
	Id   string
}

func (i *Info) equal(o *Info) bool {
	return i.Id == o.Id
}

func (i *Info) String() string {
	return fmt.Sprintf("Network:%v:%v Id:%v", i.Addr.Network(), i.Addr.String(), i.Id)
}

func newNodeInfo(addr net.TCPAddr) *Info {
	i := new(Info)
	i.Addr = addr
	i.Id = hashAddr(addr)

	return i
}
