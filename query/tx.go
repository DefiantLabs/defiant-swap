package query

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/cobra"
)

func SubmitTxAwaitResponse(clientCtx client.Context, msgs []sdk.Msg, txGas uint64, gasPrices string) {
	txf := BuildTxFactory(clientCtx, txGas, gasPrices)
	txf, txfErr := PrepareFactory(clientCtx, clientCtx.GetFromName(), txf)
	cobra.CheckErr(txfErr)

	txBuilder, err := tx.BuildUnsignedTx(txf, msgs...)
	cobra.CheckErr(err)

	txBuilder.SetFeeGranter(clientCtx.GetFeeGranterAddress())
	err = tx.Sign(txf, clientCtx.GetFromName(), txBuilder, true)
	cobra.CheckErr(err)

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	cobra.CheckErr(err)

	resp, err := clientCtx.BroadcastTxSync(txBytes)
	cobra.CheckErr(err)
	fmt.Printf("TX result code: %d. Waiting for TX with hash %s\n", resp.Code, resp.TxHash)

	tx1resp, err := AwaitTx(clientCtx, resp.TxHash, 5*time.Second)
	cobra.CheckErr(err)
	fmt.Printf("%+v\n", tx1resp)
}

// Get the TX by hash, waiting for it to be included in a block
func AwaitTx(clientCtx client.Context, txHash string, timeout time.Duration) (*txTypes.GetTxResponse, error) {
	var txByHash *txTypes.GetTxResponse
	var txLookupErr error
	startTime := time.Now()
	timeBetweenQueries := 100

	txClient := txTypes.NewServiceClient(clientCtx)

	for txByHash == nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		if time.Since(startTime) > timeout {
			cancel()
			return nil, txLookupErr
		}

		txByHash, txLookupErr = txClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
		if txLookupErr != nil {
			time.Sleep(time.Duration(timeBetweenQueries) * time.Millisecond)
		}
		cancel()
	}

	return txByHash, nil
}

func GetGasFee(numRoutes int) uint64 {
	return uint64(numRoutes * 200000)
}

var (
	// Variables used for retries
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 400)
	RtyErr    = retry.LastErrorOnly(true)
)

func GetKeyAddress(clientCtx client.Context, keyName string) (sdk.AccAddress, error) {
	info, err := clientCtx.Keyring.Key(keyName)
	if err != nil {
		return nil, err
	}
	return info.GetAddress(), nil
}

func PrepareFactory(clientCtx client.Context, keyName string, txf tx.Factory) (tx.Factory, error) {
	var (
		err      error
		from     sdk.AccAddress
		num, seq uint64
	)

	// Get key address and retry if fail
	if err = retry.Do(func() error {
		from, err = GetKeyAddress(clientCtx, keyName)
		if err != nil {
			return err
		}
		return err
	}, RtyAtt, RtyDel, RtyErr); err != nil {
		return tx.Factory{}, err
	}

	// Set the account number and sequence on the transaction factory and retry if fail
	if err = retry.Do(func() error {
		if err = txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
			return err
		}
		return err
	}, RtyAtt, RtyDel, RtyErr); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()

	if initNum == 0 || initSeq == 0 {
		if err = retry.Do(func() error {
			num, seq, err = txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
			if err != nil {
				return err
			}
			return err
		}, RtyAtt, RtyDel, RtyErr); err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}

func BuildTxFactory(clientContext client.Context, gas uint64, gasPrices string) tx.Factory {
	txf := newFactoryCLI(clientContext, gasPrices, gas)
	return txf
}

// NewFactoryCLI creates a new Factory.
func newFactoryCLI(clientCtx client.Context, gasPrices string, gas uint64) tx.Factory {
	f := tx.Factory{}

	f = f.WithChainID(clientCtx.ChainID)
	f = f.WithKeybase(clientCtx.Keyring)
	f = f.WithAccountRetriever(clientCtx.AccountRetriever)
	f = f.WithTxConfig(clientCtx.TxConfig)
	f = f.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	f = f.WithGas(gas)
	f = f.WithGasPrices(gasPrices)

	if clientCtx.SignModeStr == flags.SignModeLegacyAminoJSON {
		//fmt.Println("Default sign-mode 'direct' not supported by Ledger, using sign-mode 'amino-json'.")
		f = f.WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	}

	return f
}
