package systemd

var WireguardPeerTemplate string = "" +
	"[WireGuardPeer]\n" +
	"PublicKey = __PEERPUBLICKEY__\n" +
	"AllowedIPs = __PEERALLOWEDIPS__\n" +
	"Endpoint = __PEERENDPOINT__\n" +
	"PersistentKeepalive = __KEEPALIVE__\n" +
	"RouteTable = __PEERROUTETABLE__\n" +
	"RouteMetric = __PEERROUTEMETRIC__\n"
