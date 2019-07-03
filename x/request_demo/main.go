package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nomasters/hashmap/pkg/payload"
	"github.com/nomasters/hashmap/pkg/sig"
)

func main() {
	s := []sig.Signer{sig.GenNaclSign()}
	m := []byte("hello, world")
	p, err := payload.Generate(m, s)
	if err != nil {
		log.Fatal(err)
	}
	pb, err := payload.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	baseURL := "http://localhost:3000"
	contentType := "application/protobuf"
	resp, err := http.Post(baseURL, contentType, bytes.NewReader(pb))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	fmt.Println("successful post for:   ", p.Endpoint())

	resp2, err := http.Get(fmt.Sprintf("%v/%v", baseURL, p.Endpoint()))
	if err != nil {
		log.Fatal(err)
	}
	pbResp, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal(err)
	}
	pResp, err := payload.Unmarshal(pbResp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get is successful for: ", pResp.Endpoint())
	fmt.Println("payload response data: ", string(pResp.Data))
}
