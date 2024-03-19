package main

import (
	"context"
	"log"
	"time"

	"github.com/s809616134/go-blocker/crypto"
	"github.com/s809616134/go-blocker/node"
	"github.com/s809616134/go-blocker/proto"
	"github.com/s809616134/go-blocker/util"
	"google.golang.org/grpc"
)

func main() {
	makeNode(":3000", []string{}, true)
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"}, false)
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"}, false)

	for {
		time.Sleep(time.Millisecond * 100)
		makeTransaction()
	}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValdidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "Blocker-1",
		ListenAddr: listenAddr,
	}
	if isValdidator {
		cfg.PrivateKey = crypto.GeneratPrivateKey()
	}
	n := node.NewNode(cfg)
	go n.Start(listenAddr, bootstrapNodes)
	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)
	privKey := crypto.GeneratPrivateKey()

	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Address().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  90,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
