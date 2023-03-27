package client

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

func randomIP(name string) (ip string) {
	var b [8]byte
	for i, v := range name[0:8] {
		b[i] = byte(v)
	}
	var i int = 0
	for _, v := range name[7:] {
		if i > 7 {
			i -= 8
		}
		b[i] = b[i] + byte(v)
		i++
	}
	s := binary.LittleEndian.Uint64(b[:])
	rnd := rand.New(rand.NewSource(int64(s)))
	var lo [1]byte
	rnd.Read(lo[:])
	ip = "10.0.0." + fmt.Sprint(lo[0])
	return ip
}
