package server

import (
	"github.com/bytemare/ecc"
	"github.com/bytemare/secret-sharing/keys"
	"strings"
)

const (
	ristretto255Sha512 = "Ristretto255-SHA512"
	p256Sha256         = "P256-SHA256"
	p384Sha384         = "P384-SHA384"
	p521Sha512         = "P521-SHA512"
	edwards25519Sha512 = "Edwards25519-SHA512"
	secp256k1Sha256    = "Secp256k1-SHA256"
)

type CreatePrivateKeyRequest struct {
	PrimeOrderGroup string `json:"prime_order_group"`
}

type CreatePrivateKeyResponse struct {
	PrivateKeyHex string `json:"private_key_hex"`
}

type SplitPrivateKeyRequest struct {
	PrivateKeyHex string `json:"private_key_hex"`
	Threshold     uint16 `json:"threshold"`
	Total         uint16 `json:"total"`
}

type SplitPrivateKeyResponse struct {
	KeyShars []KeyShare `json:"key_shars"`
	Comms    []string   `json:"comms"`
}

type KeyShare struct {
	Group   string `json:"group"`
	SK      string `json:"sk"`
	Public  string `json:"public"`
	Details *keys.KeyShare
}

func GroupToString(group ecc.Group) string {
	switch group {
	case ecc.Ristretto255Sha512:
		return ristretto255Sha512
	case ecc.P256Sha256:
		return p256Sha256
	case ecc.P384Sha384:
		return p384Sha384
	case ecc.P521Sha512:
		return p521Sha512
	case ecc.Edwards25519Sha512:
		return edwards25519Sha512
	case ecc.Secp256k1Sha256:
		return secp256k1Sha256
	}
	return ""
}

func StringToGroup(group string) ecc.Group {
	switch strings.ToLower(group) {
	case strings.ToLower(ristretto255Sha512):
		return ecc.Ristretto255Sha512
	case strings.ToLower(p256Sha256):
		return ecc.P256Sha256
	case strings.ToLower(p384Sha384):
		return ecc.P384Sha384
	case strings.ToLower(p521Sha512):
		return ecc.P521Sha512
	case strings.ToLower(edwards25519Sha512):
		return ecc.Edwards25519Sha512
	case strings.ToLower(secp256k1Sha256):
		return ecc.Secp256k1Sha256
	}
	return 0
}
