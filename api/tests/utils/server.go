package test_utils

import (
	"strconv"

	"github.com/Klaushayan/azar/api"
)

func GenerateConfig(port string) *api.Config {
	var config api.Config = *api.NewConfig()
	mappedPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	config.Database.Port = mappedPort
	config.Database.Name = "azar_test"

	return &config
}