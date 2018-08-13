package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nomasters/hashmap"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/nacl/sign"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no args given. must use analyze|gen-key|gen-payload")
	}
	arg := os.Args[1]
	switch arg {
	case "analyze":
		analyze()
	case "gen-key":
		genKey()
	case "gen-payload":
		genPayload()
	default:
		log.Fatal("invalid arg")
	}
}

func analyze() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	p := hashmap.Payload{}
	if err := json.Unmarshal(input, &p); err != nil {
		log.Fatalf("invalid payload: %v\n", err)
	}

	// Outputs Payload as Indented JSON string
	fmt.Println("\nPayload\n-------\n")
	payload, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(payload))

	// Outputs Message as Indented JSON string
	fmt.Println("\nMessage\n-------\n")
	m, err := p.GetMessage()
	if err != nil {
		log.Fatal(err)
	}

	message, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(message))

	// Outputs Data as string
	fmt.Println("\nData\n----\n")
	data, err := m.DataBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	fmt.Println("\nChecker\n-------\n")

	fmt.Println("Verify Payload      : " + verifyChecker(p))
	fmt.Println("Validate TTL        : " + ttlChecker(*m))
	fmt.Println("Validate Timestamp  : " + timeStampChecker(*m))
	fmt.Println("Validate Data Size  : " + dataSizeChecker(*m))
}

func verifyChecker(p hashmap.Payload) string {
	status := "PASS"
	if err := p.Verify(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func ttlChecker(m hashmap.Message) string {
	status := "PASS"
	if err := m.ValidateTTL(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func timeStampChecker(m hashmap.Message) string {
	status := "PASS"
	if err := m.ValidateTimeStamp(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func dataSizeChecker(m hashmap.Message) string {
	status := "PASS"
	if err := m.ValidateDataSize(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
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
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	pk, err := base64.StdEncoding.DecodeString(string(text))
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
