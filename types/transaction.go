package types

import (
	"crypto/sha256"

	pb "github.com/golang/protobuf/proto"
	"github.com/s809616134/go-blocker/crypto"
	"github.com/s809616134/go-blocker/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		if len(input.Signature) == 0 {
			panic("the transaction has no signature")
		}

		sig := crypto.SignatureFromBytes(input.Signature)
		pubKey := crypto.PublicKeyFromBytes(input.PublicKey)

		// we don't have signature in tx when we are signing
		// we should verify tx without signature

		// TODO: make sure we don't run into issue after verification
		// cause we set the signature to nil
		tempSig := input.Signature
		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
		input.Signature = tempSig
	}
	return true
}
