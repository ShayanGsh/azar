package main

import (
	"github.com/Klaushayan/azar/api"
)

func main() {
	c, e := api.LoadConfigFromEnv()
	if e != nil {
		panic(e)
	}

	s := api.NewServer(c)
	s.Start()
}