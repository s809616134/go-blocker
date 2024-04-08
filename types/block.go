package types

import (
	"bytes"
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
	pb "github.com/golang/protobuf/proto"
	"github.com/s809616134/go-blocker/crypto"
	"github.com/s809616134/go-blocker/proto"
	"github.com/sirupsen/logrus"
)

type TxHash struct {
	hash []byte
}

func NewTxHash(hash []byte) TxHash {
	return TxHash{hash: hash}
}

func (h TxHash) CalculateHash() ([]byte, error) {
	return h.hash, nil
}

func (h TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(h.hash, other.(TxHash).hash)
	return equals, nil
}

func VerifyBlock(b *proto.Block) bool {
	if len(b.Transactions) > 0 {
		if !VerifyRootHash(b) {
			logrus.Error("INVALID root hash")
			return false
		}
	}
	if len(b.PublicKey) != crypto.PubKeyLen {
		logrus.Error("INVALID public key length")
		return false
	}
	if len(b.Signature) != crypto.SignatureLen {
		logrus.Error("INVALID signature length")
		return false
	}

	var (
		sig    = crypto.SignatureFromBytes(b.Signature)
		pubKey = crypto.PublicKeyFromBytes(b.PublicKey)
		// Verify the block header
		hash = HashBlock(b)
	)
	return sig.Verify(pubKey, hash)
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	if len(b.Transactions) > 0 {
		tree, err := GetMerkleTree(b)
		if err != nil {
			panic(err)
		}
		b.Header.RootHash = tree.MerkleRoot()
	}

	hash := HashBlock(b)
	sig := pk.Sign(hash)
	b.PublicKey = pk.Public().Bytes()
	b.Signature = sig.Bytes()

	return sig
}

func VerifyRootHash(b *proto.Block) bool {
	tree, err := GetMerkleTree(b)
	if err != nil {
		return false
	}

	valid, err := tree.VerifyTree()
	if err != nil {
		return false
	}

	if !valid {
		return false
	}

	return bytes.Equal(b.Header.RootHash, tree.MerkleRoot())
}

func GetMerkleTree(b *proto.Block) (*merkletree.MerkleTree, error) {
	list := make([]merkletree.Content, len(b.Transactions))
	for i := 0; i < len(b.Transactions); i++ {
		list[i] = NewTxHash(HashTransaction(b.Transactions[i]))
	}

	// Create a new Merkle Tree from the list of Content
	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// HashBlock return a SHA256 of the header
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
