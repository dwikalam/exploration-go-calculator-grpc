package config

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	serverConfig serverConfig
}

func New() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("load godotenv failed: %v", err)
	}

	var (
		serverConfig serverConfig

		err error
	)

	serverConfig, err = newServerConfig()
	if err != nil {
		return Config{}, fmt.Errorf("creating new server config failed: %v", err)
	}

	return Config{
		serverConfig: serverConfig,
	}, nil
}

func (c *Config) GetServerNetwork() string {
	return c.serverConfig.network
}

func (c *Config) GetServerHost() string {
	return c.serverConfig.host
}

func (c *Config) GetServerPort() int {
	return c.serverConfig.port
}

func (c *Config) GetServerAddress() string {
	return net.JoinHostPort(c.GetServerHost(), strconv.Itoa(c.GetServerPort()))
}

type serverConfig struct {
	network string
	host    string
	port    int
}

func newServerConfig() (serverConfig, error) {
	network := getEnvValue("SERVER_NETWORK", "tcp")
	host := getEnvValue("SERVER_HOST", "localhost")
	port, err := strconv.Atoi(getEnvValue("SERVER_PORT", "8080"))
	if err != nil {
		return serverConfig{}, fmt.Errorf("failed converting port value to int from env: %v", err)
	}

	return serverConfig{
		network: network,
		host:    host,
		port:    port,
	}, nil
}

func getEnvValue(key string, placeholder string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return placeholder
}
