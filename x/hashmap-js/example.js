const hashmap = require('./index.js');
const now = require('nano-time');

var key = "5yGLD+NlUB9gQ76XU/GbEF2BuLxRyLBIuw5HSt+b9gOgHe6nBhkt4bumuTcxe6fMOe19GaKNvcUgnfWqJokDNA=="

var opts = {
    endpoint: "https://prototype.hashmap.sh",
    multiHash: "2DrjgbD6zUx2svjd4NcXfsTwykspqEQmcC2WC7xeBUyPcBofuo",
}

// get payload example
p1 = new hashmap.Payload(opts)
p1.get()
.then(payload => console.log(p1.getMessage()))
.catch(err => console.log(err))

// post payload example
var opts = { endpoint: "https://prototype.hashmap.sh" }
p2 = new hashmap.Payload(opts)
p2.generate(key, "hello, world it is: " + now())
p2.post()
.then(resp => console.log(resp))
.catch(err => console.log(err))
