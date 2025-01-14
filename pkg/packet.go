package pkg

import (
	"math/bits"
)

type Packet [8]byte

type PacketType uint8

const (
	PacketTypeData PacketType = 0
	PacketTypeError PacketType = 1
	PacketTypeWrite PacketType = 2
	PacketTypeRead PacketType = 3
)

// TODO: Test the setters

func (p Packet) Type() PacketType {
	return PacketType(p[0] & 3)
}

func (p *Packet) SetType(typ PacketType) {
	p[0] = (p[0] & 0xFC) | (uint8(typ) & 0x03)
}

func (p Packet) Chip() uint8 {
	return (p[0] >> 2) | (p[1] << 6)
}

func (p *Packet) SetChip(chip uint8) {
	p[0] = (p[0] & 0x03) | (chip << 2)
	p[1] = (p[1] & 0xFC) | (chip >> 6)
}

func (p Packet) Channel() uint8 {
	return p[1] >> 2
}

func (p *Packet) SetChannel(channel uint8) {
	p[1] = (p[1] & 0x03) | (channel << 2)
}

func (p Packet) Timestamp() uint32 {
	return uint32(p[2]) |
		(uint32(p[3]) << 8) |
		(uint32(p[4]) << 16) |
		(uint32((p[5] & 0x7F)) << 24)
}

func (p *Packet) SetTimestamp(timestamp uint32) {
	p[2] = byte(timestamp & 0xFF)
	p[3] = byte((timestamp >> 8) & 0xFF)
	p[4] = byte((timestamp >> 16) & 0xFF)
	p[5] = (p[5] & 0x80) | byte((timestamp >> 24) & 0x7F)
}

func (p Packet) First() bool {
	return p[5] >> 7 == 1
}

func (p *Packet) SetFirst(first bool) {
	if first {
		p[5] = (p[5] & 0x7F) | 0x70
	} else {
		p[5] = p[5] & 0x7F
	}
}

func (p Packet) Data() uint8 {
	return p[6]
}

func (p *Packet) SetData(data uint8) {
	p[6] = data
} 

func (p Packet) TrigType() uint8 {
	return p[7] & 3
}

func (p *Packet) SetTrigType(trigtype uint8) {
	p[7] = (p[7] & 0xFC) | (trigtype & 0x03)
}

func (p Packet) LocalFifoFlags() uint8 {
	return (p[7] >> 2) & 3
}

func (p *Packet) SetLocalFifoFlags(flags uint8) {
	p[7] = (p[7] & 0xF3) & ((flags & 0x03) << 2)
}

func (p Packet) SharedFifoFlags() uint8 {
	return (p[7] >> 4) & 3
}

func (p *Packet) SetSharedFifoFlags(flags uint8) {
	p[7] = (p[7] & 0xCF) & ((flags & 0x03) << 4)
}

func (p Packet) Downstream() bool {
	return (p[7] >> 6) & 1 == 1
}

func (p *Packet) SetDownstream(downstream bool) {
	if downstream {
		p[7] = (p[7] & 0xBF) | 0x40
	} else {
		p[7] = (p[7] & 0xBF)
	}
}

func (p Packet) ParityBit() uint8 {
	return p[7] >> 7
}

func (p *Packet) SetParityBit(b uint8) {
	p[7] = (p[7] & 0x7F) | (b << 7)
}

func (p *Packet) UpdateParity() {
	if ! p.ValidParity() {
		if p.ParityBit() == 1 {
			p.SetParityBit(0)
		} else {
			p.SetParityBit(1)
		}
	}
}

func (p Packet) ValidParity() bool {
	onesCount := 0
	for i, b := range p {
		if i == 7 {
			onesCount = onesCount + bits.OnesCount(uint(b & 0x7F)) // Skip parity bit
		} else {
			onesCount = onesCount + bits.OnesCount(uint(b))
		}
	}
	return (1 - (onesCount % 2)) == int(p.ParityBit())
}
