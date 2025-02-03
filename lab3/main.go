package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

const BlockSize = 8

const (
	C1 = 0x01010104
	C2 = 0x01010101
)

type gostCipher struct {
	key  [32]byte
	sbox [8][16]byte
}

type gammaMode struct {
	cipher *gostCipher
	Y, Z   uint32
}

func NewCipher(key []byte, sbox [8][16]byte) (*gostCipher, error) {
	if len(key) != 32 {
		return nil, errors.New("gost28147: invalid key size")
	}

	c := new(gostCipher)
	copy(c.key[:], key)
	c.sbox = sbox

	return c, nil
}

func (c *gostCipher) Encrypt(dst, src []byte) {
	var block [BlockSize]byte
	copy(block[:], src)

	for i := 0; i < 32; i++ {
		keyPart := c.key[i%8]
		block = c.round(block, keyPart)
	}

	copy(dst, block[:])
}

func (c *gostCipher) round(block [BlockSize]byte, key byte) [BlockSize]byte {
	var result [BlockSize]byte
	for i := 0; i < BlockSize; i++ {
		result[i] = c.sbox[i][(block[i]^key)&0x0F]
	}
	val := binary.LittleEndian.Uint64(result[:])
	val = (val << 11) | (val >> (64 - 11))
	binary.LittleEndian.PutUint64(result[:], val)
	return result
}

func NewGammaMode(c *gostCipher, S [BlockSize]byte) *gammaMode {
	g := new(gammaMode)
	g.cipher = c

	g.Y = binary.LittleEndian.Uint32(S[:4])
	g.Z = binary.LittleEndian.Uint32(S[4:])

	return g
}

func addModulo(a, b uint32) uint32 {
	res := a + b
	if res < a || res < b {
		res++
	}
	return res
}

func (g *gammaMode) Crypt(dst, src []byte) {
	for i := 0; i < len(src); i += BlockSize {
		gamma := g.generateGamma()

		for j := 0; j < BlockSize && i+j < len(src); j++ {
			dst[i+j] = src[i+j] ^ gamma[j]
		}
	}
}

func (g *gammaMode) generateGamma() [BlockSize]byte {
	Yj := addModulo(g.Y, C2)
	Zj := addModulo(g.Z, C1)

	var gamma [BlockSize]byte
	binary.LittleEndian.PutUint32(gamma[:4], Yj)
	binary.LittleEndian.PutUint32(gamma[4:], Zj)
	g.cipher.Encrypt(gamma[:], gamma[:])

	g.Y = Yj
	g.Z = Zj

	return gamma
}

func EncryptWithSync(cipher *gostCipher, plaintext []byte, S [BlockSize]byte) []byte {
	ciphertext := make([]byte, BlockSize+len(plaintext))

	copy(ciphertext[:BlockSize], S[:])

	gmode := NewGammaMode(cipher, S)
	gmode.Crypt(ciphertext[BlockSize:], plaintext)

	return ciphertext
}

func DecryptWithSync(cipher *gostCipher, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	S := [BlockSize]byte{}
	copy(S[:], ciphertext[:BlockSize])

	plaintext := make([]byte, len(ciphertext)-BlockSize)
	gmode := NewGammaMode(cipher, S)
	gmode.Crypt(plaintext, ciphertext[BlockSize:])

	return plaintext, nil
}

var (
	file = flag.String("file", "", "Абсолютный или относительный путь к файлу")
	enc  = flag.Bool("enc", false, "Включить режим кодирования")
)

func main() {
	key := make([]byte, 32)
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

	fmt.Printf("Ключ шифрования: %v\n", key)
	if *enc {
		fmt.Println("Кодирование " + strings.Split((*file), "\\")[len(strings.Split((*file), "\\"))-1] + "...")
		encoded := EncryptWithSync(cipher, bytes, S)
		newFilePath := strings.Split(path.Base(*file), ".")[0] + " enc." + strings.Split(path.Base(*file), ".")[len(strings.Split(path.Base(*file), "."))-1]
		os.Remove(newFilePath)

		if err := os.WriteFile(newFilePath, []byte(encoded), os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Println("Файл закодирован и сохранен как " + strings.Split((newFilePath), "\\")[len(strings.Split((newFilePath), "\\"))-1])
	} else {
		fmt.Println("Декодирование " + strings.Split((*file), "\\")[len(strings.Split((*file), "\\"))-1] + "...")
		decrypted, err := DecryptWithSync(cipher, bytes)
		if err != nil {
			panic(err)
		}
		newFilePath := strings.Split(path.Base(*file), ".")[0] + " dec." + strings.Split(path.Base(*file), ".")[len(strings.Split(path.Base(*file), "."))-1]
		os.Remove(newFilePath)
		if err := os.WriteFile(newFilePath, decrypted, os.ModeExclusive); err != nil {
			panic(err)
		}
		fmt.Println("Файл декодирован и сохранен как " + strings.Split((newFilePath), "\\")[len(strings.Split((newFilePath), "\\"))-1])
	}
}
