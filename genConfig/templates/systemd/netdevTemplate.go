package systemd

var NetDevTemplate string = "" +
	"[NetDev]\n" +
	"Name = __WGINTNAME__\n" +
	"Kind = wireguard\n" +
	"Description = __WGINTDESC__\n"
