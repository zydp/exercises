package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Info struct {
	A uint16
	B int8
	C uint8
	//D string
}

/* Package binary implements simple translation between numbers and byte sequences and encoding and decoding of varints. */

func main() {
	writebuf := new(bytes.Buffer)
	data := Info{A: 61374, B: -54, C: 254}
	//data.D = "TestMessage"
	err := binary.Write(writebuf, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return
	}

	fmt.Printf("%x\n", writebuf.Bytes())

	readerbuf := bytes.NewReader(writebuf.Bytes())
	var tpro Info
	err = binary.Read(readerbuf, binary.LittleEndian, &tpro)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(tpro)
}
