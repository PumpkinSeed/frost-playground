package core

import (
	"github.com/bytemare/ecc"
	"github.com/bytemare/frost"
	"log/slog"
	"testing"
)

func TestVerification(t *testing.T) {
	message := []byte("test")
	var aggSig = &frost.Signature{}
	if err := aggSig.DecodeHex("0703557589164a02309d191d386791682232ddf2c24bfc34f08aec627419ccfe2a54c324d9c865ea68416c20c2ef87f497fedc213f59798afef2c69c7b5d92a05057"); err != nil {
		t.Fatal(err)
	}

	verificationKey := ecc.Secp256k1Sha256.Base()
	if err := verificationKey.DecodeHex("0242bf8dcf4eb95ba7977f0e0c385c205c68e10c1f2047133b3a381d44f980351f"); err != nil {
		t.Fatal(err)
	}

	err := frost.VerifySignature(frost.Secp256k1, message, aggSig, verificationKey)
	if err != nil {
		slog.Error("failed to verify signature", "error", err)
		return
	} else {
		slog.Info("signature verified")
	}
}
