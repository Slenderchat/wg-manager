package genConfig

import (
	"fmt"
	"regexp"
	"strings"
	"wg-manager/client"
	"wg-manager/genConfig/templates/systemd"
)

func genClientNetDev(config *serverConfig) (section string) {
	section = systemd.NetDevTemplate
	section = strings.ReplaceAll(section, "__WGINTNAME__", "vpn")
	section = strings.ReplaceAll(section, "__WGINTDESC__", config.InterfaceDescription+`.`+config.Domain)
	return
}

func genClientWireguard(config *serverConfig, client *client.Client) (section string) {
	section = systemd.WireguardTemplate
	section = strings.ReplaceAll(section, "__LISTENPORT__", config.Port)
	section = strings.ReplaceAll(section, "__SERVERPRIVATEKEY__", client.PrivateKey)
	rTemplate := `(?m)[\r\n]+^.*__SUBSTR__.*$`
	if config.FirewallMark != "" {
		section = strings.ReplaceAll(section, "__FIREWALLMARK__", config.FirewallMark)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__FIREWALLMARK__"))
		section = r.ReplaceAllString(section, "")
	}
	if client.RouteTable != "" {
		section = strings.ReplaceAll(section, "__SERVERROUTETABLE__", client.RouteTable)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__SERVERROUTETABLE__"))
		section = r.ReplaceAllString(section, "")
	}
	if client.RouteMetric != "" {
		section = strings.ReplaceAll(section, "__SERVERROUTEMETRIC__", client.RouteMetric)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__SERVERROUTEMETRIC__"))
		section = r.ReplaceAllString(section, "")
	}
	return
}

func genClientWireguardPeer(config *serverConfig) (section string, err error) {
	section = fmt.Sprintf("\n// %v\n", config.InterfaceDescription+`.`+config.Domain)
	section += systemd.WireguardPeerTemplate
	section = strings.ReplaceAll(section, "__PEERPUBLICKEY__", config.PublicKey)
	section = strings.ReplaceAll(section, "__PEERALLOWEDIPS__", "0.0.0.0/0")
	section = strings.ReplaceAll(section, "__PEERENDPOINT__", "hel.sl-chat.ru:"+config.Port)
	rTemplate := `(?m)[\r\n]+^.*__SUBSTR__.*$`
	if config.Keepalive != "" {
		section = strings.ReplaceAll(section, "__KEEPALIVE__", config.Keepalive)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__KEEPALIVE__"))
		section = r.ReplaceAllString(section, "")
	}
	r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__PEERROUTETABLE__"))
	section = r.ReplaceAllString(section, "")
	r = regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__PEERROUTEMETRIC__"))
	section = r.ReplaceAllString(section, "")
	return
}

func GenClientConfig(client *client.Client) (config string, err error) {
	sConfig, e := readConfig()
	if e != nil {
		err = fmt.Errorf("failed when generating server config: %v", e)
		return
	}
	config = genClientNetDev(sConfig)
	config += genClientWireguard(sConfig, client)
	peers, e := genClientWireguardPeer(sConfig)
	if e != nil {
		err = fmt.Errorf("error generating server config: %v", e)
		config = ""
		return
	}
	config += peers
	return
}
