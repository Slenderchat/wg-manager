package client

type Client struct {
	Name        string
	PrivateKey  string
	PublicKey   string
	IP          string
	Endpoint    string
	RouteTable  string
	RouteMetric string
}

// Generates new client object with new pair of private and public keys
func NewClient(name string) *Client {
	priv, pub := genSecrets()
	ip := randomIP(name)
	return &Client{Name: name, PrivateKey: priv, PublicKey: pub, IP: ip}
}
