// This is an experiment with a streamlined encoding setup.

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gogo/protobuf/proto"
	hashmap "github.com/nomasters/hashmap/x/hashmap"

	ptypes "github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/nacl/sign"
)

type Version uint8
type Method uint8

const (
	_ Version = iota
	V0_1_0
)

const (
	_ Method = iota
	NaClSign
)

const (
	MaxLen = 512
)

type Entry struct {
	Version   Version
	Method    Method
	TimeStamp time.Time
	TTL       time.Duration
	Pub       [32]byte
	Sig       [64]byte
	Len       uint16
	Data      []byte
}

func main() {
	pk, sk, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	message := []byte("ipfs/zb2rherYQC4ZJw2gZD4kYNoRcjjTy4C5HKrZVdRK2pgStJT87")

	fmt.Println("\n-- old format raw --\n")
	oldEntryRaw := []byte(`{"data":"eyJtZXNzYWdlIjoiYVhRZ2QyOXlhM009IiwidGltZXN0YW1wIjoxNTUxODkyOTk3NDQ5MDAwMDAwLCJzaWdNZXRob2QiOiJuYWNsLXNpZ24tZWQyNTUxOSIsInZlcnNpb24iOiIwLjAuMSIsInR0bCI6ODY0MDB9","sig":"aPHttyXucZslTD51boL9V0Glsfvv6ZQMaqo/JSBErKJ5Os/CZjwBN/PH5qA+MO+vuTyz+aCSr7TYMJ+TBj8cAw==","pubkey":"/7lFDizHeQMZpEVTtKPhZ/Gn7KcZoBa+bCJnh5gsRp8="}`)
	fmt.Println(oldEntryRaw)
	fmt.Println("\n-- old format len --\n")
	fmt.Println(len(oldEntryRaw))
	fmt.Println("\n-- old format hex --\n")
	fmt.Printf("%x\n", oldEntryRaw)
	fmt.Printf("\n\n\n")

	entry := Entry{
		Version:   V0_1_0,
		Method:    NaClSign,
		TimeStamp: time.Now(),
		TTL:       time.Duration(15 * time.Minute),
		Pub:       *pk,
		Len:       uint16(len(message)),
		Data:      message,
	}

	entry.Sign(sk)

	entry2, err := NewEntryFromBytes(entry.Bytes())
	if err != nil {
		panic(err)
	}

	if err := entry2.Verify(); err != nil {
		panic(err)
	}
	fmt.Println("\n-- binary blob raw --\n")
	fmt.Println(entry2.Bytes())
	fmt.Println("\n-- binary blob len --\n")
	fmt.Println(len(entry2.Bytes()))
	fmt.Println("\n-- binary blob hex --\n")
	fmt.Printf("%x\n", entry2.Bytes())
	fmt.Printf("\n\n\n")

	pbTimeNow, _ := ptypes.TimestampProto(entry.TimeStamp)

	var protoEntry hashmap.Entry
	protoEntry.Version = hashmap.Entry_V1
	protoEntry.Method = hashmap.Entry_NACLSIGN
	protoEntry.Timestamp = pbTimeNow
	protoEntry.Ttl = ptypes.DurationProto(time.Duration(15 * time.Minute))
	protoEntry.Pub = entry.Pub[:]
	protoEntry.Sig = entry.Sig[:]
	protoEntry.Len = uint32(len(message))
	protoEntry.Data = message
	pbOut, err := proto.Marshal(&protoEntry)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n-- protobuff raw -- \n")
	fmt.Println(pbOut)
	fmt.Println("\n-- protobuff raw len --\n")
	fmt.Println(len(pbOut))
	fmt.Println("\n-- protobuff hex --\n")
	fmt.Printf("%x\n", pbOut)
}

func (e Entry) Verify() error {
	c := [][]byte{
		e.Sig[:],
		e.BytesToSign(),
	}
	signedMessage := bytes.Join(c, nil)

	_, ok := sign.Open(nil, signedMessage, &e.Pub)
	if ok {
		return nil
	}
	return errors.New("verify failed")
}

