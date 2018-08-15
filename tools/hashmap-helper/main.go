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
	fmt.Println("\nData\n-------\n")
	d, err := p.GetData()
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	// Outputs Data as string
	fmt.Println("\nMessage\n----\n")
	message, err := d.MessageBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(message))

	fmt.Println("\nChecker\n-------\n")

	fmt.Println("Verify Payload      : " + verifyChecker(p))
	fmt.Println("Validate TTL        : " + ttlChecker(*d))
	fmt.Println("Validate Timestamp  : " + timeStampChecker(*d))
	fmt.Println("Validate Data Size  : " + dataSizeChecker(*d))
}

func verifyChecker(p hashmap.Payload) string {
	status := "PASS"
	if err := p.Verify(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func ttlChecker(d hashmap.Data) string {
	status := "PASS"
	if err := d.ValidateTTL(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func timeStampChecker(d hashmap.Data) string {
	status := "PASS"
	if err := d.ValidateTimeStamp(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func dataSizeChecker(d hashmap.Data) string {
	status := "PASS"
	if err := d.ValidateMessageSize(); err != nil {
		status = "FAIL - " + err.Error()
	}
	return status
}

func genPayload() {

	m := flag.StringP("message", "m", `{"content":"hello, world. This is data stored in HashMap."}`, "message to be stored in data of payload")
	ttl := flag.Int64P("ttl", "t", hashmap.DataTTLDefault, "ttl in seconds for payload")
	ts := flag.Int64P("timestamp", "s", time.Now().Unix(), "timestamp for message in unix-time")

	flag.Parse()

	// Create the Message for Payload, and marshal the JSON
	message := base64.StdEncoding.EncodeToString([]byte(*m))
	d := hashmap.Data{
		Message:   message,
		Timestamp: *ts,
		TTL:       *ttl,
		SigMethod: hashmap.DefaultSigMethod,
		Version:   hashmap.Version,
	}
	data, err := json.Marshal(d)
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

	s := sign.Sign(nil, data, &privateKey)[:64]
	sig := base64.StdEncoding.EncodeToString(s)

	p := hashmap.Payload{
		Data:      base64.StdEncoding.EncodeToString(data),
		Signature: sig,
		PublicKey: base64.StdEncoding.EncodeToString(publicKey[:]),
	}

	payload, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	if _, valid := sign.Open(nil, append(s, data...), &publicKey); !valid {
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
