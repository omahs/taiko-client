package rpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestBlockTx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := os.Getenv("L1_NODE_WS_ENDPOINT")
	l1Client, err := NewEthClient(ctx, url, time.Second*20)
	assert.NoError(t, err)

	priv := os.Getenv("L1_PROPOSER_PRIVATE_KEY")
	sk, err := crypto.ToECDSA(common.FromHex(priv))
	assert.NoError(t, err)

	chainID, err := l1Client.ChainID(ctx)
	assert.NoError(t, err)

	opts, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	assert.NoError(t, err)
	opts.Context = ctx
	//opts.NoSend = true

	balance, err := l1Client.BalanceAt(ctx, opts.From, nil)
	assert.NoError(t, err)
	t.Logf("address: %s, balance: %s", opts.From.String(), balance.String())

	tx, err := l1Client.TransactBlobTx(opts, []byte("s"))
	assert.NoError(t, err)

	receipt, err := bind.WaitMined(ctx, l1Client, tx)
	assert.NoError(t, err)
	t.Log(receipt)
}