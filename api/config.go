package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// If the config is loaded from a json file, DatabaseHost, DatabasePort, DatabaseUsername, and DatabasePassword will be ignored.
type Config struct {
	Port int `json:"port" envconfig:"PORT" default:"5000"`
	Host string `json:"host" envconfig:"HOST" default:"localhost"`
	JWTSecret string `json:"secret" envconfig:"JWT_SECRET" default:"secret"`
	JWTExpiration int `json:"expiration" default:"3600"` // in seconds
	Database Database `json:"database"`

	// ENV VARIABLES ONLY
	DatabaseHost string `json:"-" envconfig:"DB_HOST" default:"localhost"`
	DatabasePort int `json:"-" envconfig:"DB_PORT" default:"5432"`
	DatabaseUsername string `json:"-" envconfig:"DB_USER" default:"postgres"`
	DatabasePassword string `json:"-" envconfig:"DB_PASS" default:""`
	DatabaseName string `json:"-" envconfig:"DB_NAME" default:"azar"`
}

type Database struct {
		Host string `json:"host"`
		Port int `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name string `json:"name"`
}

func (db *Database) ToConnString() (string) {
    var buf bytes.Buffer
	fmt.Fprintf(&buf, "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", db.Host, db.Port, db.Username, db.Password, db.Name, "disable")
    return buf.String()
}

func (c *Config) ToConnString() (string) {
	return c.Database.ToConnString()
}

func (c *Config) Address() (string) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s:%d", c.Host, c.Port)
	return buf.String()
}

func (c *Config) SetJWTExpiration(t time.Duration) {
	c.JWTExpiration = int(t.Seconds())
}

func (c *Config) SetJWTExpirationFromSeconds(t int) {
	c.JWTExpiration = t
}

func (c *Config) GetJWTExpiration() (time.Duration) {
	return time.Duration(c.JWTExpiration) * time.Second
}

func NewConfig() (*Config) {
	return &Config{
		Port: 5000,
		Host: "localhost",
		JWTSecret: "secret",
		JWTExpiration: 3600,
		Database: Database{
			Host: "localhost",
			Port: 5432,
			Username: "postgres",
			Password: "postgres",
			Name: "azar",},
	}
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	config.SetJWTExpiration(time.Duration(config.JWTExpiration) * time.Second)
	return &config, nil
}

func LoadConfigFromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("azar", &config)
	if err != nil {
		return nil, err
	}

	config.Database.Host = config.DatabaseHost
	config.Database.Port = config.DatabasePort
	config.Database.Username = config.DatabaseUsername
	config.Database.Password = config.DatabasePassword
	config.Database.Name = config.DatabaseName

	return &config, nil
}