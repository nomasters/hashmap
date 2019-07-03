package main

import (
	"encoding/binary"
	"fmt"
	"log"

	xmss "github.com/danielhavir/go-xmss"
)

func main() {
	prm := xmss.SHA2_10_256
	prv, pub := xmss.GenerateXMSSKeypair(prm)
	fmt.Printf("%x\n", *prv)
	fmt.Printf("%x\n", *pub)
	fmt.Println(len(*prv))
	fmt.Println(len(*pub))
	fmt.Println("")

	fmt.Printf("%x%x\n", (*prv)[100:], (*prv)[68:100])
	var p [64]byte
	copy(p[:32], (*prv)[100:])
	copy(p[32:], (*prv)[68:100])
	fmt.Printf("%x\n", p)
	// msg := []byte("hello, world")

	// for i := 1022; i < 1023; i++ {
	// 	log.Println("verifying counter:", i)
	// 	bumpCounter(uint64(i), prv)
	// 	signAndVerify(prm, prv, pub, msg)
	// }
}

func signAndVerify(prm *xmss.Params, prv *xmss.PrivateXMSS, pub *xmss.PublicXMSS, msg []byte) {
	sig := *prv.Sign(prm, msg)
	m := make([]byte, prm.SignBytes()+len(msg))

	fmt.Println(sig)
	if !xmss.Verify(prm, m, sig, *pub) {
		log.Fatal("failed to verify")
	}
	log.Println(string(m)[prm.SignBytes():])
}

func bumpCounter(i uint64, prv *xmss.PrivateXMSS) {
	pb := []byte(*prv)
	copy(pb[:4], uint64ToBytes(i)[4:])
	copy(*prv, pb)
}

func uint64ToBytes(t uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, t)
	return b
}
