package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	burst := 0
	messagePair := 0
	for {
		for i := 0; i < 8; i++ {
			go sendReqPacket(burst, messagePair)
			time.Sleep(2 * time.Second)
			messagePair += 1
		}
		burst += 1
		messagePair = 0
		time.Sleep(4 * time.Minute)
	}
}

func sendReqPacket(burst int, messagePair int) {

	ntpServer := "localhost"
	// ntpServer := "time.google.com"
	// ntpServer := "34.69.18.67"
	conn, err := net.Dial("udp", ntpServer+":10123")
	if err != nil {
		fmt.Println("Error connecting to NTP server:", err)
		return
	}
	defer conn.Close()

	ntpReq := make([]byte, 48)

	ntpReq[0] = 0x1B
	orgTime := time.Now().UnixNano()
	binary.BigEndian.PutUint32(ntpReq[24:], uint32(orgTime/1e9)+2208988800)
	binary.BigEndian.PutUint32(ntpReq[28:], uint32(orgTime%(orgTime/1e9)*10))

	_, err = conn.Write(ntpReq)
	if err != nil {
		fmt.Println("Error sending NTP request:", err)
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		fmt.Println("Error setting read deadline:", err)
		return
	}

	ntpResp := make([]byte, 48)
	_, err = conn.Read(ntpResp)
	t4 := time.Now().UnixNano()
	if err != nil {
		fmt.Printf("Error receiving NTP response-%v", err)
	}

	recvSeconds := binary.BigEndian.Uint32(ntpResp[32:36])
	recvFraction := binary.BigEndian.Uint32(ntpResp[36:40])
	recvNanoSeconds := (int64(recvSeconds)-2208988800)*1e9 + int64(recvFraction)/10

	seconds := binary.BigEndian.Uint32(ntpResp[40:44])
	fraction := binary.BigEndian.Uint32(ntpResp[44:48])
	nanoSeconds := (int64(seconds)-2208988800)*1e9 + int64(fraction)/10

	T1 := orgTime
	T2 := recvNanoSeconds
	T3 := nanoSeconds
	T4 := t4

	// fmt.Printf("\n T1 -> %v ", time.Unix(0, orgTime))
	// fmt.Printf("\n T2 -> %v ", recvTime)
	// fmt.Printf("\n T3-> %v ", ntpTime)
	// fmt.Printf("\n T4-> %v", time.Unix(0, t4))

	fmt.Printf("\n BurstNo,%v,MessageNo,%v,Delay,%v,Offset,%v", burst, messagePair, float32((float32(T4-T1)-float32(T3-T2))/1e6), 0.5*float32((float32(T2-T1)+float32(T3-T4))/1e6))
}
