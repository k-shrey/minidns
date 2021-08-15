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
		question.Read(buf)
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
