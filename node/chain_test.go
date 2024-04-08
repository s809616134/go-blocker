package node

import (
	"testing"

	"github.com/s809616134/go-blocker/crypto"
	"github.com/s809616134/go-blocker/proto"
	"github.com/s809616134/go-blocker/types"
	"github.com/s809616134/go-blocker/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomBLock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratPrivateKey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	assert.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	for i := 0; i < 100; i++ {
		b := randomBLock(t, chain)
		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		block := randomBLock(t, chain)
		// hash the Block when all block info is complete
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestAddBlockWithTxInsufficientFunds(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBLock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(godSeed)
		recipient = crypto.GeneratPrivateKey().Public().Address().Bytes()
	)

	// our address got 1000 from genesis block
	// this is the genesis block tx hash
	prevTx, err := chain.txStore.Get("35cbbf7481d21f7763864d7632f1c24f07d11ff34cedb6fad79d65b39b05b995")
	assert.Nil(t, err)

	// get the 1000 tokens from genesis output
	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		// send too much to recipient
		{
			Amount:  1001,
			Address: recipient,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.NotNil(t, chain.AddBlock(block))
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBLock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(godSeed)
		recipient = crypto.GeneratPrivateKey().Public().Address().Bytes()
	)

	// our address got 1000 from genesis block
	// this is the genesis block tx hash
	prevTx, err := chain.txStore.Get("35cbbf7481d21f7763864d7632f1c24f07d11ff34cedb6fad79d65b39b05b995")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		// send 100 to recipient
		{
			Amount:  100,
			Address: recipient,
		},
		// get the rest 900 back to our address
		{
			Amount:  900,
			Address: privKey.Public().Address().Bytes(),
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	require.Nil(t, chain.AddBlock(block))
}
