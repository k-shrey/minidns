package main

// import "fmt"

type DnsQuestion struct {
	name  string
	qtype QueryType
}

func NewDnsQuestion(name string, qtype QueryType) DnsQuestion {
	q := DnsQuestion{
		name:  name,
		qtype: qtype,
	}

	return q
}

func (d *DnsQuestion) Read(buf *Packet) {
	d.name, _ = buf.readQName()
	d.qtype = QueryType(buf.readU16())
	buf.readU16()
	// fmt.Println(k)
}

func (d *DnsQuestion) Write(buf *Packet) {
	buf.writeQName(d.name)
	
	typenum := d.qtype
	buf.writeU16(uint16(typenum))
	buf.writeU16(1)
}
