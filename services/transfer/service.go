package transfer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ubiq/go-ubiq/core/types"
	"math/big"
)

type BNBTransfer struct {
	From	common.Address
	To 		common.Address
	Value	*big.Int
	Raw 	types.Log
}
