package bsc

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/ubiq/go-ubiq/core/types"
	"math/big"
	"strings"
)

type KeyPair struct {
	wallet *Wallet
}

func NewKeyPair(hex string) (*KeyPair, error) {
	hex = strings.TrimPrefix(hex, "0x")

	wallet := NewWallet()
	_, err := wallet.ImportHex(hex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import private key")
	}

	return &KeyPair{
		wallet: wallet,
	}, nil
}

func (kp *KeyPair) Address() common.Address {
	return kp.wallet.Addresses(context.Background())[0]
}

func (kp *KeyPair) SignTx(tx *types.Transaction, chainID *big.Int,
	) (*types.Transaction, error) {
	return kp.wallet.SignTx(kp.Address(), tx, chainIDx)
}