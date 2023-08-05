package parser

import (
	"errors"
	"fmt"
	"os"
	"sort"
)

const PACKET_SIZE = 188
// This is the decimal of 0x47
const SYNC_BYTE = 71
const HEADER_BYTE_SIZE = 4

/*
Steps to ParseMPEG
- group data by packets 188 bytes
- check if sync byte exists if it equals to 0x47
- check PID:
	- get 5 bits from second byte and last 8 bits from thirdbyte
*/
func ParseMPEG(byteList []byte) {
	firstPacketIgnore := true
	dataStreamPacket := dataStreamToPackets(byteList)
	pidByteHashTable := make(map[uint16]int)
	var pidArray []uint16
	// We need to handle the first packet to see if the sync byte 
	// is in the middle
	firstPacket := dataStreamPacket[0]
	firstPacketSyncByteIndex := locateSyncByte(firstPacket)
	if firstPacketSyncByteIndex != -1 {
		firstPacket = firstPacket[firstPacketSyncByteIndex:]
	}
	// If the packet is incomplete then we will ignore
	if (len(firstPacket) < 4) {
		firstPacketIgnore = true
	}
	// This will check packets sync byte
	checkPacketsSyncByte(dataStreamPacket, firstPacketIgnore)
	for i, packet := range dataStreamPacket {
		result, err := getPacketPid(packet)
		if err != nil {
			fmt.Printf("Not enough bytes in packet %d", i)
		}
		pidByteHashTable[result] = pidByteHashTable[result] + 1
	}
	for key := range pidByteHashTable {
		pidArray = append(pidArray, key)
	}
	sort.Slice(pidArray, func(i, j int) bool {
		return pidArray[i] < pidArray[j]
	})
	for _, pid := range pidArray {
		formattedHex := fmt.Sprintf("0x%X", pid)
		fmt.Printf(formattedHex)
		fmt.Println("")
	}
}

/*
This will take in the whole data stream and place them into PACKET_SIZE
*/
func dataStreamToPackets(dataStream []byte)[][]byte{
	numPackets := (len(dataStream) + PACKET_SIZE - 1) / PACKET_SIZE
	packets := make([][]byte, numPackets)
	for i := 0; i < numPackets; i++ {
		start := i * PACKET_SIZE
		end := (i + 1) * PACKET_SIZE
		if end > len(dataStream) {
			end = len(dataStream)
		}
		packets[i] = dataStream[start:end]
	}

	return packets
}

/*
The purpose of this function is to locate the sync byte in a packet
if it finds the sync byte it will return the index, if not it will
return -1
*/
func locateSyncByte(packet []byte) int {
	for i, currByte := range packet{
		if currByte == 71 {
			return i
		}
	}
	return -1
}

/*
This will loop through the packets to see if all packets have a sync byte
there is also a parameter for firstPacketIgnore, and if set, it will discard
the first packet. If there is a packet with no sync byte it will kill the
program and display the message
*/
func checkPacketsSyncByte (packetArray [][]byte, firstPacketIgnore bool) {
	for i, packet := range packetArray {
		if (firstPacketIgnore) && (i == 0) {
			continue
		}
		if !isSyncBytePresent(packet) {
			fmt.Printf("Error: No sync byte present in packet %d, offset %d", i, (i * PACKET_SIZE))
			os.Exit(1)
		}
	}
}

/*
This will check if a syncbyte is present in a packet
*/
func isSyncBytePresent(packet []byte) bool {
	if packet[0] == SYNC_BYTE {
		return true
	}
	return false
}

/*
The purpose of this function is retrieve the pid of the packets
*/
func getPacketPid (packet []byte) (uint16,error) {
	// error check, it means that not enough bytes for a header
	if len(packet) < HEADER_BYTE_SIZE {
		return 0, errors.New("Not enough Bytes")
	}

	pid := uint16(packet[1]&0x1F)<<8 | uint16(packet[2])

	// pidByte := byte(pid & 0xFF)
	return pid, nil
}