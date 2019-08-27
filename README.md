# chord : Minimal DFS implementing part of the Chord algorithm

## Overview

Aims to fully implement the algorithm but for the time it still lacks a lot of features other implementations have. And most notably it doesn't support fingers thus has worst time O(n) time complexity for any call. This means it has major problems on retaining stability especially when nodes leave the network.

## Install

```
go get github.com/rand0x0m/chord
```

## Example

```
go run $GOPATH/src/github.com/rand0x0m/chord/examples/three-nodes-network/main.go
```

## License

MIT.

