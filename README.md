![hashmap logo](images/hashmap-logo-neon.png)

# hashmap

hashmap server is a light-weight cryptographically signed public key value store inspired by JWT and IPNS.

[![CircleCI][1]][2] [![Go Report Card][3]][4] [![codecov][5]][6]

[1]: https://circleci.com/gh/nomasters/hashmap.svg?style=svg
[2]: https://circleci.com/gh/nomasters/hashmap
[3]: https://goreportcard.com/badge/github.com/nomasters/hashmap
[4]: https://goreportcard.com/report/github.com/nomasters/hashmap
[5]: https://codecov.io/gh/nomasters/hashmap/branch/master/graph/badge.svg
[6]: https://codecov.io/gh/nomasters/hashmap

## Summary

NOTE: this repo is in the midst of a large refactor for the next alpha version, the stable initial version can be found in (v0.0.4 here)[https://github.com/nomasters/hashmap/tree/v0.0.4]

The goal of `hashmap` is to serve as a key-value store for cryptographically signed payloads namespaced by the hash of the public key used to verify the signature. This allows for:

- a unique pubkey hash name spacing (one `ed25519` key pair per key-value entry)
- only the private key holder can update its corresponding public key hash endpoint
- the validity of the endpoint and the response authenticity are verifiable by all parties

The design goals were to allow a submitter to randomly generate `ed25519` keys to use for submission. Hashmap doesn't keep any user data or auth systems other than ensuring that a submitted key-value pair are cryptographically authentic and properly formatted.

A tool like `hashmap` is useful as a light-weight and mobile device friendly mutable storage endpoint. This tool is heavily inspired by IPNS, the mutable store used by the IPFS project to point to specific IPFS hashes.

One benefit of `hashmap` being decoupled from IPFS specifically is that its `message` store supports client-side encryption and therefore, a submitter who obfuscates the source IP through an anomemity network like TOR and encrypts its `message` client-side can can publicly store mutable data without the `hashmap` server having knowledge or origin of the submission nor contents of the message.

features:

- endpoints are hashed using the multi-hash format and default to `blake2b-256` hashes
- signed data uses nacl sign which leverages `ed25519`
- signed payload submission and acceptance are strictly enforced based on signature validity, message size, and date-stamp accuracy
- values stored in hashmap have a max time-to-live of 1 week and default to 24 hours.

The structure of a properly formatted payload submission is a follows:

```bash
{
  "data":   "BASE_64_ENCODED_STRING",
  "sig":    "BASE_64_ENCODED_STRING",
  "pubkey": "BASE_64_ENCODED_STRING"
}
```

The integrity of the contents of `data` are verifiable with the `sig` and the `pubkey`. But to know which signature method to use will require decoding the contents of data.

The contents of `data` decoded is:

```bash
{
  "message":   "BASE_64_ENCODED_STRING",
  "timestamp": 1535211171864296000,
  "ttl":       86400,
  "sigMethod": "nacl-sign-ed25519",
  "version":   "0.0.1"
}
```

- `message` is a `BASE_64_ENCODED_STRING` provided by the submitter. This may contain anything as long as the message bytes is less than 512 bytes.
- `timestamp` is unix-time in nanoseconds. Hashmap Server allows MaxSubmitDrift of 15 seconds. This prevents old payloads from overwriting newer payloads.
- `ttl` - is the time to live in seconds for the payload. If no TTL is set (or it is set to 0) hashmap defaults to 24 hours (86400 seconds ). A Maximum TTL of 1 week (604800 seconds) is permitted. Any TTL greater than 604800 will be rejected.
- `sigMethod` outlines the method used to verify the signature. Currently only `nacl-sign-ed25519` is supported
- `version` is used for handling potentially breaking changes in the future, but it isn't currently analyzed for acceptance.

For a hashmap payload to be accepted:

- Signature must be valid
- Timestamp must be within the time drift threshold
- TTL must be valid
- Message size must be valid

If a payload is accepted, the server will respond with something similar to:

```bash
{
  "endpoint": "2Drjgb7y6LmSaGZhw5pJrhpB4MrgVcajMXwnb8yWUavavaBSHo"
}
```

This endpoint is the `blake2b-256` multi-hash of the submitted public key in the payload. This is the endpoint that can be used to retrieve the submitted payload.

Hashmap responds in the exact payload as that which was submitted:

```bash
{
  "data":   "BASE_64_ENCODED_STRING",
  "sig":    "BASE_64_ENCODED_STRING",
  "pubkey": "BASE_64_ENCODED_STRING"
}
```

This means that any requestor can independently verify that:

- the payload data is valid based on the signature
- the TTL has not expired
- the pubkey multihash matches the pubkey contained in the payload

## Notes

This is a very early and incomplete prototype of the `hashmap` server and basic tools. Currently the backing store defaults to a simple in-memory store but redis is also supported. This code needs test coverage and possibly some rethinking on some of the structures, but this is working as an MVP.

## Basic instructions

While in development, the easiest way to run the `hashmap` CLI tool is to run

```bash
./scripts/build.sh
```

you can run the hashmap server locally without TLS certificate files from the cli with:

```bash
hashmap run --TLS=false
```

The server runs on `localhost:3000`

You can test sending a properly formatted json payload by using curl (read below to find out how to use `hashmap generate` to generate a key and a payload)

```bash
curl -X POST http://localhost:3000 -d @payload.json
```

This will respond with a multihash base58 encoded pubkey hash

You can use this hash to query hashmap like this:

`curl http://localhost:3000/2DrjgbD2fh4CL6HX5qYqKf7ULr3hwJXQgYn9sTCQSLpHAqj5n2`

## Other instructions

Also included in `hashmap` command are tools to make it easier to generate and `ed25519` private key as well as generate a properly formatted payload for submitting to the hashmap server.

## Generating an `ed25519` private key

You can generate a key to `stdout` encoded to base64 with:

```bash
hashmap generate key
```

## Saving the private key

If you'd like to save that key to a file for future use, its as easy as:

```bash
hashmap generate key > priv.key
```

## Generating a Payload

If you'd like to generate a payload with defaults use:

```bash
hashmap generate payload < priv.key
```

you can also change the default payload data for `data`, `timestamp`, and `ttl`. To look at the CLI options you can use the `help` flag

```bash
hashmap generate payload --help
```

an example of modifying the inputs is as follows:

```bash
hashmap generate payload --message="{\"hello\":\"world\"}" --timestamp=1535211171864296000 --ttl=5 < priv.key
```

you can save this ouput to a file as follows:

```bash
hashmap generate payload < priv.key > payload.json
```

## Analyzing a Payload

To analyze a payload, you can run the analyzer as follows:

```bash
hashmap analyze < payload.json

Payload
-------

{
  "data": "eyJtZXNzYWdlIjoiZXlKamNubHdkRzhpT25zaWJXVjBhRzlrSWpvaWMyVmpjbVYwWW05NElpd2lZMmgxYm1zaU9qRTJNREF3ZlN3aWFHRnphQ0k2SW1VeFpqRTFNelJrTWpaa01pSXNJbkJoZVd4dllXUWlPaUoxYnpoYVVXVk5WV3hwT0dJclNXOUtRMUppTDI1UGQwMVhVR2RwT0dWRFpFVkNUWGRrYW1wNlIyWTNVVTFqWkdnMlRXVlFVbW9yV25CM1VGZHpOVmhxZFdOcE1tRXpSVFUzUTFWUGRIUjNTRUZIZHowaWZRPT0iLCJ0aW1lc3RhbXAiOjE1MzUyMTExNzE4NjQyOTYwMDAsInR0bCI6MzAwLCJzaWdNZXRob2QiOiJuYWNsLXNpZ24tZWQyNTUxOSIsInZlcnNpb24iOiIwLjAuMSJ9",
  "sig": "s5tHkDsMfUQbsJGJf80lKdkI4kTEf9BJ3a1PCqSnYuH0aZ8kjeVG2SoGeeRsaXMNV4nDYlBddJkeQkGWort3Cg==",
  "pubkey": "9DHEtybGV0LxFFV3WBFMZzrvyme/3o1pP9Q6SD5t3NA="
}

Data
----

{
  "message": "eyJjcnlwdG8iOnsibWV0aG9kIjoic2VjcmV0Ym94IiwiY2h1bmsiOjE2MDAwfSwiaGFzaCI6ImUxZjE1MzRkMjZkMiIsInBheWxvYWQiOiJ1bzhaUWVNVWxpOGIrSW9KQ1JiL25Pd01XUGdpOGVDZEVCTXdkamp6R2Y3UU1jZGg2TWVQUmorWnB3UFdzNVhqdWNpMmEzRTU3Q1VPdHR3SEFHdz0ifQ==",
  "timestamp": 1535211171864296000,
  "ttl": 300,
  "sigMethod": "nacl-sign-ed25519",
  "version": "0.0.1"
}

Message
-------

{"crypto":{"method":"secretbox","chunk":16000},"hash":"e1f1534d26d2","payload":"uo8ZQeMUli8b+IoJCRb/nOwMWPgi8eCdEBMwdjjzGf7QMcdh6MePRj+ZpwPWs5Xjuci2a3E57CUOttwHAGw="}

Checker
-------

Verify Payload         : PASS
Validate TTL           : PASS
Validate Timestamp     : PASS
Validate Message Size  : PASS
```
