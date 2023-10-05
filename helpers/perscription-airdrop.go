package helpers

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/unlocker"
	"github.com/mrz1836/go-whatsonchain"
)

func PrescriptionAirdrop(perscOwnerAddress, inscTxId, inscWIF, fundingWIF string, c whatsonchain.ClientInterface) (string, string, error) {
	bsvTx := bt.NewTx()

	inscTxHex, err := c.GetRawTransactionData(context.Background(), inscTxId)
	if err != nil {
		return "", "", err
	}
	inscTx, err := bt.NewTxFromString(inscTxHex)
	if err != nil {
		return "", "", err
	}
	err = bsvTx.From(inscTxId,
		0,
		inscTx.Outputs[0].LockingScriptHexString(),
		inscTx.Outputs[0].Satoshis)
	if err != nil {
		return "", "", err
	}

	fundingWif, err := wif.DecodeWIF(fundingWIF)
	if err != nil {
		return "", "", err
	}
	// get public key bytes and address
	fundingPubkey := fundingWif.SerialisePubKey()
	fundingScript, err := bscript.NewP2PKHFromPubKeyBytes(fundingPubkey)
	if err != nil {
		return "", "", err
	}
	fundingAddr, err := bscript.NewAddressFromPublicKeyString(hex.EncodeToString(fundingPubkey), true)
	if err != nil {
		return "", "", err
	}
	utxos, err := c.AddressUnspentTransactions(context.Background(), fundingAddr.AddressString)
	if err != nil {
		return "", "", err
	}

	if len(utxos) < 1 {
		return "", "", errors.New("no UTXOs to use for funding")
	}

	utxoTxInfo, err := c.GetTxByHash(context.Background(), utxos[0].TxHash)
	if err != nil {
		return "", "", err
	}

	err = bsvTx.From(utxos[0].TxHash,
		0,
		utxoTxInfo.Vout[0].ScriptPubKey.Hex,
		uint64(utxos[0].Value))
	if err != nil {
		return "", "", err
	}

	err = bsvTx.PayToAddress(perscOwnerAddress, 1)
	if err != nil {
		return "", "", err
	}

	if utxos[0].Value < 5 {
		return "", "", errors.New("insufficient funds")
	}

	err = bsvTx.PayTo(fundingScript, uint64(utxos[0].Value-5))
	if err != nil {
		return "", "", err
	}

	// err = bsvTx.Change(fundingScript, bt.NewFeeQuote())
	// if err != nil {
	// 	return "", "", err
	// }

	inscWif, err := wif.DecodeWIF(inscWIF)
	if err != nil {
		return "", "", err
	}

	ug := &unlocker.Getter{PrivateKey: inscWif.PrivKey}
	u, err := ug.Unlocker(context.Background(), inscTx.Outputs[0].LockingScript)
	if err != nil {
		return "", "", err
	}
	err = bsvTx.FillInput(context.Background(), u, bt.UnlockerParams{InputIdx: 0})
	if err != nil {
		return "", "", err
	}

	ug2 := &unlocker.Getter{PrivateKey: fundingWif.PrivKey}
	u2, err := ug2.Unlocker(context.Background(), fundingScript) // reuse fundingScript
	if err != nil {
		return "", "", err
	}
	err = bsvTx.FillInput(context.Background(), u2, bt.UnlockerParams{InputIdx: 1})
	if err != nil {
		return "", "", err
	}

	fmt.Println(bsvTx.TxID())

	return bsvTx.String(), bsvTx.TxID(), nil
}
