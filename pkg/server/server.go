package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bytemare/ecc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/private-key", CreatePrivateKey)
	r.Post("/split-private-key", SplitPrivateKey)

	slog.Info("starting server", slog.String("address", ":3000"))
	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func CreatePrivateKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dataRaw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var group ecc.Group
	var req CreatePrivateKeyRequest
	if errJSON := json.Unmarshal(dataRaw, &req); errJSON != nil {
		slog.WarnContext(ctx, "failed to unmarshal request, setup default with secp256k1", slog.Any("error", errJSON))
		group = ecc.Secp256k1Sha256
	} else {
		switch strings.ToLower(req.PrimeOrderGroup) {
		case strings.ToLower(ristretto255Sha512):
			group = ecc.Ristretto255Sha512
		case strings.ToLower(p256Sha256):
			group = ecc.P256Sha256
		case strings.ToLower(p384Sha384):
			group = ecc.P384Sha384
		case strings.ToLower(p521Sha512):
			group = ecc.P521Sha512
		case strings.ToLower(edwards25519Sha512):
			group = ecc.Edwards25519Sha512
		case strings.ToLower(secp256k1Sha256):
			group = ecc.Secp256k1Sha256
		default:
			slog.ErrorContext(ctx, "invalid prime order group, setup default with secp256k1", slog.String("prime_order_group", req.PrimeOrderGroup))
		}
	}

	data, err := json.Marshal(CreatePrivateKeyResponse{
		PrivateKeyHex: group.NewScalar().Random().Hex(),
	})
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

func SplitPrivateKey(w http.ResponseWriter, r *http.Request) {

}
