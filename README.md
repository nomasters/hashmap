# hashmap
a light-weight cryptographically signed key value store inspired by IPNS

This is a very early and incomplete prototype of the hashMap server. Currently the DB is in-memory only and is deleted between runs. This code needs test coverage and possibly some rethinking on some of the structures, but this is working as an MVP.


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

You can test sending a properly formatted json payload by using curl (read below to find out how to use `hashmap-helper` to generate a payload)

```
curl -X POST http://localhost:3000 -d @payload.json
```
This will respond with a multihash base58 encoded pubkey hash

you can use this hash to query hashMap like this:

`curl http://localhost:3000/2DrjgbD2fh4CL6HX5qYqKf7ULr3hwJXQgYn9sTCQSLpHAqj5n2`

# Helper instructions

also included in this repo is a tool called `hashmap-helper` this makes it easier to generate and `ed25519` private key as well as generate a properly formatted payload for submitting to the hashmap server. to install

```
./scripts/build-helper.sh
```

You can generate a key to `stdOut` like this:

```
hashmap-helper gen-key
```

if you'd like to save that key to a file for future use, its as easy as:

```
hashmap-helper gen-key > priv.key
```

if you'd like to generate a payload with defaults use:

```
hashmap-helper gen-payload < priv.key
```

you can also change the default payload data for `data`, `timestamp`, and `ttl`. To look at the CLI options you can use the `help` flag


```
hashmap-helper gen-payload --help
```

an example of modifying the inputs is as follows:

```
hashmap-helper gen-payload --data="{\"hello\":\"world\"}" --timestamp=1534121771 --ttl=5 < priv.key
```
