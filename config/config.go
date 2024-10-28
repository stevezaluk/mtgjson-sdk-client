package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	IPAddress string `json:"ip_address"`
	Port      string `json:"port"`
	UseSSL    bool   `json:"use_ssl"`
}

func (c Config) BuildUri() string {
	prefix := "http://"
	if c.UseSSL {
		prefix = "https://"
	}

	ret := []string{prefix, c.IPAddress, ":", c.Port}
	return strings.Join(ret, "")
}

func ParseConfig(path string) Config {
	// convert the path to an absolute path
	if strings.HasPrefix(path, "~") {
		home := os.Getenv("HOME")
		path = strings.ReplaceAll(path, "~", home)
	}

	file, err := os.Open(path)
	if err == os.ErrNotExist {
		fmt.Println("[error] Failed to find config file: ", path)
		panic(1)
	}

	if err == os.ErrPermission {
		fmt.Println("[error] Invalid permissions for opening the config file: ", path)
		panic(1)
	}

	bytes, _ := io.ReadAll(file)

	var ret Config
	errors := json.Unmarshal([]byte(bytes), &ret)
	if errors != nil {
		fmt.Println("[error] Failed to parse config file: ", errors.Error())
		panic(1)
	}

	return ret
}

func ParseEnv() (Config, []string) {
	envs := map[string]string{
		"MTGJSON_IP":      os.Getenv("MTGJSON_IP"),
		"MTGJSON_PORT":    os.Getenv("MTGJSON_PORT"),
		"MTGJSON_USE_SSL": os.Getenv("MTGJSON_USE_SSL"),
	}

	var invalid []string
	for key, value := range envs {
		if value == "" {
			invalid = append(invalid, key)
		}
	}

	use_ssl, err := strconv.ParseBool(envs["MTGJSON_USE_SSL"])
	if err != nil {
		invalid = append(invalid, "MTGJSON_USE_SSL")
	}

	if len(invalid) != 0 {
		return Config{}, invalid
	}

	ret := Config{
		IPAddress: envs["MTGJSON_IP"],
		Port:      envs["MTGJSON_PORT"],
		UseSSL:    use_ssl,
	}

	return ret, invalid
}
