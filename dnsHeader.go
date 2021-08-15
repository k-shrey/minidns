package main

type DnsHeader struct {
	id                  uint16
	recursionDesired    bool  // 1 bit
	truncatedMessage    bool  // 1 bit
	authoritativeAnswer bool  // 1 bit
	opcode              uint8 // 4 bits, but no representation in golang
	response            bool  // 1 bit

	rescode            uint8 // 4 bits, but no representation in golang
	checkingDisabled   bool  // 1 bit
	authedData         bool  // 1 bit
	z                  bool  // 1 bit
	recursionAvailable bool  // 1 bit

	questions            uint16 // 16 bits
	answers              uint16 // 16 bits
	authoritativeEntries uint16 // 16 bits
	resourceEntries      uint16 // 16 bits
}

func NewDnsHeader() *DnsHeader {
	h := new(DnsHeader)
	h.id = 0
	h.recursionDesired = false
	h.truncatedMessage = false
	h.authoritativeAnswer = false
	h.opcode = 0
	h.response = false
	h.rescode = uint8(NOERROR)
	h.checkingDisabled = false
	h.authedData = false
	h.z = false
	h.recursionAvailable = false
	h.questions = 0
	h.answers = 0
	h.authoritativeEntries = 0
	h.resourceEntries = 0

	return h
}

func (h *DnsHeader) read(buf *Packet) {
	// fmt.Printf("packet is: %#v", buf.buf)
	h.id = buf.readU16()

	flags := buf.readU16()
	// first half (MSBs) of flags
	first := flags >> 8
	// second half (LSBs)
	second := flags & 0xFF
	h.recursionDesired = (first & 1) > 0
	h.truncatedMessage = (first & (1 << 1)) > 0
	h.authoritativeAnswer = (first & (1 << 2)) > 0
	h.opcode = uint8((first >> 3) & 0xF)
	h.response = (first & (1 << 7)) > 0

	h.rescode = uint8(ResultCodeFromNumber(uint8(second & 0x0F)))
	h.checkingDisabled = (second & (1 << 4)) > 0
	h.authedData = (second & (1 << 5)) > 0
	h.z = (second & (1 << 6)) > 0
	h.recursionAvailable = (second & (1 << 7)) > 0

	h.questions = buf.readU16()
	h.answers = buf.readU16()
	h.authoritativeEntries = buf.readU16()
	h.resourceEntries = buf.readU16()
	// fmt.Printf("inside from buf: %#v", h)
}

func (h *DnsHeader) write(buf *Packet) {
	buf.writeU16(h.id)

	flags := uint8(0)
	if h.recursionDesired {
		flags |= 1
	}
	if h.truncatedMessage {
		flags |= (1 << 1)
	}
	if h.authoritativeAnswer {
		flags |= (1 << 2)
	}

	flags |= h.opcode << 3

	if h.response {
		flags |= (1 << 4)
	}

	buf.write(flags)

	//next part
	flags = 0
	flags |= h.rescode

	if h.checkingDisabled {
		flags |= (1 << 4)
	}
	if h.authedData {
		flags |= (1 << 5)
	}
	if h.z {
		flags |= (1 << 6)
	}
	if h.recursionAvailable {
		flags |= (1 << 7)
	}

	buf.writeU16(h.questions)
	buf.writeU16(h.answers)
	buf.writeU16(h.authoritativeEntries)
	buf.writeU16(h.resourceEntries)

}
