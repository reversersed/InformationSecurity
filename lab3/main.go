package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	BlockSize = 8
	C1        = 0x01010104
	C2        = 0x01010101
)

type gostCipher struct {
	key  [32]byte
	sbox [8][16]byte
}

func NewCipher(key []byte, sbox [8][16]byte) (*gostCipher, error) {
	if len(key) != 32 {
		return nil, errors.New("gost28147: invalid key size: " + strconv.Itoa(len(key)))
	}

	c := new(gostCipher)
	copy(c.key[:], key)
	c.sbox = sbox

	return c, nil
}
func (c *gostCipher) encryptBlock(block [BlockSize]byte, subkey uint32) [BlockSize]byte {
	left := binary.LittleEndian.Uint32(block[:4])
	right := binary.LittleEndian.Uint32(block[4:])

	rightWithKey := right + subkey

	var sboxResult uint32
	for i := 0; i < 4; i++ {
		byteVal := byte(rightWithKey >> (8 * i))
		sboxResult |= uint32(c.sbox[i][byteVal&0x0F]) << (8 * i)
	}

	sboxResult = (sboxResult << 11) | (sboxResult >> (32 - 11))

	newRight := left ^ sboxResult

	left = right
	right = newRight

	var result [BlockSize]byte
	binary.LittleEndian.PutUint32(result[:4], left)
	binary.LittleEndian.PutUint32(result[4:], right)

	return result
}

func (c *gostCipher) EncryptGamma(plaintext []byte, S [BlockSize]byte) []byte {
	Y := binary.LittleEndian.Uint32(S[:4])
	Z := binary.LittleEndian.Uint32(S[4:])
	var gamma [BlockSize]byte

	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += BlockSize {
		gamma, Y, Z = c.generateGamma(Y, Z)

		for j := 0; j < BlockSize && i+j < len(plaintext); j++ {
			ciphertext[i+j] = plaintext[i+j] ^ gamma[j]
		}
	}

	return ciphertext
}

func (c *gostCipher) generateGamma(Y, Z uint32) ([BlockSize]byte, uint32, uint32) {
	Yj := Y + C2
	Zj := (Z+C1-1)%(1<<32-1) + 1

	var gamma [BlockSize]byte
	binary.LittleEndian.PutUint32(gamma[:4], Yj)
	binary.LittleEndian.PutUint32(gamma[4:], Zj)
	gamma = c.encryptBlock(gamma, 0)

	return gamma, Yj, Zj
}

func (c *gostCipher) DecryptGamma(ciphertext []byte, S [BlockSize]byte) []byte {
	return c.EncryptGamma(ciphertext, S)
}

var (
	file = flag.String("file", "", "Абсолютный или относительный путь к файлу")
	dec  = flag.Bool("dec", false, "Включить режим декодирования")
)

func main() {
	key := []byte{255, 126, 235, 54, 45, 27, 15, 69, 228, 14, 88, 148, 8, 91, 99, 42, 52, 54, 12, 65, 24, 55, 127, 246, 126, 109, 195, 121, 12, 5, 0, 8}
	S := [BlockSize]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}

	sbox := [8][16]byte{
		{4, 10, 9, 2, 13, 8, 0, 14, 6, 11, 1, 12, 7, 15, 5, 3},
		{14, 11, 4, 12, 6, 13, 15, 10, 2, 3, 8, 1, 0, 7, 5, 9},
		{5, 8, 1, 13, 10, 3, 4, 2, 14, 15, 12, 7, 6, 0, 9, 11},
		{7, 13, 10, 1, 0, 8, 9, 15, 14, 4, 6, 12, 11, 2, 5, 3},
		{6, 12, 7, 1, 5, 15, 13, 8, 4, 10, 9, 14, 0, 3, 11, 2},
		{4, 11, 10, 0, 7, 2, 1, 13, 3, 6, 8, 5, 9, 12, 15, 14},
		{13, 11, 4, 1, 3, 15, 5, 9, 0, 10, 14, 7, 6, 8, 2, 12},
		{1, 15, 13, 0, 5, 7, 10, 4, 9, 2, 3, 14, 6, 11, 8, 12},
	}

	cipher, err := NewCipher(key, sbox)
	if err != nil {
		panic(err)
	}

	flag.Parse()
	if len(*file) == 0 {
		flag.PrintDefaults()
		return
	}

	bytes, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	if !*dec {
		fmt.Println("Кодирование " + strings.Split((*file), "\\")[len(strings.Split((*file), "\\"))-1] + "...")
		encoded := cipher.EncryptGamma(bytes, S)
		newFilePath := strings.Split(path.Base(*file), ".")[0] + " enc." + strings.Split(path.Base(*file), ".")[len(strings.Split(path.Base(*file), "."))-1]
		os.Remove(newFilePath)

		if err := os.WriteFile(newFilePath, []byte(encoded), os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Println("Файл закодирован и сохранен как " + strings.Split((newFilePath), "\\")[len(strings.Split((newFilePath), "\\"))-1])
	} else {
		fmt.Println("Декодирование " + strings.Split((*file), "\\")[len(strings.Split((*file), "\\"))-1] + "...")
		decrypted := cipher.DecryptGamma(bytes, S)

		newFilePath := strings.Split(path.Base(*file), ".")[0] + " dec." + strings.Split(path.Base(*file), ".")[len(strings.Split(path.Base(*file), "."))-1]
		os.Remove(newFilePath)
		if err := os.WriteFile(newFilePath, decrypted, os.ModeExclusive); err != nil {
			panic(err)
		}
		fmt.Println("Файл декодирован и сохранен как " + strings.Split((newFilePath), "\\")[len(strings.Split((newFilePath), "\\"))-1])
	}
}
