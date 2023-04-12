package levis

import (
	"encoding/binary"
	"fmt"
)

func Equal(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func DecodeToBytes(s []uint16) []byte {

	result := make([]byte, 0)

	for _, v := range s {
		result = append(result, byte(v>>8&0xFF))
		result = append(result, byte(v&0xFF))
	}

	fmt.Printf("debug DecodeToBytes: %s, [% X]\n", result, result)
	return result
}

func DecodeToChars(s []uint16) []byte {

	result := make([]byte, 0)

	for _, v := range s {
		result = append(result, byte(v&0xFF))
		result = append(result, byte(v>>8&0xFF))
	}

	fmt.Printf("debug DecodeToChars: %s, [% X]\n", result, result)
	return result
}

func EncodeFromChars(s []byte) []uint16 {

	copyS := make([]byte, len(s))

	copy(copyS, s)

	if len(copyS)%2 != 0 {
		copyS = append(copyS, 0x00)
	}

	result := make([]uint16, 0)

	for i := range make([]int, len(copyS)/2) {
		idx := 2 * i
		value := []byte{copyS[idx], copyS[idx+1]}
		result = append(result, binary.LittleEndian.Uint16(value))
	}

	fmt.Printf("debug EncodeFromChars: %v\n", result)
	return result
}

func EncodeFromBytes(s []byte) []uint16 {

	copyS := make([]byte, len(s))

	copy(copyS, s)

	if len(copyS)%2 != 0 {
		copyS = append(copyS, 0x00)
	}

	result := make([]uint16, 0)

	for i := range make([]int, len(copyS)/2) {
		idx := 2 * i
		value := []byte{copyS[idx], copyS[idx+1]}
		result = append(result, binary.BigEndian.Uint16(value))
	}

	fmt.Printf("debug EncodeFromBytes: %v\n", result)
	return result
}

func EncodeToChars(s []byte) []byte {

	copyS := make([]byte, 0)
	copyS = append(copyS, s...)
	if len(copyS)%2 != 0 {
		copyS = append(copyS, 0x00)
	}

	// fmt.Printf("debug: %s, [% X]\n", copyS, copyS)

	result := make([]byte, 0)

	for {
		if len(copyS) < 2 {
			break
		}
		result = append(result, copyS[1])
		result = append(result, copyS[0])

		copyS = copyS[2:]
	}

	fmt.Printf("debug EncodeToChars: %v\n", result)
	return result

}

func GenerateMessage(slaveId int, funcCode int, data []byte) []byte {

	adu := make([]byte, 0x00)

	adu = append(adu, byte(slaveId))
	adu = append(adu, byte(funcCode))
	adu = append(adu, byte(len(data)))
	adu = append(adu, data...)

	crcc := new(crc)
	crcc.reset().pushBytes(adu)
	checksum := crcc.value()

	fmt.Printf("encode crc: %X\n", checksum)

	adu = append(adu, byte(checksum))
	adu = append(adu, byte(checksum>>8))

	length := len(adu)

	crcc.reset().pushBytes(adu[0 : length-2])
	checksum_ := uint16(adu[length-1])<<8 | uint16(adu[length-2])

	fmt.Printf("decode crc: %X\n", checksum_)

	return adu
}
