package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bytemare/ecc"
	secretsharing "github.com/bytemare/secret-sharing"
	"github.com/bytemare/secret-sharing/keys"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
		group = StringToGroup(req.PrimeOrderGroup)
		if group == 0 {
			slog.ErrorContext(ctx, "invalid prime order group, setup default with secp256k1", slog.String("prime_order_group", req.PrimeOrderGroup))
			group = ecc.Secp256k1Sha256
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
	if err := scalar.DecodeHex(scalar.Hex()); err != nil {
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

	// TODO
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

	response, err := json.Marshal(SplitPrivateKeyResponse{
		KeyShars: shares,
		Comms:    nil,
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
