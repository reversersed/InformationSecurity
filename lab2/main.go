package main

import (
	"errors"
	"fmt"
)

const (
	base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	paddingChar = '='
)

func encodeBase64(data []byte) string {
	var result []byte
	var buffer uint32
	var i int

	for i = 0; i+2 < len(data); i += 3 {
		buffer = (uint32(data[i]) << 16) | (uint32(data[i+1]) << 8) | uint32(data[i+2])

		result = append(result, base64Chars[(buffer>>18)&0x3F])
		result = append(result, base64Chars[(buffer>>12)&0x3F])
		result = append(result, base64Chars[(buffer>>6)&0x3F])
		result = append(result, base64Chars[buffer&0x3F])
	}

	if i < len(data) {
		buffer = uint32(data[i]) << 16
		if i+1 < len(data) {
			buffer |= uint32(data[i+1]) << 8
		}

		result = append(result, base64Chars[(buffer>>18)&0x3F])
		result = append(result, base64Chars[(buffer>>12)&0x3F])

		if i+1 < len(data) {
			result = append(result, base64Chars[(buffer>>6)&0x3F])
		} else {
			result = append(result, paddingChar)
		}
		result = append(result, paddingChar)
	}

	return string(result)
}

func decodeBase64(data string) ([]byte, error) {
	var result []byte
	var buffer uint32
	var count int

	charMap := make(map[byte]uint32)
	for i, char := range base64Chars {
		charMap[byte(char)] = uint32(i)
	}

	for i := 0; i < len(data); i++ {
		char := data[i]
		if char == paddingChar {
			break
		}

		value, exists := charMap[char]
		if !exists {
			return nil, errors.New("invalid base64 character")
		}

		buffer = (buffer << 6) | value
		count++

		if count == 4 {
			result = append(result, byte((buffer>>16)&0xFF))
			result = append(result, byte((buffer>>8)&0xFF))
			result = append(result, byte(buffer&0xFF))
			buffer = 0
			count = 0
		}
	}

	if count > 0 {
		if count == 2 {
			buffer <<= 12
			result = append(result, byte((buffer>>16)&0xFF))
		} else if count == 3 {
			buffer <<= 6
			result = append(result, byte((buffer>>16)&0xFF))
			result = append(result, byte((buffer>>8)&0xFF))
		} else {
			return nil, errors.New("invalid base64 length")
		}
	}

	return result, nil
}

func main() {
	data := []byte("Hello, World!")
	encoded := encodeBase64(data)
	fmt.Println("Encoded:", encoded)

	decoded, err := decodeBase64(encoded)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	fmt.Println("Decoded:", string(decoded))
}
