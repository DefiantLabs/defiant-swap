package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/DefiantLabs/JunoswapArbitrageCLI/query"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInstantiate(t *testing.T) {
	query.InitAddressPrefixes()

	clientCtx := query.GetTxClient("uni-5", "https://rpc.uni.juno.deuslabs.fi:443", "/home/kyle/.juno", "test", "testnetkey")
	address := clientCtx.FromAddress.String()

	msg := query.JunoswapInstantiateContract{
		Token1Denom:          &query.Token{Native: "ujunox"},
		Token2Denom:          &query.Token{Native: "uusdcx"},
		LpTokenCodeID:        2658,
		ProtocolFeeRecipient: address,
		ProtocolFeePercent:   "0.2",
		LpFeePercent:         "0.1",
	}
	b, _ := json.Marshal(msg)

	fmt.Printf("Msg JSON: %s\n", b)

	req := &types.MsgInstantiateContract2{
		Sender: address,
		Admin:  address,
		CodeID: 2658,
		Label:  "woo!",
		Msg:    b,
		Salt:   []byte("test instantiate, baby!"),
	}

	fees := "7500ujunox" //Note that gas*gasPrices MUST equal fees, or the TX will error
	query.SubmitTxAwaitResponse(clientCtx, []sdk.Msg{req}, 300000, "0.001ujunox", fees)
}
