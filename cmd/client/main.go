package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/x-junkang/connected/internal/configure"
	"github.com/x-junkang/connected/internal/protocol"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", configure.GlobalObject.Host, configure.GlobalObject.TCPPort))
	if err != nil {
		fmt.Println("连接出错")
		return
	}
	header := &protocol.MarsHeader{
		HeaderLength: 20,
		Sequence:     1,
		BodyLength:   5,
	}
	data := []byte{'h', 'e', 'l', 'l', 'o'}
	for {
		binary.Write(conn, binary.LittleEndian, header)
		binary.Write(conn, binary.LittleEndian, data)

		fmt.Println("done 1")
		var respHead protocol.MarsHeader
		err = binary.Read(conn, binary.LittleEndian, &respHead)
		if err != nil {
			return
		}
		bodyLen := int(header.BodyLength)
		resp := make([]byte, bodyLen)
		n, err := io.ReadFull(conn, resp)
		if err != nil {
			return
		}
		fmt.Println(n, string(resp))
		time.Sleep(1 * time.Second)
	}
}
