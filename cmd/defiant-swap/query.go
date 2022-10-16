package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/DefiantLabs/OsmosisArbitrageBot/query"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

var genTokenUrl = "http://arb.defiantlabs.net:8080/api/token"
var simulateSwapUrl = "http://arb.defiantlabs.net:8080/api/secured/estimateswap"
var defiantRpc = "http://arb.defiantlabs.net:26657"
var defaultChain = "osmosis-1"

var ledgerCmd = &cobra.Command{
	Use:   "ledger <keyname>",
	Short: "Configure ledger device. Omit keyname to list ledger keys. Use --delete flag to delete the key.",
	Long:  `Will check that a ledger key exists and print all ledger keys, or prompt to add one`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyName := ""
		ctx, err := client.GetClientQueryContext(cmd)
		cobra.CheckErr(err)
		kb := ctx.Keyring

		if len(args) > 0 {
			keyName = args[0]
			if delete {
				errD := kb.Delete(keyName)
				cobra.CheckErr(errD)
				return nil
			} else {
				fmt.Printf("Adding key %s -- accept/reject request on Ledger device", keyName)
			}
		} else {
			kri, err := kb.List()
			cobra.CheckErr(err)
			fmt.Printf("Listing ledger keys, %d keys total\n", len(kri))
			counter := 1
			for _, v := range kri {
				fmt.Printf("(%d/%d)\n", counter, len(kri))
				fmt.Printf("Name: %s\nType: %s\nAddress: %s\nPubkey: %s\n", v.GetName(), v.GetType(), v.GetAddress(), v.GetPubKey())
			}
			return nil
		}

		_, err = kb.Key(keyName)
		buf := bufio.NewReader(ctx.Input)
		outputFormat := ctx.OutputFormat

		if err == nil {
			//  account exists, ask for user confirmation
			response, err2 := input.GetConfirmation(fmt.Sprintf("override the existing name %s", keyName), buf, cmd.ErrOrStderr())
			if err2 != nil {
				return err2
			}

			if !response {
				return errors.New("aborted")
			}

			err2 = kb.Delete(keyName)
			if err2 != nil {
				return err2
			}
		}

		coinType, _ := cmd.Flags().GetUint32("coin-type")
		account, _ := cmd.Flags().GetUint32("account")
		index, _ := cmd.Flags().GetUint32("index")
		hd.CreateHDPath(coinType, account, index)

		bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()
		keyringAlgos, _ := kb.SupportedAlgorithms()
		algoStr, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
		algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
		cobra.CheckErr(err)

		info, err := kb.SaveLedgerKey(keyName, algo, bech32PrefixAccAddr, coinType, account, index)
		if err != nil {
			return err
		}

		return printCreate(cmd, info, false, "", outputFormat)
	},
}

