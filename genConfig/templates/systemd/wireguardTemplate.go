package systemd

var WireguardTemplate string = "\n" +
	"[WireGuard]\n" +
	"ListenPort = __LISTENPORT__\n" +
	"PrivateKey = __SERVERPRIVATEKEY__\n" +
	"FirewallMark = __FIREWALLMARK__\n" +
	"RouteTable = __SERVERROUTETABLE__\n" +
	"RouteMetric = __SERVERROUTEMETRIC__\n"
