package main

// import "fmt"

type DnsPacket struct {
	header      DnsHeader
	questions   []DnsQuestion
	answers     []DnsRecord
	authorities []DnsRecord
	resources   []DnsRecord
}

func NewDnsPacket() DnsPacket {
	dnsPacket := DnsPacket{
		header:      DnsHeader{},
		questions:   []DnsQuestion{},
		answers:     []DnsRecord{},
		authorities: []DnsRecord{},
		resources:   []DnsRecord{},
	}
	return dnsPacket
}

func FromBuffer(buf *Packet) DnsPacket {
	result := NewDnsPacket()
	result.header.read(buf)
	// fmt.Printf("inside from buf: %#v", result.header)
	for i := 0; i < int(result.header.questions); i++ {
		question := NewDnsQuestion("", UNKNOWN)
		question.read(buf)
		result.questions = append(result.questions, question)
	}

	// by, _ := buf.get(buf.pos + 1)

	// fmt.Printf("now next byte is: %#x", by)
	for i := 0; i < int(result.header.answers); i++ {
		answers := ReadDnsRecord(buf)
		result.answers = append(result.answers, answers)
	}

	for i := 0; i < int(result.header.authoritativeEntries); i++ {
		authority := ReadDnsRecord(buf)
		result.authorities = append(result.authorities, authority)
	}

	for i := 0; i < int(result.header.resourceEntries); i++ {
		resource := ReadDnsRecord(buf)
		result.resources = append(result.resources, resource)
	}

	return result
}

func (d *DnsPacket) toBuffer(buf *Packet) {
	d.header.questions = uint16(len(d.questions))
	d.header.answers = uint16(len(d.answers))
	d.header.authoritativeEntries = uint16(len(d.authorities))
	d.header.resourceEntries = uint16(len(d.resources))

	d.header.write(buf)

	for _, q := range d.questions {
		q.write(buf)
	}
	for _, a := range d.answers {
		a.writeDnsRecord(buf)
	}
	for _, a := range d.authorities {
		a.writeDnsRecord(buf)
	}
	for _, r := range d.resources {
		r.writeDnsRecord(buf)
	}

}
