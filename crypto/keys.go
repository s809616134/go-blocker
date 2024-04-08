package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	PrivKeyLen   = 64
	SignatureLen = 64
	PubKeyLen    = 32
	SeedLen      = 32
	AddressLen   = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func NewPrivateKeyFromSeedStr(s string) *PrivateKey {
	seedBytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(seedBytes)
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != SeedLen {
		panic("invalid seed length, must be 32")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func GeneratPrivateKey() *PrivateKey {
	seed := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, PubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func PublicKeyFromBytes(b []byte) *PublicKey {
	if len(b) != PubKeyLen {
		panic("length of the publicKey bytes not equal to 64")
	}
	return &PublicKey{
		key: ed25519.PublicKey(b),
	}
}

func (p *PublicKey) Address() Address {
	return Address{
		value: p.key[len(p.key)-AddressLen:],
	}
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func SignatureFromBytes(b []byte) *Signature {
	if len(b) != SignatureLen {
		panic("length of the bytes not equal to 64")
	}

	return &Signature{
		value: b,
	}
}

func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	ok := ed25519.Verify(pubKey.key, msg, s.value)
	if !ok {
		fmt.Println("invalid signature of ed25519 Verify()")
	}
	return ok
}

type Address struct {
	value []byte
}

func AddressFromBytes(b []byte) Address {
	if len(b) != AddressLen {
		panic("length of the (address) bytes not equal to 20")
	}
	return Address{
		value: b,
	}
}

func (a Address) Bytes() []byte {
	return a.value
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}
