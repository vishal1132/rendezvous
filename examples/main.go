package main

import (
	"crypto/md5"
	"log"

	"github.com/vishal1132/rendezvous"
)

func main() {
	r := rendezvous.New(md5.New(), "abcd", "efgh", "ijkl", "k", "m", "n", "o", "p")
	log.Println(r.GetNTop(5, []byte("abcd")))
}
