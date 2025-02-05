package main

import (
	"encoding/binary"
	"fmt"
	"math/big"
)

const blockSize = 64

var (
	pi      = []byte{252, 238, 221, 17, 207, 110, 49, 22, 251, 196, 250, 218, 35, 197, 4, 77, 233, 119, 240, 219, 147, 46, 153, 186, 23, 54, 241, 187, 20, 205, 95, 193, 249, 24, 101, 90, 226, 92, 239, 33, 129, 28, 60, 66, 139, 1, 142, 79, 5, 132, 2, 174, 227, 106, 143, 160, 6, 11, 237, 152, 127, 212, 211, 31, 235, 52, 44, 81, 234, 200, 72, 171, 242, 42, 104, 162, 253, 58, 206, 204, 181, 112, 14, 86, 8, 12, 118, 18, 191, 114, 19, 71, 156, 183, 93, 135, 21, 161, 150, 41, 16, 123, 154, 199, 243, 145, 120, 111, 157, 158, 178, 177, 50, 117, 25, 61, 255, 53, 138, 126, 109, 84, 198, 128, 195, 189, 13, 87, 223, 245, 36, 169, 62, 168, 67, 201, 215, 121, 214, 246, 124, 34, 185, 3, 224, 15, 236, 222, 122, 148, 176, 188, 220, 232, 40, 80, 78, 51, 10, 74, 167, 151, 96, 115, 30, 0, 98, 68, 26, 184, 56, 130, 100, 159, 38, 65, 173, 69, 70, 146, 39, 94, 85, 47, 140, 163, 165, 125, 105, 213, 149, 59, 7, 88, 179, 64, 134, 172, 29, 247, 48, 55, 107, 228, 136, 217, 231, 137, 225, 27, 131, 73, 76, 63, 248, 254, 141, 83, 170, 144, 202, 216, 133, 97, 32, 113, 103, 164, 45, 43, 9, 91, 203, 155, 37, 208, 190, 229, 108, 82, 89, 166, 116, 210, 230, 244, 180, 192, 209, 102, 175, 194, 57, 75, 99, 182}
	tau     = []byte{0, 8, 16, 24, 32, 40, 48, 56, 1, 9, 17, 25, 33, 41, 49, 57, 2, 10, 18, 26, 34, 42, 50, 58, 3, 11, 19, 27, 35, 43, 51, 59, 4, 12, 20, 28, 36, 44, 52, 60, 5, 13, 21, 29, 37, 45, 53, 61, 6, 14, 22, 30, 38, 46, 54, 62, 7, 15, 23, 31, 39, 47, 55, 63}
	matrixA = [64]uint64{
		0x8e20faa72ba0b470, 0x6c022c38f90a4c07, 0xa011d380818e8f40, 0x0ad97808d06cb404,
		0x90dab52a387ae76f, 0x092e94218d243cba, 0x9d4df05d5f661451, 0x18150f14b9ec46dd,
		0x86275df09ce8aaa8, 0xe230140fc0802984, 0x456c34887a3805b9, 0x9bcf4486248d9f5d,
		0xe4fa2054a80b329c, 0x492c024284fbaec0, 0x70a6a56e2440598e, 0x07e095624504536c,
		0x47107ddd9b505a38, 0x3601161cf205268d, 0x5086e740ce47c920, 0x05e23c0468365a02,
		0x486dd4151c3dfdb9, 0x8a174a9ec8121e5d, 0xc0a878a0a1330aa6, 0x0c84890ad27623e0,
		0x439da0784e745554, 0x71180a8960409a42, 0xac361a443d1c8cd2, 0xc3e9224312c8c1a0,
		0x727d102a548b194e, 0xaa16012142f35760, 0x3853dc371220a247, 0x8d70c431ac02a736,
		0xad08b0e0c3282d1c, 0x1b8e0b0e798c13c8, 0x2843fd2067adea10, 0x8c711e02341b2d01,
		0x24b86a840e90f0d2, 0x4585254f64090fa0, 0x60543c50de970553, 0x0642ca05693b9f70,
		0xafc0503c273aa42a, 0xb60c05ca30204d21, 0x561b0d22900e4669, 0xeffa11af0964ee50,
		0x39b008152acb8227, 0x550b8e9e21f7a530, 0x1ca76e95091051ad, 0xc83862965601dd1b,
		0xd8045870ef14980e, 0x83478b07b2468764, 0x14aff010bdd87508, 0x46b60f011a83988e,
		0x125c354207487869, 0xaccc9ca9328a8950, 0x302a1e286fc58ca7, 0x0321658cba93c138,
		0xd960281e9d1d5215, 0x5b068c651810a89e, 0x2b838811480723ba, 0xf97d86d98a327728,
		0x9258048415eb419d, 0xa48b474f9ef5dc18, 0x0edd37c48a08a6d8, 0x641c314b2b8ee083,
	}
	C = [12][8]int{
		{0xb1085bda, 0x1ecadae9, 0xebcb2f81, 0xc0657c1f, 0x2f6a7643, 0x2e45d016, 0x714eb88d, 0x7585c4fc},
		{0x6fa3b58a, 0xa99d2f1a, 0x4fe39d46, 0x0f70b5d7, 0xf3feea72, 0x0a232b98, 0x61d55e0f, 0x16b50131},
		{0xf574dcac, 0x2bce2fc7, 0x0a39fc28, 0x6a3d8435, 0x06f15e5f, 0x529c1f8b, 0xf2ea7514, 0xb1297b7b},
		{0xef1fdfb3, 0xe81566d2, 0xf948e1a0, 0x5d71e4dd, 0x488e857e, 0x335c3c7d, 0x9d721cad, 0x685e353f},
		{0x4bea6bac, 0xad474799, 0x9a3f410c, 0x6ca92363, 0x7f151c1f, 0x1686104a, 0x359e35d7, 0x800fffbd},
		{0xae4faeae, 0x1d3ad3d9, 0x6fa4c33b, 0x7a3039c0, 0x2d66c4f9, 0x5142a46c, 0x187f9ab4, 0x9af08ec6},
		{0xf4c70e16, 0xeeaac5ec, 0x51ac86fe, 0xbf240954, 0x399ec6c7, 0xe6bf87c9, 0xd3473e33, 0x197a93c9},
		{0x9b1f5b42, 0x4d93c9a7, 0x03e7aa02, 0x0c6e4141, 0x4eb7f871, 0x9c36de1e, 0x89b4443b, 0x4ddbc49a},
		{0x378f5a54, 0x1631229b, 0x944c9ad8, 0xec165fde, 0x3a7d3a1b, 0x25894224, 0x3cd955b7, 0xe00d0984},
		{0xabbedea6, 0x80056f52, 0x382ae548, 0xb2e4f3f3, 0x8941e71c, 0xff8a78db, 0x1fffe18a, 0x1b336103},
		{0x7bcd9ed0, 0xefc889fb, 0x3002c6cd, 0x635afe94, 0xd8fa6bbb, 0xebab0761, 0x20018021, 0x14846679},
		{0x378ee767, 0xf11631ba, 0xd21380b0, 0x0449b17a, 0xcda43c32, 0xbcdf1d77, 0xf82012d4, 0x30219f9b},
	}
)

