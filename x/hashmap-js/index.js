const nacl = require('tweetnacl');
const rp = require('request-promise');
const now = require('nano-time');

const maxMessageBytes  = 512
const defaultSigMethod = 'nacl-sign-ed25519'
const dataTTLDefault   = 86400  // 1 day in seconds
const dataTTLMax       = 604800 // 1 week in seconds
const version          = '0.0.1'

class Payload {
    constructor(opts={}) {
        if (opts.endpoint) this.endpoint = opts.endpoint
        if (opts.multiHash) this.multiHash = opts.multiHash
    }
    // generate takes a base64 encoded key, a message string, and opts object
    // and creates a properly formatted and signed payload
    // it returns a JSON encoded string and sets the class
    // internal state for use with other methods
    generate(key, message=' ', opts={}) {
        const ttl = dataTTLDefault
        if (opts && opts.ttl) {
            ttl = opts.ttl
        }
        if (ttl > dataTTLMax) {
            throw "invalide ttl, exceeds max"
        }
        const data = {
            message: Buffer.from(message, 'ascii').toString('base64'),
            timestamp: Number(now()),
            sigMethod: defaultSigMethod,
            version: version,
            ttl: ttl,
        }
        // dataBytes takes a byte buffer of the data object that has been stringified
        const dataBytes = Buffer.from(JSON.stringify(data), 'ascii')
        const privKey = Buffer.from(key, 'base64');
        const pubKey = privKey.slice(32,64)
        const signedMessage = nacl.sign(dataBytes, privKey)
        const sig = (new Buffer.from(signedMessage.slice(0,64))).toString('base64')

        const p = {
            data: dataBytes.toString('base64'),
            pubkey: pubKey.toString('base64'),
            sig: sig.toString('base64'),
        }
        this.validate(p)
        return JSON.stringify(p)
    }
    get(multiHash, endpoint) {
        if (!multiHash && !this.multiHash) { throw "missing multiHash" }
        if (multiHash) this.multiHash = multiHash
        if (!endpoint && !this.endpoint) { throw "missing endpoint" }
        if (endpoint) this.endpoint = endpoint
        var opts = {
            uri: this.endpoint + '/' + this.multiHash,
            json: true
        }
        return rp(opts)
        .then(resp => {
            this.validate(resp)
            return resp
        })
        .catch(err => { throw err })
    }
    post(endpoint) {
        if (!endpoint && !this.endpoint) { throw "missing endpoint" }
        if (!this.raw) { throw "missing payload" }
        if (endpoint) this.endpoint = endpoint
        var opts = {
            uri: this.endpoint,
            method: 'POST',
            body: this.raw,
            json: true,
        }
        return rp(opts)
    }
    import(raw) {
        this.validate(JSON.parse(raw))
    }
    validate(p, opts={}) {
        this.raw     = p
        this.data    = Buffer.from(p.data, 'base64');
        this.sig     = Buffer.from(p.sig, 'base64');
        this.pubkey  = Buffer.from(p.pubkey, 'base64');
        this.validateMessage()
        this.validateSig()
    }
    validateMessage() {
        if (this.getMessageBytes().length > maxMessageBytes) {
            throw "message length exceeds max threshold"
        }
    }
    validateSig() {
        var signedData = Buffer.concat([this.sig, this.data]);
        if (!nacl.sign.open(signedData, this.pubkey)) {
            throw "signature validation failed"
        }
    }
    getData() {
        return JSON.parse(this.data.toString('ascii'))
    }
    getMessageBytes() {
        var m = this.getData().message
        return Buffer.from(m, 'base64');
    }
    getMessage() {
        var buf = this.getMessageBytes()
        return buf.toString('ascii');
    }
}

module.exports = {
    Payload
};