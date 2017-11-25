package main

import (
	"acerwei/gmailbox/central"
	"flag"
	"fmt"
)

var (
	opt          = flag.String("o", "retrieve", "retrieve | decode ")
	fileToDecode = flag.String("file", "", "file to decode")
	path         = flag.String("path", "./data", "root path to store emails")
	encypt       = flag.String("encrypt", "blowfish", "encryption algorithm (simple | blowfish)")
	encKey       = flag.String("enckey", "", "encryption key")
	startDate    = flag.String("startDate", "2015/11/11", "start date (inclusive), format: YYYY-MM-DD")
	endDate      = flag.String("endDate", "2016/06/01", "end date (exclusive), format: YYYY-MM-DD")
)

func main() {
	flag.Parse()
	err := central.Initialize(*path, *encypt, *encKey)
	if err != nil {
		fmt.Printf("[ERROR] initialize gmailbox error: %v", err)
	}
	if *opt == "retrieve" {
		central.Retrieve(*startDate, *endDate)
	} else if *opt == "decode" {
		central.DecodeMail(*fileToDecode)
	} else {
		fmt.Printf("[ERROR] unsuportted operation")
	}
}
