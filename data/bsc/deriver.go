package bsc

import (
	"encoding/hex"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"math"
)

type Deriver struct {
	key *hdkeychain.ExtendedKey
}

func NewDeriver(src string) (*Deriver, error) {
	key, err := hdkeychain.NewKeyFromString(src)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse key")
	}
	return &Deriver{key}, nil
}

func (d *Deriver) ChildPrivate(i uint32) (string, error) {
	child, err := d.key.Derive(i)
	if err != nil {
		return "", err
	}

	priv, err := child.ECPrivKey()
	if err != nil {
		return "", errors.Wrap(err, "failed to get private key")
	}

	return hex.EncodeToString(priv.Serialize()), nil
}

func (d *Deriver) ChildAddress(i uint64) (string, error) {
	if i >= math.MaxUint32 {
		panic("child overflow")
	}
	child, err := d.key.Derive(uint32(i))
	if err != nil {
		return  "", err
	}

	public, err := child.ECPubKey()
	if err != nil {
		return "", errors.Wrap(err, "failed to get public key")
	}
	address := crypto.PubkeyToAddress(*public.ToECDSA())

	return address.Hex(), nil
}
