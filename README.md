# hashmap
a light-weight cryptographically signed key value store inspired by IPNS

[![circleci][1]][2] [![Go Report Card][3]][4]

[1]: https://circleci.com/gh/nomasters/hashmap.svg?style=shield&circle-token=46ac657a268fef44dc132ef2241291c51811edd2
[2]: https://circleci.com/gh/nomasters/hashmap
[3]: https://goreportcard.com/badge/github.com/nomasters/hashmap
[4]: https://goreportcard.com/report/github.com/nomasters/hashmap



`hashmap` is a light-weight cryptographically signed key-value store inspired by IPNS. The purpose of this tool is to allow any user to generate a unique `ed25519` private key and use the corresponding hash of the public key as the verified REST enppoint for a key-value store.


## Notes

This is a very early and incomplete prototype of the `hashmap` server and basic tools. Currently the DB is in-memory only and is deleted between runs. This code needs test coverage and possibly some rethinking on some of the structures, but this is working as an MVP.


## Basic instructions


While in development, the easiest way to run the `hashmap` CLI tool is to run

```
./scripts/build.sh
```

you can run the hashmap server from the cli with:

```
hashmap run
```

The server runs on `localhost:3000`

You can test sending a properly formatted json payload by using curl (read below to find out how to use `hashmap generate` to generate a key and a payload)

```
curl -X POST http://localhost:3000 -d @payload.json
```
This will respond with a multihash base58 encoded pubkey hash

you can use this hash to query hashMap like this:

`curl http://localhost:3000/2DrjgbD2fh4CL6HX5qYqKf7ULr3hwJXQgYn9sTCQSLpHAqj5n2`

## Other instructions

Also included in `hashmap` command are tools to make it easier to generate and `ed25519` private key as well as generate a properly formatted payload for submitting to the hashmap server.


## Generating an `ed25519` private key

You can generate a key to `stdout` encoded to base64 with:

```
hashmap generate key
```

## Saving the private key

If you'd like to save that key to a file for future use, its as easy as:

```
hashmap generate key > priv.key
```

## Generating a Payload

If you'd like to generate a payload with defaults use:

```
hashmap generate payload < priv.key
```

you can also change the default payload data for `data`, `timestamp`, and `ttl`. To look at the CLI options you can use the `help` flag

```
hashmap generate payload --help
```

an example of modifying the inputs is as follows:

```
hashmap generate payload --message="{\"hello\":\"world\"}" --timestamp=1534121771 --ttl=5 < priv.key
```

you can save this ouput to a file as follows:

```
hashmap generate payload < priv.key > payload.json
```

## Analyzing a Payload

To analyze a payload, you can run the analyzer as follows:

```
hashmap analyze < payload.json

Payload
-------

{
  "data": "eyJkYXRhIjoiZXlKb1pXeHNieUk2SW5kdmNteGtJbjA9IiwidGltZXN0YW1wIjoxNTM0MTYyMzgzLCJ0dGwiOjg2NDAwLCJzaWdNZXRob2QiOiJuYWNsLXNpZ24tZWQyNTUxOSIsInZlcnNpb24iOiIwLjAuMSJ9",
  "sig": "h7clARjoYeh3Mmg7EOsKb0QVpvhKUYymFeZ7tFIyGqdNd5mt/QMmvtO/fWy9/nYbcDXQ0+37VFmhpBjMEFXlAQ==",
  "pubkey": "z0CRLsemGDadzmzA9/3R3e4JkEtVZLOD+gAU7EtychQ="
}

Data
----

{
  "message": "eyJoZWxsbyI6IndvcmxkIn0=",
  "timestamp": 1534162383,
  "ttl": 86400,
  "sigMethod": "nacl-sign-ed25519",
  "version": "0.0.1"
}

Message
-------

{"hello":"world"}

Checker
-------

Verify Payload         : PASS
Validate TTL           : PASS
Validate Timestamp     : PASS
Validate Message Size  : PASS
```
