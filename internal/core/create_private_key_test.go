package core

import (
	"fmt"
	"github.com/bytemare/secp256k1"
	"github.com/bytemare/secret-sharing/keys"
	"log/slog"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/bytemare/ecc"
	secretsharing "github.com/bytemare/secret-sharing"
)

type disallowEqual [0]func()

type Scalar struct {
	_ disallowEqual
	*secp256k1Scalar
}

type secp256k1Scalar struct {
	scalar *internalScalar
}

type internalScalar struct {
	_      disallowEqual
	scalar big.Int
}

func TestScalar(t *testing.T) {
	var hexData = "d903a3359bb6ec8ffacd7f15c382f21da19a971adb9787667adb5e8a09b44723"

	bigNum := new(big.Int)

	// SetString returns (result *Int, success bool)
	bigNum, ok := bigNum.SetString(hexData, 16)
	if !ok {
		fmt.Println("Invalid number string")
		return
	}
	fmt.Println(bigNum.String())
}

func TestSplit(t *testing.T) {
	var hexData = "d903a3359bb6ec8ffacd7f15c382f21da19a971adb9787667adb5e8a09b44723"
	var threshold uint16 = 2
	var total uint16 = 4

	g := ecc.Secp256k1Sha256
	scalar := g.NewScalar()
	if err := scalar.DecodeHex(hexData); err != nil {
		t.Fatal(err)
	}

	privateKeyShares, poly, err := secretsharing.ShardReturnPolynomial(g, scalar, threshold, total)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(privateKeyShares)
	t.Log(poly)
}

func TestSplitWithSmallNumbers(t *testing.T) {
	var threshold uint16 = 3
	var total uint16 = 5

	var privateKey uint64 = 5433

	var g = ecc.Secp256k1Sha256

	// Private key
	scalar := g.NewScalar()
	scalar.SetUInt64(privateKey)

	// Create polynomial
	polynomial := make(secretsharing.Polynomial, threshold)
	i := uint16(0)
	if scalar != nil {
		polynomial[0] = scalar.Copy()
		i++
	}
	polynomial[i] = g.NewScalar().SetUInt64(1234)
	i++
	polynomial[i] = g.NewScalar().SetUInt64(7734)
	//for ; i < threshold; i++ {
	//	polynomial[i] = g.NewScalar().SetUInt64(random(10, 3000))
	//}

	// Verification key
	basePoint := secp256k1.Element{}
	verificationKey := g.Base().Multiply(polynomial[0])

	// Key shares
	secretKeyShares := make([]*keys.KeyShare, total)

	for j := uint16(1); j <= total; j++ {
		ids := g.NewScalar().SetUInt64(uint64(j))
		value := polynomial[len(polynomial)-1].Copy()
		for k := len(polynomial) - 2; k >= 0; k-- {
			value.Multiply(ids)
			polyk := polynomial[k]
			value.Add(polyk)
		}

		yi := value

		keyShare := &keys.KeyShare{
			Secret:          yi,
			VerificationKey: verificationKey,
			PublicKeyShare: keys.PublicKeyShare{
				PublicKey:     g.Base().Multiply(yi),
				VssCommitment: nil,
				ID:            j,
				Group:         g,
			},
		}
		secretKeyShares[j-1] = keyShare
	}
	for i, share := range secretKeyShares {
		fmt.Printf("Share %d: %d\n", i+1, hexToUint64(share.Secret.Hex()))
	}
	//fmt.Println("Verification Key:", hexToUint64(verificationKey.Hex()))
	fmt.Println(10)
}

func hexToUint64(hexStr string) uint64 {
	// Remove optional "0x" or "0X" prefix
	if len(hexStr) >= 2 && (hexStr[0:2] == "0x" || hexStr[0:2] == "0X") {
		hexStr = hexStr[2:]
	}
	fmt.Println("Hex String:", hexStr)

	value, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("invalid hex string: %s", err))
		return 0
	}
	return value
}

func random(min, max int) uint64 {
	return uint64(rand.Intn(max-min) + min)
}
