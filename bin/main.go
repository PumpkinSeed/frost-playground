package main

import (
	"log/slog"

	"github.com/PumpkinSeed/frost-playground/pkg/server"
	"github.com/bytemare/ecc"
	"github.com/bytemare/frost"
	"github.com/bytemare/frost/debug"
	"github.com/bytemare/secret-sharing/keys"
)

func main() {
	server.Run()

	g := ecc.Secp256k1Sha256

	scalar := g.NewScalar().Random()

	slog.Info("random scalar", "scalar", scalar.Hex())

	scalar2 := g.NewScalar()
	if err := scalar2.DecodeHex(scalar.Hex()); err != nil {
		slog.Error("failed to decode hex", "error", err)
	}

	slog.Info("decoded scalar", "scalar", scalar2.Hex())

	secretKeyShares, verificationKey, _ := debug.TrustedDealerKeygen(frost.Secp256k1, scalar2, 4, 7)

	for _, sk := range secretKeyShares {
		slog.Info("secret key share", "ID", sk.ID, "secret", sk.Secret.Hex())
	}

	publicKeyShares := make([]*keys.PublicKeyShare, len(secretKeyShares))
	for i, sk := range secretKeyShares {
		publicKeyShares[i] = sk.Public()
	}

	configuration := &frost.Configuration{
		Ciphersuite:           frost.Secp256k1,
		Threshold:             4,
		MaxSigners:            7,
		VerificationKey:       verificationKey,
		SignerPublicKeyShares: publicKeyShares,
	}

	if err := configuration.Init(); err != nil {
		slog.Error("failed to initialize configuration", "error", err)
	}
}
