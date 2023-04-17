package main

import (
	"github.com/Klaushayan/azar/api"
	"log"
)

func main() {
	c, e := api.LoadConfigFromEnv()
	if e != nil {
		panic(e)
	}

	s := api.NewServer(c)
	s.MigrationCheck()
	s.Start()
	log.Println("Started")
}