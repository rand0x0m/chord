package main

import (
	"fmt"
	"net"
	"time"

	"github.com/rand0x0m/chord"
)

var dfs1, dfs2, dfs3 chord.DFS

var testAddrs = [6]net.TCPAddr{
	{IP: net.ParseIP("127.0.0.1"), Port: 5483},
	{IP: net.ParseIP("127.0.0.1"), Port: 5484},
	{IP: net.ParseIP("127.0.0.1"), Port: 5485},
}

var testValues = [6]string{
	"TestValue",
	"TestValue87",
	"TestValue90",
	"TestValue20",
	"TestValue6",
	"TestValue65",
}

func main() {

	dfs1, err := chord.NewDFS(&chord.Config{&testAddrs[0], nil})
	if err != nil {
		panic(err)
	}

	dfs2, err := chord.NewDFS(&chord.Config{&testAddrs[1], &testAddrs[0]})
	if err != nil {
		panic(err)
	}

	dfs3, err := chord.NewDFS(&chord.Config{&testAddrs[2], &testAddrs[0]})
	if err != nil {
		panic(err)
	}

	if err := dfs2.Connect(); err != nil {
		panic(0)
	}

	time.Sleep(5 * time.Second)

	if err := dfs3.Connect(); err != nil {
		panic(0)
	}

	time.Sleep(5 * time.Second)

	for _, v := range testValues {
		err := dfs1.Put(v)
		if err != nil {
			panic(err)
		}
	}

	for _, v := range testValues {
		v, err := dfs1.Get(v)
		fmt.Println(v)
		if err != nil {
			panic(err)
		}
	}

	dfs1.Shutdown()

	time.Sleep(3 * time.Second)

	dfs2.Shutdown()

	time.Sleep(3 * time.Second)

	dfs3.Shutdown()
}
