package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
)

func main() {
	var signingTypes []string
	flag.StringSliceVarP(&signingTypes, "signing-type","t", []string{}, "declares the signing type for keygen")
	flag.Parse()
	fmt.Println(signingTypes)
}