package wallet

import (
	"encoding/json"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/stretchr/testify/assert"
)

func TestMarshalWallet(t *testing.T) {
	w := NewWallet("jeremy")
	w.KeyRing = append(w.KeyRing, NewKeypair(crypto.NewEd25519(), []byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}))
	expected := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"vega/ed25519","tainted":false,"meta":null}]}`
	m, err := json.Marshal(&w)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(m))
}

func TestUnMarshalWallet(t *testing.T) {
	w := Wallet{}
	marshalled := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"vega/ed25519","tainted":false,"meta":null}]}`
	err := json.Unmarshal([]byte(marshalled), &w)
	assert.NoError(t, err)
	assert.Len(t, w.KeyRing, 1)
	assert.Equal(t, []byte{1, 2, 3, 4}, w.KeyRing[0].pubBytes)
	assert.Equal(t, []byte{4, 3, 2, 1}, w.KeyRing[0].privBytes)
	assert.Equal(t, "vega/ed25519", w.KeyRing[0].Algorithm.Name())
}

func TestUnMarshalWalletErrorInvalidAlgorithm(t *testing.T) {
	w := Wallet{}
	marshalled := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"notanalgorithm","tainted":false,"meta":null}]}`
	err := json.Unmarshal([]byte(marshalled), &w)
	assert.EqualError(t, err, crypto.ErrUnsupportedSignatureAlgorithm.Error())
}
