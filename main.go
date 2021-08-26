package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
)

func main() {

	qname := "www.yahoo.com"
	qtype := A
	// buffer := make([]uint8, 512)

	packet := NewDnsPacket()
	packet.header.id = 12345
	packet.header.questions = 1
	packet.header.recursionDesired = true
	packet.header.z = true
	packet.questions = append(packet.questions, NewDnsQuestion(qname, qtype))

	reqBuf := NewPacket()
	packet.toBuffer(&reqBuf)

	// pwd, _ := os.Getwd()
	// fmt.Println("DIRECOTY IS", pwd)
	e := ioutil.WriteFile("/packet.dmp", reqBuf.buf[:], 0644)
	if e != nil {
		panic(e)
	}

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	// fmt.Println("got google connection, writing now")
	num, err := conn.Write(reqBuf.buf[0:reqBuf.pos])
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Printf("Wrote %#v bytes\n", num)
	fmt.Printf("%#v", reqBuf.buf[0:reqBuf.pos])
	resBuf := NewPacket()

	num, err = bufio.NewReader(conn).Read(resBuf.buf[:])
	fmt.Printf("Read %#v bytes\n", num)

	resPacket := FromBuffer(&resBuf)

	for _, q := range resPacket.questions {
		fmt.Printf("\nquestions: %#v\n", q)
	}
	fmt.Println("num answers", len(resPacket.answers))
	for _, a := range resPacket.answers {
		fmt.Printf("\nanswers: %#v\n", a)
	}
	for _, au := range resPacket.authorities {
		fmt.Println("authorities : ", au)
	}
	for _, r := range resPacket.resources {
		fmt.Println("resources : ", r)
	}

}
