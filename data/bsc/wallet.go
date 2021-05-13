package bsc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/ubiq/go-ubiq/core/types"
	"math/big"
)

var (
	ErrNoKey = errors.New("wallet doesn't have requested key")
)

type Wallet struct {
	hd    	bool
	master  *Deriver
	keys  	map[common.Address]ecdsa.PrivateKey
}

func NewHDWallet(hdPrivate string, n uint64) (*Wallet, error) {
	master, err := NewDeriver(hdPrivate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init key deriver")
	}

	wallet := &Wallet{
		hd: true,
		master: master,
		keys: make(map[common.Address]ecdsa.PrivateKey),
	}

	if err := wallet.extend(n); err != nil {
		return nil, errors.Wrap(err, "failed to extend master")
	}

	return wallet, nil
}

func (wallet *Wallet) Import(raw []byte) (common.Address, error) {
	priv, err := crypto.ToECDSA(raw)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to convert private key")
	}
	address := crypto.PubkeyToAddress(priv.PublicKey)
	wallet.keys[address] = *priv
	return address, nil
}

func (wallet *Wallet) ImportHex(data string) (common.Address, error) {
	raw, err := hex.DecodeString(data)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "failed to decode string")
	}
	return wallet.Import(raw)
}

func (wallet *Wallet) extend(i uint64) error {
	for uint64(len(wallet.keys)) < i {
		child, err := wallet.master.ChildPrivate(uint32(len(wallet.keys)))
		if err != nil {
			return errors.Wrap(err, "failed to extend child")
		}

		raw, err := hex.DecodeString(child)
		if err != nil {
			return errors.Wrap(err, "failed decode private key")
		}

		if _, err := wallet.Import(raw); err != nil {
			return errors.Wrap(err, "failed to import key")
		}
	}

	return nil
}

func (wallet *Wallet) Addresses(ctx context.Context) (result []common.Address) {
	for addr := range wallet.keys {
		result = append(result, addr)
	}
	return result
}

func (wallet *Wallet) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int,
	) (*types.Transaction, error) {
	key, ok := wallet.keys[address]
	if !ok {
		return nil, ErrNoKey
	}
	return SignTxWithPrivate(&key, tx, chainID)
}

func SignTxWithPrivate(key *ecdsa.PrivateKey, tx *types.Transaction, chainID *big.Int,
	) (*types.Transaction, error) {
	return types.SignTx(tx, types.NewEIP155Signer(chainID), key)
}

func (wallet *Wallet) HasAddress(address common.Address) bool {
	_, ok := wallet.keys[address]
	return ok
}