type GOST struct {
}

func NewGost() *GOST {
	return new(GOST)
}

func (g *GOST) sBlock(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = pi[b]
	}
	return result
}

func (g *GOST) pBlock(data []byte) []byte {
	result := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		result[i] = data[tau[i]]
	}
	return result
}

func (g *GOST) lBlock(data []byte) []byte {
	result := make([]byte, blockSize)
	for i := 0; i < 8; i++ {
		chunk := binary.LittleEndian.Uint64(data[i*8 : (i+1)*8])
		var transformed uint64
		for j := 0; j < 64; j++ {
			if (chunk>>j)&1 == 1 {
				transformed ^= matrixA[j]
			}
		}
		binary.LittleEndian.PutUint64(result[i*8:(i+1)*8], transformed)
	}
	return result
}

func (g *GOST) xBlock(a, b []byte) []byte {
	result := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func (g *GOST) F(h, N []byte) []byte {
	data := g.xBlock(h, N)
	data = g.sBlock(data)
	data = g.pBlock(data)
	data = g.lBlock(data)
	return data
}

func (g *GOST) E(K, m []byte) []byte {
	keys := make([][]byte, 13)
	keys[0] = K

	for i := 1; i < 13; i++ {
		var c int64
		for _, v := range C[i-1] {
			c = (c << 8) | int64(v)
		}
		keys[i] = g.F(keys[i-1], g.intToVec512(big.NewInt(c)))
	}

	data := m
	for i := 0; i < 13; i++ {
		data = g.xBlock(keys[i], data)
		data = g.sBlock(data)
		data = g.pBlock(data)
		data = g.lBlock(data)
	}
	return data
}

func (g *GOST) g(h, m, N []byte) []byte {
	K := g.F(h, N)
	eKm := g.E(K, m)
	return g.xBlock(g.xBlock(eKm, h), m)
}

func (g *GOST) intToVec512(n *big.Int) []byte {
	vec := make([]byte, blockSize)
	n.FillBytes(vec)
	return vec
}

func (g *GOST) vec512ToInt(vec []byte) *big.Int {
	return new(big.Int).SetBytes(vec)
}

func (g *GOST) addMod512(a, b []byte) []byte {
	aInt := g.vec512ToInt(a)
	bInt := g.vec512ToInt(b)
	sum := new(big.Int).Add(aInt, bInt)
	mod := new(big.Int).Exp(big.NewInt(2), big.NewInt(512), nil)
	sum.Mod(sum, mod)
	return g.intToVec512(sum)
}

func (g *GOST) padLastBlock(M []byte) []byte {
	messageBitLength := uint64(len(M)) * 8
	zeroPaddingLength := (blockSize - (len(M)+1)%blockSize) % blockSize
	if zeroPaddingLength < 8 {
		zeroPaddingLength += blockSize
	}

	paddedBlock := make([]byte, len(M)+1+zeroPaddingLength+8)
	copy(paddedBlock, M)
	paddedBlock[len(M)] = 0x80

	binary.LittleEndian.PutUint64(paddedBlock[len(paddedBlock)-8:], messageBitLength)
	return paddedBlock
}

func (g *GOST) Sum(M []byte) []byte {
	h := make([]byte, blockSize)
	N := make([]byte, blockSize)
	Sigma := make([]byte, blockSize)

	for len(M) >= blockSize {
		m := M[len(M)-blockSize:]
		M = M[:len(M)-blockSize]

		h = g.g(h, m, N)
		N = g.addMod512(N, g.intToVec512(big.NewInt(512)))
		Sigma = g.addMod512(Sigma, m)
	}

	if len(M) > 0 {
		paddedM := g.padLastBlock(M)
		h = g.g(h, paddedM, N)
		N = g.addMod512(N, g.intToVec512(big.NewInt(int64(len(M)*8))))
		Sigma = g.addMod512(Sigma, paddedM)
	}

	h = g.g(h, N, make([]byte, blockSize))
	h = g.g(h, Sigma, make([]byte, blockSize))

	return h
}

func main() {
	gost := NewGost()
	sum := gost.Sum([]byte("323130393837363534333231303938373635343332313039383736353433323130393837363534333231303938373635343332313039383736353433323130"))
	fmt.Printf("%x\n", sum)
}
