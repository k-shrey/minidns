package main

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

func (d *DnsQuestion) read(buf *Packet) {
	d.name, _ = buf.readQName()
	// fmt.Println("d.name is", d.name)
	d.qtype = QueryType(buf.readU16())
	// fmt.Println("d.qtype is", d.qtype)
	buf.readU16()
	// fmt.Println("class is", class)
}

func (d *DnsQuestion) write(buf *Packet) {
	buf.writeQName(d.name)

	typenum := d.qtype
	buf.writeU16(uint16(typenum))
	buf.writeU16(1)
}
