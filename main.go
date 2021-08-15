package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	file, err := ioutil.ReadFile("response_packet.txt")
	if err != nil {
		fmt.Println("cannot read file")
	}

	buffer := NewPacket()
	copy(buffer.buf[:], file)

	packet := FromBuffer(&buffer)
	fmt.Printf("header: %#v\n", packet.header)

	for _, q := range packet.questions {
		fmt.Printf("\nquestions: %#v\n", q)
	}
	// fmt.Println("num answers", len(packet.answers))
	for _, a := range packet.answers {
		fmt.Printf("\nanswers: %#v\n", a)
	}
	for _, au := range packet.authorities {
		fmt.Println("authorities : ", au)
	}
	for _, r := range packet.resources {
		fmt.Println("resources : ", r)
	}

}
