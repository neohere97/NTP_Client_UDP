package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	ntpPort = 10123
)

func main() {
	fmt.Println("Starting NTP server...")
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", ntpPort))
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, 48)
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		buf[0] = 0x1B

		ntpTime := time.Now().UnixNano()
		fmt.Printf("\n%v ntpTime fraction \n", ntpTime)
		binary.BigEndian.PutUint32(buf[40:], uint32(ntpTime/1e9)+2208988800)
		binary.BigEndian.PutUint32(buf[44:], uint32(ntpTime%(ntpTime/1e9)*10))

		_, err = conn.WriteToUDP(buf, addr)
		fmt.Printf("\nBufferSent-> %v", buf)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			continue
		}
	}
}

func calculateChecksum(buf []byte) uint16 {
	var sum uint32
	for i := 0; i < len(buf); i += 2 {
		sum += uint32(buf[i])<<8 | uint32(buf[i+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return uint16(^sum)
}
