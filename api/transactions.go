package api

import (
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/minter-go-node/core/transaction"
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/rpc/core/types"
	"math/big"
)

type TransactionResponse struct {
	Hash     common.HexBytes   `json:"hash"`
	RawTx    string            `json:"raw_tx"`
	Height   int64             `json:"height"`
	Index    uint32            `json:"index"`
	From     string            `json:"from"`
	Nonce    uint64            `json:"nonce"`
	GasPrice *big.Int          `json:"gas_price"`
	GasCoin  types.CoinSymbol  `json:"gas_coin"`
	GasUsed  int64             `json:"gas_used"`
	Type     byte              `json:"type"`
	Data     json.RawMessage   `json:"data"`
	Payload  []byte            `json:"payload"`
	Tags     map[string]string `json:"tags"`
	Code     uint32            `json:"code,omitempty"`
	Log      string            `json:"log,omitempty"`
}

type ResultTxSearch struct {
	Txs        []*core_types.ResultTx `json:"txs"`
	TotalCount int                    `json:"total_count"`
}

func Transactions(query string) (*[]TransactionResponse, error) {
	rpcResult, err := client.TxSearch(query, false, 1, 100)
	if err != nil {
		return nil, err
	}

	result := make([]TransactionResponse, len(rpcResult.Txs))
	for i, tx := range rpcResult.Txs {
		decodedTx, _ := transaction.DecodeFromBytes(tx.Tx)
		sender, _ := decodedTx.Sender()

		tags := make(map[string]string)

		for _, tag := range tx.TxResult.Tags {
			switch string(tag.Key) {
			case "tx.type":
				tags[string(tag.Key)] = fmt.Sprintf("%X", tag.Value)
			default:
				tags[string(tag.Key)] = string(tag.Value)
			}
		}

		data, err := encodeTxData(decodedTx)
		if err != nil {
			return nil, err
		}

		result[i] = TransactionResponse{
			Hash:     common.HexBytes(tx.Tx.Hash()),
			RawTx:    fmt.Sprintf("%x", []byte(tx.Tx)),
			Height:   tx.Height,
			Index:    tx.Index,
			From:     sender.String(),
			Nonce:    decodedTx.Nonce,
			GasPrice: decodedTx.GasPrice,
			GasCoin:  decodedTx.GasCoin,
			GasUsed:  tx.TxResult.GasUsed,
			Type:     decodedTx.Type,
			Data:     data,
			Payload:  decodedTx.Payload,
			Tags:     tags,
			Code:     tx.TxResult.Code,
			Log:      tx.TxResult.Log,
		}
	}

	return &result, nil
}
