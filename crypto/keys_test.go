package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratPrivateKey()
	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)
	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), PubKeyLen)
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "0d106558cf7ef735ee8e16a116940f66936b0b06a28d8c364694fbe851cec503"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "09826edcd3dbb55c91e8e318cb98f5e1dbf1f166"
	)
	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)
	address := privKey.Public().Address()
	fmt.Println(address)
	assert.Equal(t, address.String(), addressStr)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratPrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))

	// Test wit invalid message
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	// Test wit invalid pubKey
	invalidPrivKey := GeneratPrivateKey()
	invalidPubKey := invalidPrivKey.Public()
	assert.False(t, sig.Verify(invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratPrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, AddressLen, len(address.Bytes()))
}
