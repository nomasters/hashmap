//go:generate protoc -I=pb --go_out=pb/ pb/payload.proto

package payload

// TODO

// marshal and unmarshal a payload to protobuff

// a payload should take protobuff bytes, unmarshal, and then return
// a payload interface, which support specific actions.
