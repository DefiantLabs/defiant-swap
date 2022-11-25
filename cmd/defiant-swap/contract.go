package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/DefiantLabs/JunoswapArbitrageCLI/query"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Instantiates contract",
	Long:  `Helper func for internal use`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		flagSet := cmd.Flags()
		rpcProvider := "https://rpc.uni.juno.deuslabs.fi:443"

		//clientCtx = clientCtx.WithNodeURI(rpcProvider)
		//clientCtx = clientCtx.WithChainID("uni-5")

		chain_u, chainErr := flagSet.GetString("chain-id")
		if chainErr != nil && chain_u != "" {
			clientCtx = clientCtx.WithChainID(chain_u)
		}

		if err != nil {
			fmt.Println(err.Error())
			cobra.CheckErr(err)
		}
		//flagSet := cmd.Flags()

		node, nErr := flagSet.GetString("node")
		if nErr != nil && node != "" {
			rpcProvider = node
		}

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

		fmt.Printf("keyringDir: %s\n", clientCtx.KeyringDir)
		kri, err := clientCtx.Keyring.List()
		cobra.CheckErr(err)

		fmt.Printf("Listing keyring keys, %d keys total\n", len(kri))
		counter := 1
		for _, v := range kri {
			fmt.Printf("(%d/%d)\n", counter, len(kri))
			fmt.Printf("Name: %s\nType: %s\nAddress: %s\nPubkey: %s\n", v.GetName(), v.GetType(), v.GetAddress(), v.GetPubKey())
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
			Salt:   []byte("test instantiate, baby!"),
		}

		txGas := query.GetGasFee(1)
		gasPrices := 0.001
		total := txGas * uint64(gasPrices)
		query.SubmitTxAwaitResponse(clientCtx, []sdk.Msg{req}, txGas, "0.001ujuno", fmt.Sprintf("%d%s", total, "ujuno"))

		return err
	},
}

func init() {
	contractCmd.Flags().String("chain-id", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	flags.AddTxFlagsToCmd(contractCmd)

	keysCmd := keys.AddKeyCommand()
	keysCmd.Flags().AddFlagSet(keys.Commands(".").PersistentFlags())
	keysCmd.SetArgs([]string{
		fmt.Sprintf("--%s=false", flags.FlagUseLedger),
		fmt.Sprintf("--%s=%s", cli.OutputFlag, "text"),
		fmt.Sprintf("--%s=%s", flags.FlagKeyAlgorithm, string(hd.Secp256k1Type)),
		fmt.Sprintf("--%s=%d", "coin-type", sdk.CoinType),
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
	})
	contractCmd.Flags().AddFlagSet(keysCmd.Flags())
	//ledgerCmd.Flags().BoolVar(&delete, "delete", false, "Delete the given ledger key.")
}
