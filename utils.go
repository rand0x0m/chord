package chord

import (
	"crypto/sha1"
	"io"
	"math/big"
	"net"
	"net/rpc"
)

func hashData(s string) string {
	hasher := sha1.New()
	io.WriteString(hasher, s)
	return new(big.Int).SetBytes(hasher.Sum(nil)).String()
}

func hashAddr(addr net.TCPAddr) string {
	return hashData(addr.String())
}

func exclusiveBetween(left, i, right string) bool {
	if left < i && i < right {
		return true
	}

	if i < right && right < left {
		return true
	}

	if right < left && left < i {
		return true
	}

	return false
}

func inclusiveBetween(left, i, right string) bool {
	if left <= i && i <= right {
		return true
	}

	if i <= right && right <= left {
		return true
	}

	if right <= left && left <= i {
		return true
	}

	return false
}

func inclusiveRightBetween(left, i, right string) bool {
	if left < i && i <= right {
		return true
	}

	if i <= right && right < left {
		return true
	}

	if right < left && left < i {
		return true
	}

	return false
}

func setupConnection(server *Info) (*rpc.Client, error) {
	return rpc.DialHTTP(server.Addr.Network(), server.Addr.String())
}
