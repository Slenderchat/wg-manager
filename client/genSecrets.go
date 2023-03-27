package client

import (
	"encoding/base64"
	"math/rand"

	"golang.org/x/crypto/curve25519"
)

// Generates and returns base64 encoded new public and private curve25519 keys
func genSecrets() (priv64 string, pub64 string) {
	var priv, pub [32]byte
	for i := range priv {
		priv[i] = byte(rand.Intn(256))
	}
	curve25519.ScalarBaseMult(&pub, &priv)
	priv64 = base64.StdEncoding.EncodeToString(priv[:])
	pub64 = base64.StdEncoding.EncodeToString(pub[:])
	return
}
