package cryptostore

import (
	"github.com/pkg/errors"
	crypto "github.com/tepleton/go-crypto"
)

var (
	// GenEd25519 produces Ed25519 private keys
	GenEd25519 Generator = GenFunc(genEd25519)
	// GenSecp256k1 produces Secp256k1 private keys
	GenSecp256k1 Generator = GenFunc(genSecp256)
)

// Generator determines the type of private key the keystore creates
type Generator interface {
	Generate() crypto.PrivKey
}

// GenFunc is a helper to transform a function into a Generator
type GenFunc func() crypto.PrivKey

func (f GenFunc) Generate() crypto.PrivKey {
	return f()
}

func genEd25519() crypto.PrivKey {
	return crypto.GenPrivKeyEd25519()
}

func genSecp256() crypto.PrivKey {
	return crypto.GenPrivKeySecp256k1()
}

func getGenerator(algo string) (Generator, error) {
	switch algo {
	case crypto.NameEd25519:
		return GenEd25519, nil
	case crypto.NameSecp256k1:
		return GenSecp256k1, nil
	default:
		return nil, errors.Errorf("Cannot generate keys for algorithm: %s", algo)
	}
}
