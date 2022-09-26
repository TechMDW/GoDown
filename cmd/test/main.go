package main

import "github.com/TechMDW/GoDown/internal/attacks"

func main() {
	attacks.UdpFlood()

	select {}
}
