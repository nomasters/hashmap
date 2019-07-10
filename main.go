// This is free and unencumbered software released into the public domain.

// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.

// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

// For more information, please refer to <http://unlicense.org>

package main

import "github.com/nomasters/hashmap/cmd"

func main() {
	cmd.Execute()
}

/*

This is a description of the CLI tool generally.


hashmap version

hashmap generate keys //default of ed25519,xmss_SHA2_10_256 and output current dir
hashmap generate keys --types=ed25519,xmss_SHA2_10_256
hashmap generate keys --output=/path

hashmap generate payload --message="hello, world" --ttl=5s --timestamp=now --keys=/path

$ hashmap analyze < payload
- output as TOML by default, allow out
Pubkey Hash: 	HASH
version:		v1
timestamp: 		1562017791651859000
ttl:			86400
sig bundles:
	sig1:
		alg: 	nacl_sign
		pub: 	HASH
		sig: 	HASH
	sig2:
		alg:	xmss_sha2_10_256
		pub:	HASH
		sig:	HASH

data: HASH

Valid Payload: 				TRUE
Valid Version: 				TRUE
Valid Timestamp: 			TRUE
Within Submission Window: 	FALSE
Expired TTL: 				FALSE
Valid Data Size:			TRUE
Valid Signatures: 			TRUE
Verify sig1: 				TRUE
Verify sig2: 				TRUE


hashmap --ge


- hashmap


*/
