package cmd

import (
	"os"

	params "github.com/CosmosContracts/juno/v11/app/params"
	"github.com/DefiantLabs/JunoswapArbitrageCLI/query"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := query.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(os.ExpandEnv("$HOME/") + ".juno").
		WithViper("JUNO")

	rootCmd := &cobra.Command{
		Use:   "defiant-swap",
		Short: "A CLI tool for Juno swaps that captures arbitrage to maximize swap results",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
	}

	initRootCmd(rootCmd, encodingConfig)
	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	rootCmd.AddCommand(swapCmd)
	rootCmd.AddCommand(ledgerCmd)
	rootCmd.AddCommand(contractCmd)
}
