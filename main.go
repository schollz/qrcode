package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"rsc.io/qr"
)

func main() {
	b, err := ioutil.ReadFile("test.gpg")
	q, err := qr.Encode(string(b), qr.L)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q, err)
	q.PNG()
	err = ioutil.WriteFile("test.png", q.PNG(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
