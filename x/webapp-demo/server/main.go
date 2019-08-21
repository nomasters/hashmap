package main

import (
	"log"
	"net/http"

	_ "github.com/nomasters/hashmap/x/webapp-demo/statik"
	"github.com/rakyll/statik/fs"
)


// TODO:
// move assets to /assets
// / or /:hash should resolve to /index.html /assets should resolve assets


// Before buildling, run go generate.
// Then, run the main program and visit http://localhost:8080/public/hello.txt
func main() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	http.ListenAndServe(":8080", nil)
}
