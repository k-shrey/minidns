package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {

	dnsCache := NewDnsCache(60)
	// qname := "www.yahoo.com"
	addr := net.UDPAddr{
		Port: 1053,
		IP:   net.ParseIP("127.0.0.1"),
	}

	fmt.Println("listening on udp")
	ln, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	// defer ln.Close()

	upstream, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	var wg sync.WaitGroup
	for {

		var buf [512]uint8
		fmt.Println("waiting to read from udp")
		_, addr, err := ln.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("cant read from udp??")
			log.Fatal(err)
		}
		fmt.Println("read something.")
		fmt.Println("Starting worker")
		wg.Add(1)
		go worker(&wg, true, ln, addr, buf, &upstream, dnsCache)
	}

	// wg.Wait()
}

func worker(wg *sync.WaitGroup, debug bool, conn *net.UDPConn, addr *net.UDPAddr,
	buf [512]uint8, upstream *net.Conn, cache *DnsCache) {

	defer wg.Done()

	recBuf := NewPacket()
	recBuf.buf = buf
	// _, err := bufio.NewReader(*conn).Read(recBuf.buf[:])
	// if err != nil {
	// 	fmt.Printf("Some error while reading %v", err)
	// 	return
	// }
	recPacket := FromBuffer(&recBuf)
	fmt.Printf("Decoded the packet. making another.\n")

	packet := NewDnsPacket()
	packet.header.id = recPacket.header.id
	packet.header.questions = recPacket.header.questions
	packet.header.recursionDesired = true
	packet.header.z = false

	fmt.Printf("checking cache\n")

	flag := 0
	for _, q := range recPacket.questions {
		cache.mutex.Lock()
		if cache.cache[q.name] != nil && cache.cache[q.name].value != nil {
			packet.questions = append(packet.questions, NewDnsQuestion(q.name, A))
			flag = 1
		}
		cache.mutex.Unlock()
	}
	fmt.Printf("not present in the cache.\n")

	if flag == 1 {
		resBuf := NewPacket()
		packet.toBuffer(&resBuf)
		_, err := conn.WriteToUDP(resBuf.buf[0:resBuf.pos], addr)
		if err != nil {
			fmt.Printf("Some error while writing back to client%v", err)
			return
		}
		return
	}

	// packet := NewDnsPacket()
	// packet.header.id = uint16(rand.Intn(10000))
	// packet.header.questions = 1
	// packet.header.recursionDesired = true
	// packet.header.z = false
	// packet.questions = append(packet.questions, NewDnsQuestion(qname, A))

	// reqBuf := NewPacket()
	// packet.toBuffer(&reqBuf)

	num, err := (*upstream).Write(recBuf.buf[0:recBuf.pos])
	if err != nil {
		fmt.Printf("Some error while writing%v", err)
		return
	}

	fmt.Printf("Wrote %#v bytes\n", num)
	// fmt.Printf("%#v", reqBuf.buf[0:reqBuf.pos])
	resBuf := NewPacket()

	num, err = bufio.NewReader(*upstream).Read(resBuf.buf[:])
	if err != nil {
		fmt.Printf("Some error while reading %v", err)
		return
	}
	fmt.Printf("Read %#v bytes\n", num)

	resPacket := FromBuffer(&resBuf)
	for _, a := range resPacket.answers {
		if domain := a.ARecord.domain; domain != "" {
			cache.mutex.Lock()
			cache.cache[domain] = &item{}
			cache.cache[domain].value = a.ARecord.addr
			cache.cache[domain].lastAccess = time.Now().Unix()
			cache.mutex.Unlock()
		}
	}

	_, err = conn.WriteToUDP(resBuf.buf[0:resBuf.pos], addr)
	if err != nil {
		fmt.Printf("Some error while writing back to client 2 %v", err)
		return
	}

	if debug {
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

}
