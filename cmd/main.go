package main

import (
	"github.com/ShayanGsh/azar/api"
	"log"
)

func main() {
	c, e := api.LoadConfigFromEnv()
	if e != nil {
		panic(e)
	}

	s := api.NewAPIServer(c)
	s.MigrationCheck()
	s.Start()
	log.Println("Started")
}