package query

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/CosmWasm/wasmd/app"
	params "github.com/CosmosContracts/juno/v11/app/params"
	"github.com/avast/retry-go"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

var (
	Bech32PrefixAccAddr = "juno"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = Bech32PrefixAccAddr + "pub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = Bech32PrefixAccAddr + "valoper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = Bech32PrefixAccAddr + "valoperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = Bech32PrefixAccAddr + "valcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = Bech32PrefixAccAddr + "valconspub"
)

func InitAddressPrefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	app.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	app.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

func GetTxClient(chain string, node string, home string, keyringBackend string, fromFlag string) client.Context {
	encodingConfig := MakeEncodingConfig()
	clientCtx := client.Context{
		ChainID:      chain,
		NodeURI:      node,
		KeyringDir:   home,
		GenerateOnly: false,
	}

	ctxKeyring, krErr := client.NewKeyringFromBackend(clientCtx, keyringBackend)
	if krErr != nil {
		return client.Context{}
	}

	clientCtx = clientCtx.WithKeyring(ctxKeyring)

	//Where node is the node RPC URI
	rpcClient, rpcErr := client.NewClientFromNode(node)

	if rpcErr != nil {
		return client.Context{}
	}

	fromAddr, fromName, _, err := client.GetFromFields(clientCtx, clientCtx.Keyring, fromFlag)
	if err != nil {
		return client.Context{}
	}

	clientCtx = clientCtx.WithCodec(encodingConfig.Marshaler).
		WithChainID(chain).
		WithFrom(fromFlag).
		WithFromAddress(fromAddr).
		WithFromName(fromName).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastAsync).
		WithHomeDir(home).
		WithViper("JUNO").
		WithNodeURI(node).
		WithKeyringDir(home).
		WithClient(rpcClient).
		WithSkipConfirmation(true)

	return clientCtx
}

func SubmitTxAwaitResponse(clientCtx client.Context, msgs []sdk.Msg, gas uint64, gasPrices string, fees string) {
	txf := BuildTxFactory(clientCtx, gas, gasPrices, fees)
	txf, txfErr := PrepareFactory(clientCtx, clientCtx.GetFromName(), txf)
	cobra.CheckErr(txfErr)

	txBuilder, err := tx.BuildUnsignedTx(txf, msgs...)
	cobra.CheckErr(err)

	txBuilder.SetFeeGranter(clientCtx.GetFeeGranterAddress())

	fmt.Printf("From name (tx.sign): %s", clientCtx.GetFromName())
	err = tx.Sign(txf, clientCtx.GetFromName(), txBuilder, true)
	cobra.CheckErr(err)

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	cobra.CheckErr(err)

	resp, err := clientCtx.BroadcastTxSync(txBytes)
	cobra.CheckErr(err)
	fmt.Printf("TX result code: %d. Waiting for TX with hash %s\n", resp.Code, resp.TxHash)

	tx1resp, err := AwaitTx(clientCtx, resp.TxHash, 15*time.Second)
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
		} else {
			fmt.Printf("Key address: %s", from.String())
		}
		return err
	}, RtyAtt, RtyDel, RtyErr); err != nil {
		return tx.Factory{}, err
	}

	// Set the account number and sequence on the transaction factory and retry if fail
	if err = retry.Do(func() error {
		if err = txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
			return err
		} else {
			a := txf.AccountRetriever()
			fmt.Printf("Account retriever exists? %t\n", a != nil)
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
			fmt.Printf("init num: %d\n", num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
			fmt.Printf("sequence %d\n", seq)
		}
	}

	return txf, nil
}

func BuildTxFactory(clientContext client.Context, gas uint64, gasPrices string, fees string) tx.Factory {
	txf := newFactoryCLI(clientContext, gas, gasPrices, fees)
	return txf
}

// NewFactoryCLI creates a new Factory.
func newFactoryCLI(clientCtx client.Context, gas uint64, gasPrices string, fees string) tx.Factory {
	f := tx.Factory{}

	fmt.Printf("Chain ID: %s, keyring: %+v, account retriever: %+v, txconfig: %+v\n", clientCtx.ChainID, clientCtx.Keyring, clientCtx.AccountRetriever, clientCtx.TxConfig)

	f = f.WithChainID(clientCtx.ChainID)
	f = f.WithKeybase(clientCtx.Keyring)
	f = f.WithAccountRetriever(clientCtx.AccountRetriever)
	f = f.WithTxConfig(clientCtx.TxConfig)
	f = f.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	f = f.WithGas(gas)
	//f = f.WithGasPrices(gasPrices)
	f = f.WithFees(fees)

	if clientCtx.SignModeStr == flags.SignModeLegacyAminoJSON {
		//fmt.Println("Default sign-mode 'direct' not supported by Ledger, using sign-mode 'amino-json'.")
		f = f.WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	}

	return f
}
