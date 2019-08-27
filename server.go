package chord

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
)

type Server struct {
	httpServer *http.Server
}

func newServer(bindAddr net.TCPAddr, n Node) (*Server, error) {
	RPCserver := rpc.NewServer()
	err := RPCserver.Register(&n)

	if err != nil {
		return nil, err
	}

	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux

	RPCserver.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	http.DefaultServeMux = oldMux

	listener, err := net.Listen(bindAddr.Network(), bindAddr.String())

	if err != nil {
		return nil, err
	}

	server := new(Server)
	server.httpServer = new(http.Server)
	server.httpServer.Handler = mux

	go server.httpServer.Serve(listener)

	go func() {
		<-n.Shutdown
		server.httpServer.Shutdown(context.Background())
	}()

	return server, nil
}
