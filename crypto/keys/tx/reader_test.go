package tx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tepleton/go-crypto"
	"github.com/tepleton/go-crypto/keys"
	"github.com/tepleton/go-crypto/keys/cryptostore"
	"github.com/tepleton/go-crypto/keys/storage/memstorage"
	data "github.com/tepleton/go-wire/data"
)

func TestReader(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	algo := crypto.NameEd25519
	cstore := cryptostore.New(
		cryptostore.SecretBox,
		memstorage.New(),
		keys.MustLoadCodec("english"),
	)
	type sigs struct{ name, pass string }
	u := sigs{"alice", "1234"}
	u2 := sigs{"bob", "foobar"}

	_, _, err := cstore.Create(u.name, u.pass, algo)
	require.Nil(err, "%+v", err)
	_, _, err = cstore.Create(u2.name, u2.pass, algo)
	require.Nil(err, "%+v", err)

	cases := []struct {
		tx   Sig
		sigs []sigs
	}{
		{New([]byte("first")), nil},
		{New([]byte("second")), []sigs{u}},
		{New([]byte("other")), []sigs{u2}},
		{NewMulti([]byte("m-first")), nil},
		{NewMulti([]byte("m-second")), []sigs{u}},
		{NewMulti([]byte("m-other")), []sigs{u, u2}},
	}

	for _, tc := range cases {
		tx := tc.tx

		// make sure json serialization and loading works w/o sigs
		var pre Sig
		pjs, err := data.ToJSON(tx)
		require.Nil(err, "%+v", err)
		err = data.FromJSON(pjs, &pre)
		require.Nil(err, "%+v", err)
		assert.Equal(tx, pre)

		for _, s := range tc.sigs {
			err = cstore.Sign(s.name, s.pass, tx)
			require.Nil(err, "%+v", err)
		}

		var post Sig
		sjs, err := data.ToJSON(tx)
		require.Nil(err, "%+v", err)
		err = data.FromJSON(sjs, &post)
		require.Nil(err, "%+v\n%s", err, string(sjs))
		assert.Equal(tx, post)

		if len(tc.sigs) > 0 {
			assert.NotEqual(pjs, sjs, "%s\n ------ %s", string(pjs), string(sjs))
		}
	}
}
