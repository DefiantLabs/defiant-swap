package main

import (
	"os"

	cmd "github.com/DefiantLabs/JunoswapArbitrageCLI/cmd/defiant-swap"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, os.ExpandEnv("$HOME/")+".juno"); err != nil {
		os.Exit(1)
	}
}

//This test case requires that the given Sender address (swapMsg.Sender/'multipart') has already authorized a grant for 'arb' address
//On the CLI you can authorize the grant as follows:
//osmosisd tx authz grant {arb_address} generic --from multipart --keyring-backend test --chain-id osmosis-1 --msg-type "/osmosis.gamm.v1beta1.MsgSwapExactAmountIn" --gas=250000 --gas-prices=0.005uosmo --node https://rpc-osmosis.ecostake.com:443
// func SendTx(t *testing.T) {
// 	encodingConfig := osmosis.MakeEncodingConfig()
// 	chain := "osmosis-1"
// 	rpc := "https://rpc.osmosis.zone:443"
// 	osmosisHome := "/home/kyle/.osmosisd"
// 	tradeClient := "multipart"
// 	arbClient := "arb"

// 	tradeClientContext := query.GetOsmosisTxClient(encodingConfig, chain, rpc, osmosisHome, "test", tradeClient)
// 	route := []types.SwapAmountInRoute{{PoolId: 1, TokenOutDenom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"}}
// 	msg, msgErr := query.BuildSwapExactAmountIn(*tradeClientContext, "1024000uosmo", "1", route)
// 	if msgErr != nil {
// 		return
// 	}

// 	swapMsg := msg.(*types.MsgSwapExactAmountIn)
// 	swapMsg.Sender = "osmo14tkd4079rnk7vnt0q9pg3pj44eyz8ahqrtajln"
// 	arbClientContext := query.GetOsmosisTxClient(encodingConfig, chain, rpc, osmosisHome, "test", arbClient)

// 	if arbClientContext == nil {
// 		fmt.Println("Configuration error")
// 		os.Exit(1)
// 	}

// 	msgBytes, mErr := swapMsg.Marshal()
// 	if mErr != nil {
// 		return
// 	}

// 	msgExec := &authz.MsgExec{
// 		Grantee: "osmo14mmus5h7m6vkp0pteks8wawaj4wf3sx7fy3s2r",
// 		Msgs:    []*ctypes.Any{{TypeUrl: "/osmosis.gamm.v1beta1.MsgSwapExactAmountIn", Value: msgBytes}},
// 	}

// 	msg2, msg2Err := query.BuildSwapExactAmountIn(*arbClientContext, "2000000uosmo", "1", route)
// 	if msg2Err != nil {
// 		return
// 	}

// 	txf := query.PresignBuildTxFactory(arbClientContext, 3)
// 	txf, txfErr := query.PresignPrepareFactory(*arbClientContext, arbClient, txf)
// 	if txfErr != nil {
// 		return
// 	}

// 	txBuilder, err := tx.BuildUnsignedTx(txf, msgExec, msg2)
// 	if err != nil {
// 		return
// 	}

// 	txBuilder.SetFeeGranter(arbClientContext.GetFeeGranterAddress())
// 	err = tx.Sign(txf, arbClientContext.GetFromName(), txBuilder, true)
// 	if err != nil {
// 		return
// 	}

// 	txBytes, err := arbClientContext.TxConfig.TxEncoder()(txBuilder.GetTx())
// 	if err != nil {
// 		return
// 	}

// 	tx1resp, err := arbClientContext.BroadcastTxSync(txBytes)
// 	if err != nil {
// 		os.Exit(1)
// 	} else {
// 		util.Logger.Info("code", zap.Uint32("code", tx1resp.Code))
// 	}
// }
