package analyze

import (
	"encoding/json"
	"fmt"

	"github.com/nomasters/hashmap"
)

// Payload analyzes a payload and prints results to stdOut
func Payload(input []byte) error {
	p := hashmap.Payload{}
	if err := json.Unmarshal(input, &p); err != nil {
		return fmt.Errorf("invalid payload: %v\n", err)
	}

	// Outputs Payload as Indented JSON string
	fmt.Println("\nPayload\n-------\n")
	payload, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(payload))

	// Outputs Message as Indented JSON string
	fmt.Println("\nData\n-------\n")
	d, err := p.GetData()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	// Outputs Data as string
	fmt.Println("\nMessage\n----\n")
	message, err := d.MessageBytes()
	if err != nil {
		return err
	}
	fmt.Println(string(message))

	fmt.Println("\nChecker\n-------\n")

	fmt.Println("Verify Payload      : " + verifyChecker(p))
	fmt.Println("Validate TTL        : " + ttlChecker(*d))
	fmt.Println("Validate Timestamp  : " + timeStampChecker(*d))
	fmt.Println("Validate Data Size  : " + dataSizeChecker(*d))
	return nil
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
