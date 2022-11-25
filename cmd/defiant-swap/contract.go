package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/DefiantLabs/JunoswapArbitrageCLI/query"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Instantiates contract",
	Long:  `Helper func for internal use`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			fmt.Println(err.Error())
			cobra.CheckErr(err)
		}
		//flagSet := cmd.Flags()
		rpcProvider := "https://rpc.uni.juno.deuslabs.fi:443"
		clientCtx = clientCtx.WithNodeURI(rpcProvider)
		clientCtx = clientCtx.WithChainID("uni-5")

		// RPC connection endpoint
		rpcClient, rpcErr := client.NewClientFromNode(rpcProvider)
		cobra.CheckErr(rpcErr)
		clientCtx = clientCtx.WithClient(rpcClient)

		from, _ := cmd.Flags().GetString(flags.FlagFrom)
		if from != "" {
			fmt.Printf("From: %s\n", from)
		}

		address := clientCtx.GetFromAddress().String()
		fmt.Printf("Address: %s\n", address)
		if address == "" {
			return errors.New("must specify user account address")
		}

		msg := query.JunoswapInstantiateContract{
			Token1Denom:          &query.Token{Native: "ujunox"},
			Token2Denom:          &query.Token{Native: "uusdcx"},
			LpTokenCodeID:        2658,
			ProtocolFeeRecipient: address,
			ProtocolFeePercent:   "0.2",
			LpFeePercent:         "0.1",
		}
		b, _ := json.Marshal(msg)

		req := &types.MsgInstantiateContract2{
			Sender: address,
			Admin:  address,
			CodeID: 2658,
			Label:  "woo!",
			Msg:    b,
		}

		txGas := query.GetGasFee(1)
		gasPrices := "0.001ujunox"
		query.SubmitTxAwaitResponse(clientCtx, []sdk.Msg{req}, txGas, gasPrices)

		return err
	},
}

func init() {
	flags.AddTxFlagsToCmd(contractCmd)
	//contractCmd.Flags().AddFlagSet(keys.Commands(".").PersistentFlags())
}
