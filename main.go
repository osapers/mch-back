package main

import (
	"log"

	"github.com/osapers/mch-back/cmd"
)

func main() {
	if err := cmd.Launch(); err != nil {
		log.Fatal(err)
	}
}
