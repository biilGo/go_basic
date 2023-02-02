package main

import (
	"encoding/gob"
	"log"
	"os"
)

type Address struct {
	Type     string
	City     string
	Conuntry string
}

type VCard struct {
	FirstName string
	LastName  string
	Address   []*Address
	Remark    string
}

var content string

func main() {
	pa := &Address{"private", "Aartselaar", "Belgium"}

	wa := &Address{"work", "Boom", "Belgium"}

	vc := VCard{"Jan", "Kersschot", []*Address{pa, wa}, "none"}

	file, _ := os.OpenFile("D:/git_biilGo/go_basic/20221231_RW_data/format_JSON_Data/Gob_transport_data/gob2/vcard.gob", os.O_CREATE|os.O_WRONLY, 0666)

	defer file.Close()

	enc := gob.NewEncoder(file)

	err := enc.Encode(vc)

	if err != nil {
		log.Fatal("Error in encoding gob")
	}
}
