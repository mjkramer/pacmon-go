package main

// import (
// 	"encoding/binary"
// )

func Parity64(data *Packet) byte {
	// x := binary.LittleEndian.Uint64(data)
    x := uint64(data[0]) << 56 | uint64(data[1]) << 48 | uint64(data[2]) << 40 | uint64(data[3]) << 32 |
         uint64(data[4]) << 24 | uint64(data[5]) << 16 | uint64(data[6]) << 8  | uint64(data[7])
	x ^= x >> 32
	x ^= x >> 16
	x ^= x >> 8
	x ^= x >> 4
	x ^= x >> 2
	return byte(x) & 1
}
