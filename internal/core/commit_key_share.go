package core

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bytemare/secret-sharing/keys"
)

type CommitKeyShareRequest struct {
	Group           string   `json:"group"`
	VerificationKey string   `json:"verification_key"`
	SecretKeyShare  string   `json:"secret_key_share"`
	PublicKeyShares []string `json:"public_key_shares"`
	Threshold       uint16   `json:"threshold"`
	Total           uint16   `json:"total"`
}

type CommitKeyShareResponse struct {
	Signer     string `json:"signer"`
	Commitment string `json:"commitment"`
}

func CommitKeyShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req CommitKeyShareRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request", slog.Any("error", errJSON))
		http.Error(w, errJSON.Error(), http.StatusBadRequest)
		return
	}

	responseStruct, err := commitKeyShare(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseStruct)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func commitKeyShare(ctx context.Context, req CommitKeyShareRequest) (CommitKeyShareResponse, error) {
	configuration, err := composeConfiguration(ctx, composeConfigurationParams{
		VerificationKey: req.VerificationKey,
		PublicKeyShares: req.PublicKeyShares,
		Threshold:       req.Threshold,
		Total:           req.Total,
	})
	if err != nil {
		return CommitKeyShareResponse{}, err
	}

	var secretKeyShare = keys.KeyShare{}
	if err := secretKeyShare.DecodeHex(req.SecretKeyShare); err != nil {
		slog.ErrorContext(ctx, "failed to decode hex of secret", slog.Any("error", err))
		return CommitKeyShareResponse{}, err
	}

	signer, err := configuration.Signer(&secretKeyShare)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create signer", slog.Any("error", err))
		return CommitKeyShareResponse{}, err
	}

	commitment := signer.Commit()

	return CommitKeyShareResponse{
		Signer:     signer.Hex(),
		Commitment: commitment.Hex(),
	}, nil
}
