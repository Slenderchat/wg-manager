package genConfig

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"wg-manager/client"
)

type serverConfig struct {
	InterfaceName        string
	InterfaceDescription string
	Port                 string
	IP                   string
	Domain               string
	PrivateKey           string
	PublicKey            string
	Keepalive            string
	FirewallMark         string
	RouteTable           string
	RouteMetric          string
}

func (config *serverConfig) iterator() [][]string {
	return [][]string{
		{"InterfaceName", config.InterfaceName},
		{"InterfaceDescription", config.InterfaceDescription},
		{"Port", config.Port},
		{"IP", config.IP},
		{"Domain", config.Domain},
		{"PrivateKey", config.PrivateKey},
		{"PublicKey", config.PublicKey},
		{"Keepalive", config.Keepalive},
		{"FirewallMark", config.FirewallMark},
		{"RouteTable", config.RouteTable},
		{"RouteMetric", config.RouteMetric},
	}
}

func readConfig() (config *serverConfig, err error) {
	f, e := os.OpenFile("servers.json", os.O_CREATE|os.O_RDONLY, 0644)
	if e != nil {
		err = fmt.Errorf("failed to open server config file: %v", e)
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	e = dec.Decode(&config)
	if e != nil {
		if e.Error() == "EOF" {
			config = &serverConfig{}
			e = saveConfig(config)
			if e != nil {
				err = fmt.Errorf("failed to read server config file: %v", e)
				config = nil
				return
			}
			err = errors.New("your config file was unexistent, and was created for you - please adjust it")
			config = nil
			return
		} else {
			err = fmt.Errorf("failed to read server config file: %v", e)
			config = nil
			return
		}
	}
	for _, v := range config.iterator() {
		if v[1] == "" {
			switch v[0] {
			case "Port":
				var i uint64
				for {
					var fb [2]byte
					var b []byte
					rand.Read(fb[:])
					for _, v := range fb {
						b = append(b, v)
					}
					for i := 2; i < 8; i++ {
						b = append(b, 0)
					}
					i = binary.LittleEndian.Uint64(b)
					if i < 10000 {
						continue
					} else {
						break
					}
				}
				config.Port = fmt.Sprint(i)
				continue
			case "PrivateKey", "PublicKey":
				server := client.NewClient(config.InterfaceDescription)
				config.PrivateKey = server.PrivateKey
				config.PublicKey = server.PublicKey
				server = nil
				continue
			case "InterfaceDescription":
				config.InterfaceDescription, e = os.Hostname()
				if e != nil {
					err = fmt.Errorf("value of InterfaceDescription is not set, and ecountered error when trying to get machine hostname: %v", e)
					config = nil
					return
				}
			case "FirewallMark", "RouteTable", "RouteMetric", "Keepalive":
			default:
				err = fmt.Errorf("error when parsing server config file: required setting %v is unspecified or has no value, please adjust your config", v[0])
				saveConfig(config)
				config = nil
				return
			}
		}
	}
	e = saveConfig(config)
	if e != nil {
		err = fmt.Errorf("failed to read server config file: %v", e)
		config = nil
		return
	}
	return
}

func saveConfig(config *serverConfig) error {
	f, e := os.OpenFile("servers.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if e != nil {
		return fmt.Errorf("failed to open server config file: %v", e)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	e = enc.Encode(config)
	if e != nil {
		return fmt.Errorf("failed to read server config file: %v", e)
	}
	return nil
}
