package test_utils

import (
	"os"
	"strconv"
	"strings"

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
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config.MigrationPath = strings.Join(strings.Split(currentPath, "/")[:len(strings.Split(currentPath, "/"))-1], "/") + "/azar-db/migrations"

	return &config
}