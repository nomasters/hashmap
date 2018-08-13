package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nomasters/hashmap"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/nacl/sign"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no args given. must use gen-key|gen-payload")
	}
	arg := os.Args[1]
	switch arg {
	case "gen-key":
		genKey()
	case "gen-payload":
		genPayload()
	default:
		log.Fatal("invalid arg")
	}
}

func genPayload() {

	d := flag.StringP("data", "d", `{"content":"hello, world. This is data stored in HashMap."}`, "data to be stored in the message")
	ttl := flag.Int64P("ttl", "t", hashmap.MessageTTLDefault, "ttl in seconds for payload")
	ts := flag.Int64P("timestamp", "s", time.Now().Unix(), "timestamp for message in unix-time")

	flag.Parse()

	// Create the Message for Payload, and marshal the JSON
	data := base64.StdEncoding.EncodeToString([]byte(*d))
	m := hashmap.Message{
		Data:      data,
		Timestamp: *ts,
		TTL:       *ttl,
		SigMethod: hashmap.DefaultSigMethod,
		Version:   hashmap.Version,
	}
	message, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	// sign this Mashalled message with the PrivKey
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	pk, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Fatal(err)
	}
	var privateKey [64]byte
	var publicKey [32]byte
	copy(privateKey[:], pk)
	copy(publicKey[:], pk[32:])

	s := sign.Sign(nil, message, &privateKey)[:64]
	sig := base64.StdEncoding.EncodeToString(s)

	p := hashmap.Payload{
		Message:   base64.StdEncoding.EncodeToString(message),
		Signature: sig,
		PublicKey: base64.StdEncoding.EncodeToString(publicKey[:]),
	}

	payload, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	if _, valid := sign.Open(nil, append(s, message...), &publicKey); !valid {
		log.Fatal("sig failed")
	}

	fmt.Printf("%s\n", payload)

}

// quick and dirty priv key generation. disregarding pubkey and error
// pubkey can be derived from the privKey, so we don't need to keep it
func genKey() {
	_, privKey, _ := sign.GenerateKey(rand.Reader)
	p := base64.StdEncoding.EncodeToString(privKey[:])
	fmt.Printf("%s", p)
}
