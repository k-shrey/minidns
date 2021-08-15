package main

import (
	// "fmt"
	"fmt"
	"net"
)

type ARecord struct {
	domain string
	addr   net.IP
	ttl    uint32
}

type UNKNOWNRecord struct {
	domain     string
	qType      QueryType
	dataLength uint16
	ttl        uint32
}

type DnsRecord struct {
	RecordType    int
	ARecord       ARecord
	UNKNOWNRecord UNKNOWNRecord
}

func ReadDnsRecord(buf *Packet) DnsRecord {
	domain, _ := buf.readQName()
	qType := QueryType(buf.readU16())
	buf.readU16()
	ttl := buf.readU32()
	dataLength := buf.readU16()

	var record DnsRecord
	switch qType {
	case 1:
		rawAddr := buf.readU32()
		addr := net.IPv4(
			uint8((rawAddr>>24)&0xFF),
			uint8((rawAddr>>16)&0xFF),
			uint8((rawAddr>>8)&0xFF),
			uint8((rawAddr>>0)&0xFF),
		)
		record = DnsRecord{
			ARecord: ARecord{
				domain: domain,
				addr:   addr,
				ttl:    ttl,
			},
			RecordType: 1,
		}
	case 0:
		buf.step(int(dataLength))
		record = DnsRecord{
			UNKNOWNRecord: UNKNOWNRecord{
				domain:     domain,
				qType:      qType,
				ttl:        ttl,
				dataLength: dataLength,
			},
			RecordType: 0,
		}
	}

	return record
}

func (d *DnsRecord) WriteDnsRecord(buf *Packet) int {
	startPos := buf.pos

	if d.RecordType == 1 {
		buf.writeQName(d.ARecord.domain)
		buf.writeU16(1)
		buf.writeU32(d.ARecord.ttl)
		buf.writeU16(4)

		ip := d.ARecord.addr.To4()

		buf.write(ip[0])
		buf.write(ip[1])
		buf.write(ip[2])
		buf.write(ip[3])

	} else if d.RecordType == 0 {
		fmt.Println("skipping unknown record")
	}

	return buf.pos - startPos
}
