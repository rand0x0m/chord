package chord

import (
	"errors"
	"net"
	"time"
)

type Node struct {
	Info     *Info
	Db       Db
	Server   *Server
	Succ     *Info
	Pred     *Info
	Shutdown chan struct{}
}

//Put adds a value to the node's database.
func (n *Node) Put(value string, _ *struct{}) error {
	if value != "" {
		key := hashData(value)
		if inclusiveRightBetween(n.Pred.Id, key, n.Info.Id) {
			n.Db[key] = value
		}
	}

	return nil
}

//Get returns the value of the corresponding key from the database.
func (n *Node) Get(key string, value *string) error {
	*value = n.Db[key]
	return nil
}

//Delete removes the entry of the corresponding key from the database.
func (n *Node) Delete(key string, _ *struct{}) error {
	delete(n.Db, key)
	return nil
}

//Successor returns node's successor or error if the node is bootstrap.
func (n *Node) Successor(_ *struct{}, succ *Info) error {
	if n.Succ.Id == "" {
		return errors.New("not found")
	} else {
		*succ = *n.Succ
		return nil
	}
}

//Predecessor returns node's predecessor or error when the node  is bootstrap.
func (n *Node) Predecessor(_ *struct{}, pred *Info) error {
	if n.Pred.Id == "" {
		return errors.New("not found")
	} else {
		*pred = *n.Pred
		return nil
	}
}

//Notify the node about a possible new predecessor
func (n *Node) Notify(newPred Info, _ *struct{}) error {
	if n.isBootstrap() {
		if !n.Info.equal(&newPred) {
			n.Pred.Id = newPred.Id
			n.Pred.Addr = newPred.Addr
			n.Succ.Id = newPred.Id
			n.Succ.Addr = newPred.Addr
			n.migrateKeysToSuccessor()
		}
	} else if inclusiveBetween(n.Pred.Id, newPred.Id, n.Info.Id) {
		oldPredId := n.Pred.Id
		n.Pred.Id = newPred.Id
		n.Pred.Addr = newPred.Addr
		n.migrateKeysToPredecessor(oldPredId)
	}

	return nil
}

//Iteratively searches and returns the successor of a node. Returns the successor of the predecessor.
func (n *Node) FindSucc(id string, succ *Info) error {
	if n.isBootstrap() {
		*succ = *n.Info
		return nil
	}

	pred, err := remoteFindPred(n.Info, id)
	if err != nil {
		return err
	}

	temp, err := remoteSucc(pred)
	if err != nil {
		return err
	}

	*succ = *temp
	return nil
}

//Iteratively searches and returns the predecessor of a node.
func (n *Node) FindPred(id string, pred *Info) error {
	if n.isBootstrap() {
		*pred = *n.Info
		return nil
	}

	var err error
	nn := n.Info
	nnPred := n.Pred
	for !inclusiveRightBetween(nnPred.Id, id, nn.Id) {
		nn = nnPred
		nnPred, err = remotePred(nn)
		if err != nil {
			return err
		}
		if nn.equal(n.Info) {
			return errors.New("not found")
		}
	}

	*pred = *nnPred
	return nil
}

//A bootstrap node has no neighbors
func (n *Node) isBootstrap() bool {
	return n.Succ.Id == "" && n.Pred.Id == ""
}

func (n *Node) stabilize(interval time.Duration) {
	go func() {
		for {
			select {
			case <-n.Shutdown:
				return
			case <-time.After(interval):
				if !n.isBootstrap() {
					newSucc, err := remotePred(n.Succ)
					if err == nil && exclusiveBetween(n.Info.Id, newSucc.Id, n.Succ.Id) {
						n.Succ.Id = newSucc.Id
						n.Succ.Addr = newSucc.Addr
					}
				}
			}
		}
	}()
}

func (n *Node) migrateKeysToPredecessor(oldPredId string) {
	for k, v := range n.Db {
		if inclusiveRightBetween(oldPredId, k, n.Pred.Id) {
			remotePut(n.Pred, v)
			delete(n.Db, k)
		}
	}
}

func (n *Node) migrateKeysToSuccessor() {
	for k, v := range n.Db {
		if inclusiveRightBetween(n.Info.Id, k, n.Succ.Id) {
			remotePut(n.Succ, v)
			delete(n.Db, k)
		}
	}
}

func (n *Node) migrateDbToSuccessor() {
	for k, v := range n.Db {
		remotePut(n.Succ, v)
		delete(n.Db, k)
	}
}

func (n *Node) join(bootstrapAddr net.TCPAddr) error {
	var err error
	bootstrapInfo := newNodeInfo(bootstrapAddr)
	if n.Info.equal(bootstrapInfo) {
		return errors.New("cannot join itself")
	}

	//find our predecessor
	pred, err := remoteFindPred(bootstrapInfo, n.Info.Id)
	if err != nil {
		return err
	}
	n.Pred.Addr = pred.Addr
	n.Pred.Id = pred.Id

	//find our successor
	succ, err := remoteFindSucc(bootstrapInfo, n.Info.Id)
	if err != nil {
		return err
	}
	n.Succ.Addr = succ.Addr
	n.Succ.Id = succ.Id

	//let the successor now about us
	if err := remoteNotify(n.Succ, *n.Info); err != nil {
		return err
	}

	return nil
}

func (n *Node) listen() error {
	if n.Server != nil {
		return errors.New("server already running")
	}

	//start rpc server
	var err error
	n.Server, err = newServer(n.Info.Addr, *n)

	if err != nil {
		return err
	}

	//start stabilize
	n.stabilize(1 * time.Second)

	return nil
}

func (n *Node) shutdown() {
	close(n.Shutdown)
	n.migrateDbToSuccessor()
}

func NewNode(bindAddr net.TCPAddr) *Node {
	n := new(Node)
	n.Db = newDb()
	n.Info = newNodeInfo(bindAddr)
	n.Shutdown = make(chan struct{})
	n.Succ = &Info{}
	n.Pred = &Info{}

	return n
}
