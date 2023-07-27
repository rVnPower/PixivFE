package configs

import (
	"errors"
	"os"
)

var Token, BaseURL, Port, UserAgent, ProxyServer string

func parseEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)

	if !ok {
		return value, errors.New("Failed to get environment variable" + key)
	}

	return value, nil
}

func parseEnvWithDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return defaultValue
	}

	return value
}

func ParseConfig() error {
	var err error

	Token, err = parseEnv("PIXIVFE_TOKEN")

	if err != nil {
		return err
	}

	BaseURL = parseEnvWithDefault("PIXIVFE_BASEURL", "localhost")
	Port = parseEnvWithDefault("PIXIVFE_PORT", "8080")
	UserAgent = parseEnvWithDefault("PIXIVFE_USERAGENT", "Mozilla/5.0")
	ProxyServer = parseEnvWithDefault("PIXIVFE_IMAGEPROXY", "px2.rainchan.win")

	return nil
}
