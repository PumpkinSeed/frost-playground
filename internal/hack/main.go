package main

import (
	"fmt"
	"log/slog"

	"github.com/bytemare/ecc"
	"github.com/bytemare/frost"
	"github.com/bytemare/frost/debug"
	"github.com/bytemare/secret-sharing/keys"
)

func main() {
	g := ecc.Secp256k1Sha256

	scalar := g.NewScalar().Random()

	slog.Info("random scalar", "scalar", scalar.Hex())

	scalar2 := g.NewScalar()
	if err := scalar2.DecodeHex(scalar.Hex()); err != nil {
		slog.Error("failed to decode hex", "error", err)
	}

	slog.Info("decoded scalar", "scalar", scalar2.Hex())

	secretKeyShares, verificationKey, _ := debug.TrustedDealerKeygen(frost.Secp256k1, scalar2, 4, 7)
	fmt.Println("verification key", verificationKey.Hex())

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

	// Create participants
	var participants []*frost.Signer
	for _, ks := range secretKeyShares[:4] {
		ksHex := ks.Hex()

		var secretKeyShare = keys.KeyShare{}
		if err := secretKeyShare.DecodeHex(ksHex); err != nil {
			slog.Error("failed to decode hex", "error", err)
			return
		}

		signer, err := configuration.Signer(&secretKeyShare)
		if err != nil {
			slog.Error("failed to create signer", "error", err)
			return
		}

		fmt.Println("signer", signer.Hex())
		participants = append(participants, signer)
	}

	// Create commitments of the participants
	var commitments []*frost.Commitment
	for _, participant := range participants {
		commit := participant.Commit()
		fmt.Println("commitment", commit.Hex())
		commitments = append(commitments, commit)
	}

	message := []byte("Hello, world!")

	var sigShares []*frost.SignatureShare
	for _, participant := range participants[:4] {
		parHex := participant.Hex()
		fmt.Println("participant", parHex)

		var par = &frost.Signer{}
		if err := par.DecodeHex(parHex); err != nil {
			slog.Error("failed to decode hex", "error", err)
			return
		}

		sigShare, err := par.Sign(message, commitments)
		if err != nil {
			slog.Error("failed to sign message", "error", err)
			return
		}

		sigShares = append(sigShares, sigShare)
	}

	// Aggregate the signature shares
	aggSig, err := configuration.AggregateSignatures(message, sigShares, commitments, true)
	if err != nil {
		slog.Error("failed to aggregate signatures", "error", err)
		return
	}
	fmt.Println(aggSig.Hex())

	err = frost.VerifySignature(frost.Secp256k1, message, aggSig, verificationKey)
	if err != nil {
		slog.Error("failed to verify signature", "error", err)
		return
	} else {
		slog.Info("signature verified")
	}
}
