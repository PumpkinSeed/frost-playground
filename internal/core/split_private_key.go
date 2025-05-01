package core

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bytemare/ecc"
	secretsharing "github.com/bytemare/secret-sharing"
	"github.com/bytemare/secret-sharing/keys"
)

type SplitPrivateKeyRequest struct {
	PrivateKeyHex string `json:"private_key_hex"`
	Threshold     uint16 `json:"threshold"`
	Total         uint16 `json:"total"`
}

type SplitPrivateKeyResponse struct {
	KeyShars         []KeyShare `json:"key_shars"`
	VerificationKeys []string   `json:"verification_keys"`
}

type KeyShare struct {
	Group   string         `json:"group"`
	SK      string         `json:"sk"`
	Public  string         `json:"public"`
	Details *keys.KeyShare `json:"details"`
}

func SplitPrivateKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req SplitPrivateKeyRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request", slog.Any("error", errJSON))
		http.Error(w, errJSON.Error(), http.StatusBadRequest)
		return
	}

	// Decode scalar, hard-coded for now
	g := ecc.Secp256k1Sha256
	scalar := g.NewScalar()
	if err := scalar.DecodeHex(req.PrivateKeyHex); err != nil {
		slog.Error("failed to decode hex", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	privateKeyShares, poly, err := secretsharing.ShardReturnPolynomial(g, scalar, req.Threshold, req.Total)
	if err != nil {
		slog.ErrorContext(ctx, "failed to shard polynomial", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// This will be the verification key
	coms := secretsharing.Commit(g, poly)

	shares := make([]KeyShare, req.Total)
	for i, sk := range privateKeyShares {
		public := keys.PublicKeyShare{
			PublicKey:     g.Base().Multiply(sk.Secret),
			VssCommitment: coms,
			ID:            sk.ID,
			Group:         g,
		}
		shares[i] = KeyShare{
			Group:  GroupToString(g),
			SK:     sk.Secret.Hex(),
			Public: public.Hex(),
			Details: &keys.KeyShare{
				Secret:          sk.Secret,
				VerificationKey: coms[0],
				PublicKeyShare:  public,
			},
		}
	}

	var verificationKeys []string
	for _, com := range coms {
		verificationKeys = append(verificationKeys, com.Hex())
	}

	response, err := json.Marshal(SplitPrivateKeyResponse{
		KeyShars:         shares,
		VerificationKeys: verificationKeys,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