func NewEntryFromBytes(raw []byte) (Entry, error) {
	var entry Entry

	if len(raw) < 116 {
		return entry, errors.New("unexpected bytes length")
	}

	version, err := versionFromByte(raw[0])
	if err != nil {
		return entry, err
	}
	method, err := methodFromByte(raw[1])
	if err != nil {
		return entry, err
	}
	timeStamp, err := timeStampFromBytes(raw[2:10])
	if err != nil {
		return entry, err
	}
	ttl, err := ttlFromBytes(raw[10:18])
	if err != nil {
		return entry, err
	}
	pubKey, err := pubKeyFromBytes(raw[18:50])
	if err != nil {
		return entry, err
	}

	sig, err := sigFromBytes(raw[50:114])
	if err != nil {
		return entry, err
	}

	len, err := lenFromBytes(raw[114:116])
	if err != nil {
		return entry, err
	}

	entry = Entry{
		Version:   version,
		Method:    method,
		TimeStamp: timeStamp,
		TTL:       ttl,
		Pub:       pubKey,
		Sig:       sig,
		Len:       len,
		Data:      raw[116:],
	}
	return entry, nil
}

func sigFromBytes(b []byte) ([64]byte, error) {
	var sig [64]byte

	if len(b) != 64 {
		return sig, errors.New("invalid byte slice length, expecting 64 bytes")
	}
	copy(sig[:], b)
	return sig, nil
}

func lenFromBytes(b []byte) (uint16, error) {
	if len(b) != 2 {
		return 0, errors.New("invalid len bytes length, expecting 2")
	}
	return binary.BigEndian.Uint16(b), nil
}

func pubKeyFromBytes(b []byte) ([32]byte, error) {
	var pk [32]byte

	if len(b) != 32 {
		return pk, errors.New("invalid byte slice length, expecting 32 bytes")
	}
	copy(pk[:], b)
	return pk, nil
}

func versionFromByte(b byte) (Version, error) {
	v := Version(uint8(b))
	if v != V0_1_0 {
		return Version(0), errors.New("unsupported version number")
	}
	return v, nil
}

func methodFromByte(b byte) (Method, error) {
	m := Method(uint8(b))
	if m != NaClSign {
		return Method(0), errors.New("unsupported method type")
	}
	return m, nil
}

func ttlFromBytes(b []byte) (time.Duration, error) {
	if len(b) != 8 {
		return time.Duration(0), errors.New("invalid ttl")
	}
	return time.Duration(int64(binary.BigEndian.Uint64(b))), nil
}

func timeStampFromBytes(b []byte) (time.Time, error) {
	if len(b) != 8 {
		return time.Unix(0, 0), errors.New("invalid timeStamp")
	}
	return time.Unix(0, int64(binary.BigEndian.Uint64(b))), nil
}

func (e *Entry) Sign(sk *[64]byte) {
	signed := sign.Sign(nil, e.BytesToSign(), sk)
	var sig [64]byte
	copy(sig[:], signed[:64])
	e.Sig = sig
}

func (e Entry) Bytes() []byte {
	c := [][]byte{
		e.Version.Bytes(),
		e.Method.Bytes(),
		e.TimeStampBytes(),
		e.TTLBytes(),
		e.Pub[:],
		e.Sig[:],
		e.LenBytes(),
		e.Data,
	}
	return bytes.Join(c, nil)
}

func (e Entry) BytesToSign() []byte {
	c := [][]byte{
		e.Version.Bytes(),
		e.Method.Bytes(),
		e.TimeStampBytes(),
		e.TTLBytes(),
		e.LenBytes(),
		e.Data,
	}
	return bytes.Join(c, nil)
}

func (e Entry) LenBytes() []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, e.Len)
	return b
}

func (e Entry) TimeStampBytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(e.TimeStamp.UnixNano()))
	return b
}

func (e Entry) TTLBytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(e.TTL.Nanoseconds()))
	return b
}

func (v Version) Bytes() []byte {
	return []byte{uint8(v)}
}

func (m Method) Bytes() []byte {
	return []byte{uint8(m)}
}
