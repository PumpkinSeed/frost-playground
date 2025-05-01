package core

import (
	"context"
	"log/slog"

	"github.com/bytemare/ecc"
	"github.com/bytemare/frost"
	"github.com/bytemare/secret-sharing/keys"
)

type composeConfigurationParams struct {
	Group           string   `json:"group"`
	VerificationKey string   `json:"verification_key"`
	PublicKeyShares []string `json:"public_key_shares"`
	Threshold       uint16   `json:"threshold"`
	Total           uint16   `json:"total"`
}

func composeConfiguration(ctx context.Context, req composeConfigurationParams) (*frost.Configuration, error) {
	verificationKey := ecc.Secp256k1Sha256.Base()
	if err := verificationKey.DecodeHex(req.VerificationKey); err != nil {
		slog.ErrorContext(ctx, "failed to decode hex", slog.Any("error", err))
		return nil, err
	}

	var publicKeyShares []*keys.PublicKeyShare
	for _, data := range req.PublicKeyShares {
		var publicKeyShare = keys.PublicKeyShare{}
		if err := publicKeyShare.DecodeHex(data); err != nil {
			slog.ErrorContext(ctx, "failed to decode hex of public key share", slog.Any("error", err))
			return nil, err
		}
		publicKeyShares = append(publicKeyShares, &publicKeyShare)
	}

	return &frost.Configuration{
		Ciphersuite:           frost.Secp256k1,
		Threshold:             req.Threshold,
		MaxSigners:            req.Total,
		VerificationKey:       verificationKey,
		SignerPublicKeyShares: publicKeyShares,
	}, nil
}
