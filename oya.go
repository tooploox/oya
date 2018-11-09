package main

import (
	"log"

	getter "github.com/hashicorp/go-getter"
)

func main() {

	// repo, err := vcs.NewRepo("https://github.com/bilus/akasha", "/tmp/akasha")
	// if err != nil {
	// 	panic(err)
	// }
	// err = repo.Get()
	// if err != nil {
	// 	panic(err)
	// }

	err := getter.GetAny("/tmp/baz", "github.com/bilus/akasha//bin?ref=v0.4.0")
	log.Printf("Result: %v", err)

}
