package core

import (
	"encoding/json"
	"github.com/bytemare/frost"
	"io"
	"log/slog"
	"net/http"
)

type AggregateSignaturesRequest struct {
	Signatures      []string `json:"signatures"`
	Message         string   `json:"message"`
	Commitments     []string `json:"commitments"`
	VerificationKey string   `json:"verification_key"`
	PublicKeyShares []string `json:"public_key_shares"`
	Threshold       uint16   `json:"threshold"`
	Total           uint16   `json:"total"`
}

type AggregateSignaturesResponse struct {
	Signature string `json:"signature"`
}

func AggregateSignaturesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req AggregateSignaturesRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request", slog.Any("error", errJSON))
		http.Error(w, errJSON.Error(), http.StatusBadRequest)
		return
	}

	configuration, err := composeConfiguration(ctx, composeConfigurationParams{
		VerificationKey: req.VerificationKey,
		PublicKeyShares: req.PublicKeyShares,
		Threshold:       req.Threshold,
		Total:           req.Total,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var signatures []*frost.SignatureShare
	for _, sigHex := range req.Signatures {
		var sig frost.SignatureShare
		if err := sig.DecodeHex(sigHex); err != nil {
			slog.Error("failed to decode hex", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		signatures = append(signatures, &sig)
	}

	var commitments []*frost.Commitment
	for _, commitHex := range req.Commitments {
		var commit frost.Commitment
		if err := commit.DecodeHex(commitHex); err != nil {
			slog.Error("failed to decode hex", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		commitments = append(commitments, &commit)
	}

	aggSignature, err := configuration.AggregateSignatures([]byte(req.Message), signatures, commitments, true)
	if err != nil {
		slog.Error("failed to aggregate signatures", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(AggregateSignaturesResponse{Signature: aggSignature.Hex()})
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
