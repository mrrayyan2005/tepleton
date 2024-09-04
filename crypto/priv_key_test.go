package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	crypto "github.com/tepleton/go-crypto"
)

func TestGeneratePrivKey(t *testing.T) {
	testPriv := crypto.GenPrivKeyEd25519()
	testGenerate := testPriv.Generate(1)
	signBytes := []byte("something to sign")

	pub := testGenerate.PubKey()
	sig, err := testGenerate.Sign(signBytes)
	assert.NoError(t, err)
	assert.True(t, pub.VerifyBytes(signBytes, sig))
}

/*

type BadKey struct {
	PrivKeyEd25519
}

func TestReadPrivKey(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	// garbage in, garbage out
	garbage := []byte("hjgewugfbiewgofwgewr")
	XXX This test wants to register BadKey globally to go-crypto,
	but we don't want to support that.
	_, err := PrivKeyFromBytes(garbage)
	require.Error(err)

	edKey := GenPrivKeyEd25519()
	badKey := BadKey{edKey}

	cases := []struct {
		key   PrivKey
		valid bool
	}{
		{edKey, true},
		{badKey, false},
	}

	for i, tc := range cases {
		data := tc.key.Bytes()
		fmt.Println(">>>", data)
		key, err := PrivKeyFromBytes(data)
		fmt.Printf("!!! %#v\n", key, err)
		if tc.valid {
			assert.NoError(err, "%d", i)
			assert.Equal(tc.key, key, "%d", i)
		} else {
			assert.Error(err, "%d: %#v", i, key)
		}
	}
}
*/
