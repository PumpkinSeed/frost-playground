package core

import (
	"fmt"
	"math/big"
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
