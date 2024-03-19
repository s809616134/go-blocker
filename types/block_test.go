package types

import (
	"testing"

	"github.com/s809616134/go-blocker/crypto"
	"github.com/s809616134/go-blocker/util"
	"github.com/stretchr/testify/assert"
)

func TestSignVerifyBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GeneratPrivateKey()
		pubKey  = privKey.Public()
		sig     = SignBlock(privKey, block)
	)

	// 64 btyes long signature
	assert.Equal(t, 64, len(sig.Bytes()))
	// the pubKey should match the hashed block
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

	assert.Equal(t, block.PublicKey, pubKey.Bytes())
	assert.Equal(t, block.Signature, sig.Bytes())

	assert.True(t, VerifyBlock(block))

	invalidPrivKey := crypto.GeneratPrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))
}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, 32, len(hash))
}
