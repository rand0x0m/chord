package chord

import (
	"errors"
	"net"
)

type DFS struct {
	N      *Node
	Config *Config
}

type Config struct {
	BindAddr      *net.TCPAddr
	BootstrapAddr *net.TCPAddr
}

func (dfs *DFS) Put(value string) error {
	if dfs.N.isBootstrap() {
		return errors.New("disconnected")
	} else {
		succ, err := remoteFindSucc(dfs.N.Info, hashData(value))
		if err != nil {
			return err
		}
		return remotePut(succ, value)
	}
}

func (dfs *DFS) Get(value string) (string, error) {
	key := hashData(value)
	if dfs.N.isBootstrap() {
		return "", errors.New("disconnected")
	} else {
		succ, err := remoteFindSucc(dfs.N.Info, key)
		if err != nil {
			return "", err
		}

		return remoteGet(succ, key)
	}
}

func (dfs *DFS) Delete(value string) error {
	key := hashData(value)
	if dfs.N.isBootstrap() {
		return errors.New("disconnected")
	} else {
		succ, err := remoteFindSucc(dfs.N.Info, key)
		if err != nil {
			return err
		}

		return remoteDelete(succ, key)
	}
}

func (dfs *DFS) Shutdown() {
	dfs.N.shutdown()
}

func (dfs *DFS) Connect() error {
	err := dfs.N.join(*dfs.Config.BootstrapAddr)
	if err != nil {
		return err
	}

	return nil
}

func NewDFS(config *Config) (*DFS, error) {
	dfs := new(DFS)
	dfs.N = NewNode(*config.BindAddr)
	dfs.Config = config

	if err := dfs.N.listen(); err != nil {
		return nil, err
	}

	return dfs, nil
}
