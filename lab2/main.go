package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
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
			return nil, errors.New("найден неверный для base64 символ")
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
			return nil, errors.New("неправильная длина base64 кода")
		}
	}

	return result, nil
}

var (
	file = flag.String("file", "", "Абсолютный или относительный путь к файлу")
)

func main() {
	flag.Parse()
	if len(*file) == 0 {
		flag.PrintDefaults()
		return
	}

	bytes, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	fmt.Println("Декодирование " + strings.Split((*file), "\\")[len(strings.Split((*file), "\\"))-1] + "...")
	decoded, err := decodeBase64(string(bytes))
	if err != nil {
		fmt.Println("Невозможно декодировать файл. Применяется кодирование...")
		encoded := encodeBase64(bytes)
		newFilePath := strings.Split(path.Base(*file), ".")[0] + " enc." + strings.Split(path.Base(*file), ".")[len(strings.Split(path.Base(*file), "."))-1]
		os.Remove(newFilePath)

		if err := os.WriteFile(newFilePath, []byte(encoded), os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Println("Файл закодирован и сохранен как " + strings.Split((newFilePath), "\\")[len(strings.Split((newFilePath), "\\"))-1])
		return
	}
	newFilePath := strings.Split(path.Base(*file), ".")[0] + " dec." + strings.Split(path.Base(*file), ".")[len(strings.Split(path.Base(*file), "."))-1]
	os.Remove(newFilePath)
	if err := os.WriteFile(newFilePath, decoded, os.ModeExclusive); err != nil {
		panic(err)
	}
	fmt.Println("Файл декодирован и сохранен как " + strings.Split((newFilePath), "\\")[len(strings.Split((newFilePath), "\\"))-1])
}
