package genConfig

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"wg-manager/client"
	"wg-manager/genConfig/templates/systemd"
)

func genServerNetDev(config *serverConfig) (section string) {
	section = systemd.NetDevTemplate
	section = strings.ReplaceAll(section, "__WGINTNAME__", config.InterfaceName)
	section = strings.ReplaceAll(section, "__WGINTDESC__", config.InterfaceDescription+`.`+config.Domain)
	return
}

func genServerWireguard(config *serverConfig) (section string) {
	section = systemd.WireguardTemplate
	section = strings.ReplaceAll(section, "__LISTENPORT__", config.Port)
	section = strings.ReplaceAll(section, "__SERVERPRIVATEKEY__", config.PrivateKey)
	rTemplate := `(?m)[\r\n]+^.*__SUBSTR__.*$`
	if config.FirewallMark != "" {
		section = strings.ReplaceAll(section, "__FIREWALLMARK__", config.FirewallMark)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__FIREWALLMARK__"))
		section = r.ReplaceAllString(section, "")
	}
	if config.RouteTable != "" {
		section = strings.ReplaceAll(section, "__SERVERROUTETABLE__", config.RouteTable)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__SERVERROUTETABLE__"))
		section = r.ReplaceAllString(section, "")
	}
	if config.RouteMetric != "" {
		section = strings.ReplaceAll(section, "__SERVERROUTEMETRIC__", config.RouteMetric)
	} else {
		r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__SERVERROUTEMETRIC__"))
		section = r.ReplaceAllString(section, "")
	}
	return
}

func genServerWireguardPeers(config *serverConfig, clients *[]*client.Client) (section string, err error) {
	for _, v := range *clients {
		hostname, e := os.Hostname()
		if e != nil {
			err = fmt.Errorf("error generating peer configs: %v", e)
			return
		}
		if strings.Contains(v.Name, hostname) {
			continue
		}
		peerSection := fmt.Sprintf("\n// %v\n", v.Name+`.`+config.Domain)
		peerSection += systemd.WireguardPeerTemplate
		peerSection = strings.ReplaceAll(peerSection, "__PEERPUBLICKEY__", v.PublicKey)
		peerSection = strings.ReplaceAll(peerSection, "__PEERALLOWEDIPS__", v.IP+"/32")
		rTemplate := `(?m)[\r\n]+^.*__SUBSTR__.*$`
		if v.Endpoint != "" {
			peerSection = strings.ReplaceAll(peerSection, "__PEERENDPOINT__", v.Endpoint)
		} else {
			r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__PEERENDPOINT__"))
			peerSection = r.ReplaceAllString(peerSection, "")
		}
		if config.Keepalive != "" {
			peerSection = strings.ReplaceAll(peerSection, "__KEEPALIVE__", config.Keepalive)
		} else {
			r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__KEEPALIVE__"))
			peerSection = r.ReplaceAllString(peerSection, "")
		}
		if v.RouteTable != "" {
			peerSection = strings.ReplaceAll(peerSection, "__PEERROUTETABLE__", v.RouteTable)
		} else {
			r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__PEERROUTETABLE__"))
			peerSection = r.ReplaceAllString(peerSection, "")
		}
		if v.RouteMetric != "" {
			peerSection = strings.ReplaceAll(peerSection, "__PEERROUTEMETRIC__", v.RouteMetric)
		} else {
			r := regexp.MustCompile(strings.ReplaceAll(rTemplate, "__SUBSTR__", "__PEERROUTEMETRIC__"))
			peerSection = r.ReplaceAllString(peerSection, "")
		}
		section += peerSection
	}
	return
}

func GenServerConfig(clients *[]*client.Client) (config string, err error) {
	sConfig, e := readConfig()
	if e != nil {
		err = fmt.Errorf("failed when generating server config: %v", e)
		return
	}
	config = genServerNetDev(sConfig)
	config += genServerWireguard(sConfig)
	peers, e := genServerWireguardPeers(sConfig, clients)
	if e != nil {
		err = fmt.Errorf("error generating server config: %v", e)
		config = ""
		return
	}
	config += peers
	return
}
