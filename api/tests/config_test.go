package tests

import (
	"os"
	"testing"

	"github.com/Klaushayan/azar/api"
	"github.com/stretchr/testify/assert"
)

func setTestEnvs() {
	os.Setenv("AZAR_PORT", "5000")
	os.Setenv("AZAR_HOST", "localhost")
	os.Setenv("AZAR_JWT_SECRET", "thetestingsecret")
	os.Setenv("AZAR_DB_HOST", "localhost")
	os.Setenv("AZAR_DB_PORT", "5432")
	os.Setenv("AZAR_DB_USER", "postgres")
	os.Setenv("AZAR_DB_PASS", "")
	os.Setenv("AZAR_DB_NAME", "azar_test")
}

var config *api.Config

func init() {
	setTestEnvs()
}

func TestLoadConfig(t *testing.T) {
	var err error
	config, err = api.LoadConfig("config_example.json")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 5000, config.Port)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "thetestingsecret", config.JWTSecret)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "postgres", config.Database.Username)
	assert.Equal(t, "", config.Database.Password)
	assert.Equal(t, "azar_test", config.Database.Name)
}

func TestLoadConfigFromEnv(t *testing.T) {
	var err error
	config, err = api.LoadConfigFromEnv()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 5000, config.Port)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "thetestingsecret", config.JWTSecret)
	assert.Equal(t, "localhost", config.DatabaseHost)
	assert.Equal(t, 5432, config.DatabasePort)
	assert.Equal(t, "postgres", config.DatabaseUsername)
	assert.Equal(t, "", config.DatabasePassword)
	assert.Equal(t, "azar_test", config.DatabaseName)
}

func TestLoadConfigFromEnvWithMissingConfig(t *testing.T) {
	_, err := api.LoadConfig("missing_config.json")
	assert.Error(t, err)
}

func TestConfigConnString(t *testing.T) {
	connString := config.ToConnString()
	assert.Equal(t, "host=localhost port=5432 user=postgres password= dbname=azar_test sslmode=disable", connString)
}

func TestDatabaseConnString(t *testing.T) {
	connString := config.Database.ToConnString()
	assert.Equal(t, "host=localhost port=5432 user=postgres password= dbname=azar_test sslmode=disable", connString)
}

func TestConfigAddress(t *testing.T) {
	address := config.Address()
	assert.Equal(t, "localhost:5000", address)
}