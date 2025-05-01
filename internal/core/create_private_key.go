package core

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bytemare/ecc"
)

type CreatePrivateKeyRequest struct {
	PrimeOrderGroup string `json:"prime_order_group"`
}

type CreatePrivateKeyResponse struct {
	PrivateKeyHex string `json:"private_key_hex"`
}

func CreatePrivateKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := CreatePrivateKey(ctx, dataRaw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreatePrivateKey(ctx context.Context, dataRaw []byte) (CreatePrivateKeyResponse, error) {
	var group ecc.Group
	var req CreatePrivateKeyRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.WarnContext(ctx, "failed to unmarshal request, setup default with secp256k1", slog.Any("error", errJSON))
		group = ecc.Secp256k1Sha256
	} else {
		group = StringToGroup(req.PrimeOrderGroup)
		if group == 0 {
			slog.ErrorContext(ctx, "invalid prime order group, setup default with secp256k1", slog.String("prime_order_group", req.PrimeOrderGroup))
			group = ecc.Secp256k1Sha256
		}
	}

	return CreatePrivateKeyResponse{
		PrivateKeyHex: group.NewScalar().Random().Hex(),
	}, nil
}