var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Performs a swap on Osmosis, optimizing rates for users",
	Long:  `Optimizes swaps by capturing arbitrage revenue that would normally go to bots. This is a free service provided by Defiant Labs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clientCtx, _ := client.GetClientTxContext(cmd)
		flagSet := cmd.Flags()

		clientCtx = clientCtx.WithNodeURI(defiantRpc)
		clientCtx = clientCtx.WithChainID(defiantRpc)
		rpcProvider := defiantRpc
		chain := defaultChain

		node, nErr := flagSet.GetString("node")
		if nErr != nil && node != "" {
			rpcProvider = node
		}

		chain_u, chainErr := flagSet.GetString("chain-id")
		if chainErr != nil && chain_u != "" {
			chain = chain_u
		}
		clientCtx = clientCtx.WithChainID(chain)

		// RPC connection endpoint
		rpcClient, rpcErr := client.NewClientFromNode(rpcProvider)
		cobra.CheckErr(rpcErr)
		clientCtx = clientCtx.WithClient(rpcClient)

		address := clientCtx.GetFromAddress().String()
		fmt.Printf("Address: %s\n", address)
		jwt := query.JWT{}
		jwtReq := query.JWTRequest{Address: address}

		queryParams := map[string]string{}
		if hasPartnerCode {
			secret, err := GetPwd("Enter partner secret:")
			if err != nil {
				return err
			}
			queryParams["partnerSecret"] = secret
		}

		// Make a request to the server to authenticate us. The server will return a JWT.
		httpStatus, err := query.PostJson(genTokenUrl, &jwtReq, &jwt, queryParams, nil)

		if err != nil {
			return err
		} else if httpStatus != 200 {
			if jwt.Error != "" {
				return fmt.Errorf("HTTP Status %d. Reason: %s", httpStatus, jwt.Error)
			} else {
				return fmt.Errorf("HTTP Status %d", httpStatus)
			}
		}

		cobra.CheckErr(err)
		symbolIn, _ := flagSet.GetString("in")
		symbolOut, _ := flagSet.GetString("out")
		amountIn, _ := flagSet.GetString("amount-in")
		amountOut, _ := flagSet.GetString("min-amount-out")
		arbitrageWallet, _ := flagSet.GetString("arb-wallet")

		simSwapReq := &query.SimulatedSwapRequest{
			TokenInSymbol:        symbolIn,
			TokenOutSymbol:       symbolOut,
			TokenInAmount:        amountIn,
			TokenOutMinAmount:    amountOut,
			SkipWalletFundsCheck: !verifyFunds,
			UserWallet:           address,
			ArbitrageWallet:      arbitrageWallet,
		}
		result := &query.SimulatedSwapResult{}
		httpStatus, err = query.PostJson(simulateSwapUrl, simSwapReq, &result, nil, &jwt)
		if result.Error != "" {
			return errors.New(result.Error)
		} else if httpStatus != 200 {
			return fmt.Errorf("issue with request, HTTP Status %d", httpStatus)
		}

		cobra.CheckErr(err)
		fmt.Printf("Your swap would trade %s %s for %s using Osmosis pool(s) %s.\n", amountIn, symbolIn,
			result.SimulatedUserSwap.AmountOutHumanReadable, result.SimulatedUserSwap.Pools)

		if result.HasArbitrageOpportunity {
			fmt.Printf("In addition, your swap would earn %s through arbitrage that would otherwise go to bots.\n",
				result.ArbitrageSwap.EstimatedProfitHumanReadable)
			fmt.Printf("Arbitrage trade (auto-generated by Defiant Labs): %s %s through pool(s) %s, expected result: %s out.\n",
				result.ArbitrageSwap.SimulatedSwap.TokenIn.Amount.String(),
				result.ArbitrageSwap.SimulatedSwap.TokenIn.Denom,
				result.ArbitrageSwap.SimulatedSwap.Pools,
				result.ArbitrageSwap.SimulatedSwap.TokenOutAmount.String(),
			)
		}

		err = Confirm("Proceed with the swap?")
		cobra.CheckErr(err)

		if err == nil {
			msgs := []types.Msg{}

			msgUserSwap, err := query.BuildSwapExactAmountIn(clientCtx, result.SimulatedUserSwap.TokenIn, result.SimulatedUserSwap.TokenOutMinAmount, result.SimulatedUserSwap.Routes)
			cobra.CheckErr(err)
			msgs = append(msgs, msgUserSwap)
			txGas := query.GetGasFee(len(result.SimulatedUserSwap.Routes))

			fmt.Printf("Performing user's swap. Token in: %s. Minimum amount out: %s. Pool(s) %s.\n",
				result.SimulatedUserSwap.TokenIn,
				result.SimulatedUserSwap.TokenOutMinAmount,
				result.SimulatedUserSwap.Pools)

			if result.HasArbitrageOpportunity {
				arbSwap, err := query.BuildSwapExactAmountIn(clientCtx, result.ArbitrageSwap.SimulatedSwap.TokenIn,
					result.ArbitrageSwap.SimulatedSwap.TokenOutMinAmount, result.ArbitrageSwap.SimulatedSwap.Routes)
				cobra.CheckErr(err)
				msgs = append(msgs, arbSwap)
				txGas = txGas + query.GetGasFee(len(result.ArbitrageSwap.SimulatedSwap.Routes))

				fmt.Printf("Performing arbitrage swap. Token in: %s. Minimum amount out: %s. Pool(s) %s.\n",
					result.ArbitrageSwap.SimulatedSwap.TokenIn,
					result.ArbitrageSwap.SimulatedSwap.TokenOutMinAmount,
					result.ArbitrageSwap.SimulatedSwap.Pools)
			}

			txf := query.BuildTxFactory(&clientCtx, txGas)
			txf, txfErr := query.PrepareFactory(clientCtx, clientCtx.GetFromName(), txf)
			cobra.CheckErr(txfErr)

			txBuilder, err := tx.BuildUnsignedTx(txf, msgs...)
			cobra.CheckErr(err)

			txBuilder.SetFeeGranter(clientCtx.GetFeeGranterAddress())
			err = tx.Sign(txf, clientCtx.GetFromName(), txBuilder, true)
			cobra.CheckErr(err)

			txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
			cobra.CheckErr(err)

			tx1resp, err := clientCtx.BroadcastTxSync(txBytes)
			cobra.CheckErr(err)
			fmt.Printf("TX result code: %d", tx1resp.Code)
		}

		return err
	},
}

var (
	arbitrageWallet string // wallet to use for arbs. defaults to user wallet.
	tokenFrom       string // token to trade from
	tokenTo         string // token to trade to
	amountIn        string // amount you want to trade
	amountOut       string // minimum amount you'll receive
	verifyFunds     bool
	hasPartnerCode  bool

	// for ledger
	delete bool
)

func init() {
	swapCmd.Flags().StringVar(&arbitrageWallet, "arb-wallet", "", "address to use for arbitrage")
	swapCmd.Flags().StringVar(&tokenFrom, "in", "", "The token you want to trade")
	swapCmd.Flags().StringVar(&tokenTo, "out", "", "The token you want to receive")
	swapCmd.Flags().StringVar(&amountIn, "amount-in", "", "The amount to trade (in the base amount). Ex: if the token is OSMO you might put --amount-in 101.5")
	swapCmd.Flags().StringVar(&amountOut, "min-amount-out", "", "The minimum amount of the token you want to receive, format is the same as amount-in")
	swapCmd.Flags().BoolVar(&verifyFunds, "verify-funds", true, "Check that the user's wallet contains enough funds for the trade. Turn off to simulate regardless of funds.")
	swapCmd.Flags().BoolVar(&hasPartnerCode, "partner", false, "Will prompt for partner secret if --partner=true. Unlocks unlimited API requests.")

	swapCmd.MarkFlagRequired("in")
	swapCmd.MarkFlagRequired("out")
	swapCmd.MarkFlagRequired("amount-in")
	swapCmd.MarkFlagRequired("min-amount-out")
	flags.AddTxFlagsToCmd(swapCmd)

	// Ledger setup
	keysCmd := keys.AddKeyCommand()
	keysCmd.Flags().AddFlagSet(keys.Commands(".").PersistentFlags())
	keysCmd.SetArgs([]string{
		fmt.Sprintf("--%s=true", flags.FlagUseLedger),
		fmt.Sprintf("--%s=%s", cli.OutputFlag, "text"),
		fmt.Sprintf("--%s=%s", flags.FlagKeyAlgorithm, string(hd.Secp256k1Type)),
		fmt.Sprintf("--%s=%d", "coin-type", sdk.CoinType),
		// fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
	})
	ledgerCmd.Flags().AddFlagSet(keysCmd.Flags())
	ledgerCmd.Flags().BoolVar(&delete, "delete", false, "Delete the given ledger key.")

}

func Confirm(prompt string) error {
	buf := bufio.NewReader(os.Stdin)
	ok, err := input.GetConfirmation(prompt, buf, os.Stderr)

	if err != nil || !ok {
		// _, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
		return errors.New("cancelled transaction")
	}

	return nil
}

func GetPwd(prompt string) (string, error) {
	buf := bufio.NewReader(os.Stdin)
	pass, err := input.GetPassword(prompt, buf)

	if err != nil {
		return "", errors.New("cancelled transaction")
	}

	return pass, nil
}
