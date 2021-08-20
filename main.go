package main

import (
	// "bufio"
	"fmt"
	"net"
)

func main() {

	qname := "google.com"
	qtype := A
	// buffer := make([]byte, 512)

	packet := NewDnsPacket()
	packet.header.id = 12345
	packet.header.questions = 1
	packet.header.recursionDesired = true
	packet.questions = append(packet.questions, NewDnsQuestion(qname, qtype))

	reqBuf := NewPacket()
	packet.toBuffer(&reqBuf)

	// fmt.Printf("%#v\n", packet)
	udpAddr, err := net.ResolveUDPAddr("udp4", "8.8.8.8:53")
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	fmt.Println("got google connection, writing now")
	_, err = conn.Write(reqBuf.buf[0:reqBuf.pos])

	fmt.Print(reqBuf.buf[0:reqBuf.pos])
	resBuf := NewPacket()
	fmt.Println("attempting to read packet")
	// _, err = bufio.NewReader(conn).Read(resBuf.buf[:])
	conn.Read(resBuf.buf[:])

	fmt.Printf("%#v\n", resBuf.buf[:])
	resPacket := FromBuffer(&resBuf)

	// fmt.Println(resBuf.buf)

	for _, q := range resPacket.questions {
		fmt.Printf("\nquestions: %#v\n", q)
	}
	// fmt.Println("num answers", len(packet.answers))
	for _, a := range resPacket.answers {
		fmt.Printf("\nanswers: %#v\n", a)
	}
	for _, au := range resPacket.authorities {
		fmt.Println("authorities : ", au)
	}
	for _, r := range resPacket.resources {
		fmt.Println("resources : ", r)
	}

	// addr := net.UDPAddr{
	// 	Port: 45000,
	// 	IP:   net.ParseIP("127.0.0.1"),
	// }
	// ser, err := net.ListenUDP("udp", &addr)

	// if err != nil {
	// 	fmt.Printf("Error, cant list on port%v\n", err)
	// 	return
	// }

	// for {
	// 	_, remoteAddr, err := ser.ReadFromUDP(buffer)

	// }
	// file, err := ioutil.ReadFile("response_packet.txt")
	// if err != nil {
	// 	fmt.Println("cannot read file")
	// }

	// buffer := NewPacket()
	// copy(buffer.buf[:], file)

	// packet := FromBuffer(&buffer)
	// fmt.Printf("header: %#v\n", packet.header)

	// for _, q := range packet.questions {
	// 	fmt.Printf("\nquestions: %#v\n", q)
	// }
	// // fmt.Println("num answers", len(packet.answers))
	// for _, a := range packet.answers {
	// 	fmt.Printf("\nanswers: %#v\n", a)
	// }
	// for _, au := range packet.authorities {
	// 	fmt.Println("authorities : ", au)
	// }
	// for _, r := range packet.resources {
	// 	fmt.Println("resources : ", r)
	// }

}
