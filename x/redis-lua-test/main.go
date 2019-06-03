package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

var key = "a"
var value = `{
	"payload": "Ob3mbvpnROa72DTnBpz2eZn5jFfL610zUISsqyyrtF-RY2W3HqpGc4eSYBtkDKn6yhZBqHEyXwAcB1Wketu6Wg==",
	"timestamp": 1257894000000000000 
}`

var newValue = `{
	"payload": "Ob3mbvpnROa72DTnBpz2eZn5jFfL610zUISsqyyrtF-RY2W3HqpGc4eSYBtkDKn6yhZBqHEyXwAcB1Wketu6Wg==",
	"timestamp": 1257894000000000001
}`

var luaScript = `
if redis.call("EXISTS", KEYS[1]) == 1 then
	local payload = redis.call("GET", KEYS[1])
	local timestamp = cjson.decode(payload)["timestamp"]
	if timestamp > tonumber(ARGV[2]) then
		return "INVALID TIMESTAMP"
	end
end
return redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
`

func main() {
	c, _ := redis.Dial("tcp", ":6379")
	defer c.Close()

	// load up the keys
	c.Do("SET", key, value)
	conditionalSet := redis.NewScript(1, luaScript)
	// set key with value if timestamp > current timestamp, and set a 10 second TTL
	reply, err := conditionalSet.Do(c, key, newValue, 1257893000000000001, 10)
	if err != nil {
		log.Fatal(err)
	}
	switch reply.(type) {
	case string:
		fmt.Println(reply.(string))
	case []uint8:
		fmt.Println(string(reply.([]uint8)))
	default:
		fmt.Println(reply)
	}
}
