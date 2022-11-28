package cmd

import (
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"
)

type bechKeyOutFn func(keyInfo cryptokeyring.Info) (cryptokeyring.KeyOutput, error)

func printCreate(cmd *cobra.Command, info keyring.Info, showMnemonic bool, mnemonic string, outputFormat string) error {
	switch outputFormat {
	case OutputFormatText:
		cmd.PrintErrln()
		printKeyInfo(cmd.OutOrStdout(), info, keyring.MkAccKeyOutput, outputFormat)

		// print mnemonic unless requested not to.
		if showMnemonic {
			fmt.Fprintln(cmd.ErrOrStderr(), "\n**Important** write this mnemonic phrase in a safe place.")
			fmt.Fprintln(cmd.ErrOrStderr(), "It is the only way to recover your account if you ever forget your password.")
			fmt.Fprintln(cmd.ErrOrStderr(), "")
			fmt.Fprintln(cmd.ErrOrStderr(), mnemonic)
		}
	case OutputFormatJSON:
		out, err := keyring.MkAccKeyOutput(info)
		if err != nil {
			return err
		}

		if showMnemonic {
			out.Mnemonic = mnemonic
		}

		jsonString, err := keys.KeysCdc.MarshalJSON(out)
		if err != nil {
			return err
		}

		cmd.Println(string(jsonString))

	default:
		return fmt.Errorf("invalid output format %s", outputFormat)
	}

	return nil
}

func printKeyInfo(w io.Writer, keyInfo cryptokeyring.Info, bechKeyOut bechKeyOutFn, output string) {
	ko, err := bechKeyOut(keyInfo)
	if err != nil {
		panic(err)
	}

	switch output {
	case OutputFormatText:
		printTextInfos(w, []cryptokeyring.KeyOutput{ko})

	case OutputFormatJSON:
		out, err := keys.KeysCdc.MarshalJSON(ko)
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(w, string(out))
	}
}

func printInfos(w io.Writer, infos []cryptokeyring.Info, output string) {
	kos, err := cryptokeyring.MkAccKeysOutput(infos)
	if err != nil {
		panic(err)
	}

	switch output {
	case OutputFormatText:
		printTextInfos(w, kos)

	case OutputFormatJSON:
		out, err := keys.KeysCdc.MarshalJSON(kos)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "%s", out)
	}
}

func printTextInfos(w io.Writer, kos []cryptokeyring.KeyOutput) {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(out))
}
