package test_utils

import (
	"path/filepath"
	"strconv"

	"github.com/ShayanGsh/azar/api"
)

func GenerateConfig(port string) *api.Config {
	var config api.Config = *api.NewConfig()
	mappedPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	config.Database.Port = mappedPort
	config.Database.Name = "azar_test"
	config.MigrationPath, err = filepath.Abs("migrations")
	if err != nil {
		panic(err)
	}

	return &config
}