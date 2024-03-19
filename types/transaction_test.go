package types

import (
	"testing"

	"github.com/s809616134/go-blocker/crypto"
	"github.com/s809616134/go-blocker/proto"
	"github.com/s809616134/go-blocker/util"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratPrivateKey()
	fromAddress := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratPrivateKey()
	toAddress := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	// like 'balance' of the ACCOUNT model
	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			input,
		},
		Outputs: []*proto.TxOutput{
			output1, output2,
		},
	}

	// sign the tx
	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	assert.True(t, VerifyTransaction(tx))
}
