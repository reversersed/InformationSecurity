package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func generateKeys(p *big.Int, g *big.Int) (*big.Int, *big.Int) {
	cj, _ := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
	cj.Add(cj, big.NewInt(1))

	dj := new(big.Int).Exp(g, cj, p)

	return cj, dj
}

func encrypt(m *big.Int, p *big.Int, g *big.Int, dB *big.Int) (*big.Int, *big.Int) {
	k, _ := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
	k.Add(k, big.NewInt(1))

	r := new(big.Int).Exp(g, k, p)

	dBk := new(big.Int).Exp(dB, k, p)
	e := new(big.Int).Mul(m, dBk)
	e.Mod(e, p)

	return r, e
}

func decrypt(r *big.Int, e *big.Int, p *big.Int, cB *big.Int) *big.Int {
	exp := new(big.Int).Sub(p, big.NewInt(1))
	exp.Sub(exp, cB)
	rExp := new(big.Int).Exp(r, exp, p)
	m := new(big.Int).Mul(e, rExp)
	m.Mod(m, p)
	return m
}

func main() {
	p, _ := rand.Prime(rand.Reader, 128)
	g := big.NewInt(2)

	cB, dB := generateKeys(p, g)
	fmt.Printf("Секретный ключ B: %s\n", cB.String())
	fmt.Printf("Открытый ключ B: %s\n", dB.String())

	m := big.NewInt(56734542334456)

	r, e := encrypt(m, p, g, dB)
	fmt.Printf("Зашифрованное сообщение (r, e): (%s, %s)\n", r.String(), e.String())

	decryptedM := decrypt(r, e, p, cB)
	fmt.Printf("Расшифрованное сообщение: %s\n", decryptedM.String())
}
