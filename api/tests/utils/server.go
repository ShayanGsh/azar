package test_utils

import (
	"strconv"

	"github.com/Klaushayan/azar/api"
)

func GetServer(port string, config api.Config) *api.Server {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	config.Port = intPort
	s := api.NewServer(&config)
	return s
}