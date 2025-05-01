package core

import (
	"context"
	"encoding/json"
	"github.com/bytemare/frost"
	"io"
	"log/slog"
	"net/http"
)

type SignKeyShareRequest struct {
	Group       string   `json:"group"`
	Message     string   `json:"message"`
	Signer      string   `json:"signer"`
	Commitments []string `json:"commitments"`
}

type SignKeyShareResponse struct {
	Signature string `json:"signature"`
}

func SignKeyShareHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req SignKeyShareRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request", slog.Any("error", errJSON))
		http.Error(w, errJSON.Error(), http.StatusBadRequest)
		return
	}

	responseStruct, err := signKeyShare(ctx, req)

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

func signKeyShare(_ context.Context, req SignKeyShareRequest) (SignKeyShareResponse, error) {
	var participant = &frost.Signer{}
	if err := participant.DecodeHex(req.Signer); err != nil {
		slog.Error("failed to decode hex", "error", err)
		return SignKeyShareResponse{}, err
	}

	var commitments []*frost.Commitment
	for _, commitmentHex := range req.Commitments {
		var commitment = &frost.Commitment{}
		if err := commitment.DecodeHex(commitmentHex); err != nil {
			slog.Error("failed to decode hex", "error", err)
			return SignKeyShareResponse{}, err
		}

		commitments = append(commitments, commitment)
	}

	sigShare, err := participant.Sign([]byte(req.Message), commitments)
	if err != nil {
		slog.Error("failed to sign message", "error", err)
		return SignKeyShareResponse{}, err
	}

	return SignKeyShareResponse{
		Signature: sigShare.Hex(),
	}, nil
}
