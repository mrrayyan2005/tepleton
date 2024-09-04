package cryptostore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tepleton/go-keys/cryptostore"
)

func TestNoopEncoder(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	noop := cryptostore.Noop

	key := cryptostore.GenEd25519.Generate()
	key2 := cryptostore.GenSecp256k1.Generate()

	b, err := noop.Encrypt(key, "encode")
	require.Nil(err)
	assert.NotEmpty(b)

	b2, err := noop.Encrypt(key2, "encode")
	require.Nil(err)
	assert.NotEmpty(b2)
	assert.NotEqual(b, b2)

	// note the decode with a different password works - not secure!
	pk, err := noop.Decrypt(b, "decode")
	require.Nil(err)
	require.NotNil(pk)
	assert.Equal(key, pk)

	pk2, err := noop.Decrypt(b2, "kggugougp")
	require.Nil(err)
	require.NotNil(pk2)
	assert.Equal(key2, pk2)
}

func TestSecretBox(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	enc := cryptostore.SecretBox

	key := cryptostore.GenEd25519.Generate()
	pass := "some-special-secret"

	b, err := enc.Encrypt(key, pass)
	require.Nil(err)
	assert.NotEmpty(b)

	// decoding with a different pass is an error
	pk, err := enc.Decrypt(b, "decode")
	require.NotNil(err)
	require.Nil(pk)

	// but decoding with the same passphrase gets us our key
	pk, err = enc.Decrypt(b, pass)
	require.Nil(err)
	assert.Equal(key, pk)
}
