package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	// Specify the NTP server to connect to
	// ntpServer := "time.google.com"
	ntpServer := "localhost"

	conn, err := net.Dial("udp", ntpServer+":10123")
	if err != nil {
		fmt.Println("Error connecting to NTP server:", err)
		return
	}
	defer conn.Close()

	ntpReq := make([]byte, 48)

	ntpReq[0] = 0x1B

	for {
		var responses [8]time.Time
		for i := 0; i < 8; i++ {
			// Send the NTP request to the server
			_, err = conn.Write(ntpReq)
			if err != nil {
				fmt.Println("Error sending NTP request:", err)
				return
			}
			time.Sleep(2 * time.Second)
		}

		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("Error setting read deadline:", err)
			return
		}
		for i := 0; i < 8; i++ {
			ntpResp := make([]byte, 48)
			_, err = conn.Read(ntpResp)
			if err != nil {
				fmt.Printf("Error receiving NTP response-%v:%v", i, err)
				break
			}

			fmt.Printf("Buffer Received -> %v \n", ntpResp)

			seconds := binary.BigEndian.Uint32(ntpResp[40:44])
			fraction := binary.BigEndian.Uint32(ntpResp[44:48])
			fmt.Printf("Seconds-> %v, Fraction %v\n", seconds, fraction)
			ntpTime := time.Unix(int64(seconds)-2208988800, int64(fraction)/10)

			responses[i] = ntpTime
		}

		fmt.Println("NTP times:", responses)
		time.Sleep(1 * time.Minute)
	}
}
