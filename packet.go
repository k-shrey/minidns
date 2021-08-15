package main

import (
	"errors"
	"strings"
)

type Packet struct {
	buf [512]byte
	pos int
}

func NewPacket() Packet {
	packet := new(Packet)
	packet.pos = 0

	return *packet
}

func (p *Packet) step(size int) {
	p.pos += size
}

func (p *Packet) seek(pos int) {
	p.pos = pos
}

// read one byte and increment position
func (p *Packet) read() (byte, error) {
	if p.pos >= 512 {
		return 0, errors.New("end of buffer")
	}
	res := p.buf[p.pos]
	p.pos += 1

	return res, nil
}

//read one byte without incrementing position
func (p Packet) get(pos int) (byte, error) {
	if p.pos >= 512 {
		return 0, errors.New("end of buffer")
	}
	res := p.buf[pos]

	return res, nil
}

func (p Packet) getRange(start, len int) ([]byte, error) {
	if start+len >= 512 {
		return nil, errors.New("end of buffer")
	}
	return p.buf[start : start+len], nil
}

func (p *Packet) readU16() uint16 {
	r1, _ := p.read()
	r2, _ := p.read()
	res := uint16(r1)<<8 | uint16(r2)
	return res
}

func (p *Packet) readU32() uint32 {
	// res := uint32(p.pos)<<24 |
	// 	uint32(p.pos)<<16 |
	// 	uint32(p.pos)<<8 |
	// 	uint32(p.pos)<<0
	res := uint32(0)
	for i := 3; i >= 0; i-- {
		shift := 8 * i
		val, _ := p.read()
		shifted := uint32(val) << shift
		res |= shifted
	}

	return res
}

func (p *Packet) readQName() (string, error) {
	pos := p.pos
	outbuf := ""
	// handle jumps in the case of compressed records
	jumped := false
	maxJumps := 5
	jumps := 0

	delim := ""

	for {
		if jumps > maxJumps {
			return "", errors.New("limit of jumps exceeded")
		}

		//read length byte
		len, err := p.get(pos)
		if err != nil {
			return "", err
		}

		// if the 2 most significant bits are set,
		// then jump to some other offset in the packet
		if len&0xC0 == 0xC0 {

			// Update position to after the current label
			if !jumped {
				p.seek(pos + 2)
			}

			//read the next byte
			secondByte, err := p.get(pos + 1)
			if err != nil {
				return "", err
			}

			//calculate offset and jump
			offset := ((uint16(len) ^ 0xC0) << 8) | uint16(secondByte)
			pos = int(offset)
			jumped = true
			jumps += 1
			continue

		} else {
			// default case, no jumps. Read a single label and append it to the output

			// move  past the length byte
			pos += 1

			// if the length is 0, we have reached the end of the domain name
			if len == 0 {
				break
			}

			//append delimiter to the output buffer
			outbuf += delim
			// fmt.Println("pos is: ", pos)
			// extract output bytes and append it to the output buffer
			strBuffer, err := p.getRange(pos, int(len))
			if err != nil {
				return "", err
			}
			outbuf += strings.ToLower(string(strBuffer))

			delim = "."

			// move forward to the next record
			pos += int(len)

		}
	}

	if !jumped {
		p.seek(pos)
	}

	// fmt.Println("outbuf: ", *outbuf)
	return outbuf, nil
}

func (p *Packet) write(val uint8) error {
	if p.pos >= 512 {
		return errors.New("end of buffer")
	}

	p.buf[p.pos] = val
	p.pos += 1
	return nil
}

func (p *Packet) writeU16(val uint16) {
	p.write(uint8(val >> 8))
	p.write(uint8(val & 0xFF))
}

func (p *Packet) writeU32(val uint32) {
	p.write(uint8((val >> 24) & 0xFF))
	p.write(uint8((val >> 16) & 0xFF))
	p.write(uint8((val >> 8) & 0xFF))
	p.write(uint8((val >> 0) & 0xFF))
}

func (p *Packet) writeQName(qname string) error {
	for _, label := range strings.Split(qname, ".") {
		len := len(label)
		if len > 0x3f {
			return errors.New("Label is longer than 63 chars")
		}

		p.write(uint8(len))
		for _, c := range []byte(label) {
			p.write(c)
		}
	}

	p.write(0)

	return nil
}
