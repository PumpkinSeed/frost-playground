package core

import (
	"encoding/json"
	"github.com/bytemare/ecc"
	"github.com/bytemare/frost"
	"io"
	"log/slog"
	"net/http"
)

type SignatureVerificationRequest struct {
	Message         string `json:"message"`
	Signature       string `json:"signature"`
	VerificationKey string `json:"verification_key"`
}

type SignatureVerificationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SignatureVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req SignatureVerificationRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request", slog.Any("error", errJSON))
		http.Error(w, errJSON.Error(), http.StatusBadRequest)
		return
	}

	var aggSig = &frost.Signature{}
	if err := aggSig.DecodeHex(req.Signature); err != nil {
		slog.ErrorContext(ctx, "failed to decode hex of signature", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	verificationKey := ecc.Secp256k1Sha256.Base()
	if err := verificationKey.DecodeHex(req.VerificationKey); err != nil {
		slog.ErrorContext(ctx, "failed to decode hex of verification key", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = frost.VerifySignature(frost.Secp256k1, []byte(req.Message), aggSig, verificationKey)
	if err != nil {
		slog.Error("failed to verify signature", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(SignatureVerificationResponse{
		Status:  "valid",
		Message: "Valid signature",
	})
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
