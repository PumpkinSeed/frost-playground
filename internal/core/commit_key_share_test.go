package core

import (
	"github.com/bytemare/secret-sharing/keys"
	"testing"

	"github.com/bytemare/ecc"
)

func TestDecodeVerificationKey(t *testing.T) {
	var data = "02c0d2b6a57aac851e27dee6541946fb74ebf0b21ab51a0e53fa5291e643350941"

	verificationKey := ecc.Secp256k1Sha256.Base()
	if err := verificationKey.DecodeHex(data); err != nil {
		t.Error(err)
	}
}

func TestDecodeKeyShare(t *testing.T) {
	var data = "0701000000000003b67298da44f8bca065974d4c7baa038668a9bc3cc108b7c099207108f376d7aa0d6205eae4e272b12d683e024921e69f40449962d3ff4859beb67394e27a1d5503de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0"

	var secretKeyShare = keys.KeyShare{}
	if err := secretKeyShare.DecodeHex(data); err != nil {
		t.Error(err)
	}
}
