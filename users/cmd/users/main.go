package main

import (
	"fmt"
	"log"
	"users/cmd/users/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("can't init config")
	}
	fmt.Println(cfg)
}
