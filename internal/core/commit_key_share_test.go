package core

import (
	"testing"

	"github.com/bytemare/ecc"
	"github.com/bytemare/secret-sharing/keys"
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

func TestCommitKeyShare(t *testing.T) {
	ctx := t.Context()
	var req = CommitKeyShareRequest{
		VerificationKey: "03de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0",
		SecretKeyShare:  "0701000000000003b67298da44f8bca065974d4c7baa038668a9bc3cc108b7c099207108f376d7aa0d6205eae4e272b12d683e024921e69f40449962d3ff4859beb67394e27a1d5503de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0",
		PublicKeyShares: []string{
			"0701000200000003b67298da44f8bca065974d4c7baa038668a9bc3cc108b7c099207108f376d7aa03de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0023d4f1d8fa8314eed421308b292086072b05a31a77aa43775ad0d6a320e977e18",
			"0702000200000003febe57b6416c0d370288d6199c453d2fcb1846d30f11060cb00a4c6dbf0dcd0b03de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0023d4f1d8fa8314eed421308b292086072b05a31a77aa43775ad0d6a320e977e18",
			"0703000200000002332afb371142b75d61c280f023cbac075f62666d2a384898fce82369f067fb5103de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0023d4f1d8fa8314eed421308b292086072b05a31a77aa43775ad0d6a320e977e18",
			"0704000200000002c7bcf50cc96e245d668c2dede07836e9d0dd28ddd58994b55f70d7aad4977ff403de86ebf0b7aa1c9a3214ac00cb5036b430f7d4cce6849a825afbb617af3504b0023d4f1d8fa8314eed421308b292086072b05a31a77aa43775ad0d6a320e977e18",
		},
	}

	resp, err := commitKeyShare(ctx, req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.Commitment)
	t.Log(resp.Signer)
}
