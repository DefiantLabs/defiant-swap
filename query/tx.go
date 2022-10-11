package query

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

//Gas fee is based on Osmosis fees for different numbers of pools
//Exact fee (in uosmo) == .005 * the number below, so it is quite minimal.
func GetGasFee(numRoutes int) uint64 {
	switch numRoutes {
	case 1:
		return 135000
	case 2:
		return 275000
	case 3:
		return 275000
	case 4:
		return 410000
	case 5:
		return 410000
	case 6:
		return 492000
	case 7:
		return 574000
	}

	return 0
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

func BuildTxFactory(clientContext *client.Context, gas uint64) tx.Factory {
	gasPrices := "0.005uosmo"
	txf := newFactoryCLI(*clientContext, gasPrices, gas)
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
