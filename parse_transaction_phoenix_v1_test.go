package solanaparser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"os"
	"testing"
)

func TestTransaction_Phoenix_v1_ParseSwap(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("ESeR3usVhdFoj36RKSbt2rYHDJLruCAEJQaAhTqQhor5pHWGihGJytQVWx8BN14VqCnCnLvQAzUfUQQrZ2whGPp"),
		&rpc.GetParsedTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion1,
		})
	if err != nil {
		panic(err)
	}
	transaction := &rpc.ParsedTransactionWithMeta{
		Slot:        result.Slot,
		BlockTime:   result.BlockTime,
		Transaction: result.Transaction,
		Meta:        result.Meta,
	}
	tx := ParseTransaction(0, transaction)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_Phoenix_v1_ParseSwap_2(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("28mTGHKbC4NxD7AMZuW76A77vUx9XNnd4JfndAdL2eh3dJRC25QXh1Cvz9xGNwktRGicL3gW6kLV5GLT3tXPTU6j"),
		&rpc.GetParsedTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion1,
		})
	if err != nil {
		panic(err)
	}
	transaction := &rpc.ParsedTransactionWithMeta{
		Slot:        result.Slot,
		BlockTime:   result.BlockTime,
		Transaction: result.Transaction,
		Meta:        result.Meta,
	}
	tx := ParseTransaction(0, transaction)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}

func TestTransaction_Phoenix_v1_ParsePlaceLimitOrderWithFreeFunds(t *testing.T) {
	solClient := rpc.New(rpc.MainNetBeta_RPC)
	result, err := solClient.GetParsedTransaction(
		context.Background(),
		solana.MustSignatureFromBase58("3Z1XszYLjRxuKbaasJ58vHAJ6vTYonmMMDPspQqnWSK6SBm9eCvCMcUGiGaDJ2xAbzx5XxgrnXvvrX2qTuZTPc5g"),
		&rpc.GetParsedTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &rpc.MaxSupportedTransactionVersion1,
		})
	if err != nil {
		panic(err)
	}
	transaction := &rpc.ParsedTransactionWithMeta{
		Slot:        result.Slot,
		BlockTime:   result.BlockTime,
		Transaction: result.Transaction,
		Meta:        result.Meta,
	}
	tx := ParseTransaction(0, transaction)
	if tx == nil {
		panic("invalid transaction")
	}
	txJson, _ := json.MarshalIndent(tx, "", "    ")
	os.WriteFile(fmt.Sprintf("tx.json"), txJson, 0644)
}
