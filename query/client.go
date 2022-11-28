package query

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	osmosisParams "github.com/osmosis-labs/osmosis/v12/app/params"
	"github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

func BuildSwapExactAmountIn(clientCtx client.Context, tokenIn cosmosSdk.Coin, tokenOutMinAmt cosmosSdk.Int, routes []types.SwapAmountInRoute) (cosmosSdk.Msg, error) {

	msg := &types.MsgSwapExactAmountIn{
		Sender:            clientCtx.GetFromAddress().String(),
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmt,
	}

	return msg, nil
}

// chain := "osmosis-1"
// node := "https://rpc.osmosis.zone:443"
// osmosisHomeDir := "/home/kyle/.osmosisd"
//
//	keyringBackend := "test"
func GetOsmosisTxClient(encodingConfig osmosisParams.EncodingConfig, chain string, node string, osmosisHomeDir string, keyringBackend string, fromFlag string) *client.Context {
	//encodingConfig := osmosis.MakeEncodingConfig()
	clientCtx := client.Context{
		ChainID:      chain,
		NodeURI:      node,
		KeyringDir:   osmosisHomeDir,
		GenerateOnly: false,
	}

	ctxKeyring, krErr := client.NewKeyringFromBackend(clientCtx, keyringBackend)
	if krErr != nil {
		return nil
	}

	clientCtx = clientCtx.WithKeyring(ctxKeyring)

	//Where node is the node RPC URI
	rpcClient, rpcErr := client.NewClientFromNode(node)

	if rpcErr != nil {
		return nil
	}

	fromAddr, fromName, _, err := client.GetFromFields(clientCtx.Keyring, fromFlag, clientCtx.GenerateOnly)
	if err != nil {
		return nil
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
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastAsync).
		WithHomeDir(osmosisHomeDir).
		WithViper("OSMOSIS").
		WithNodeURI(node).
		WithClient(rpcClient).
		WithSkipConfirmation(true)

	return &clientCtx
}
