package server

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
	Threshold     int    `json:"threshold"`
	Total         int    `json:"total"`
}

type SplitPrivateKeyResponse struct {
}
