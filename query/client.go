package query

import (
	"encoding/json"
	"fmt"

	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
)

func BuildSwapExactAmountIn(
	clientCtx client.Context,
	amountTokenIn string,
	inDenom string,
	token1Or2 string,
	routes []SwapAmountInRoute,
	userAddr string,
	logNameSwapType string,
) []cosmosSdk.Msg {
	msgs := []cosmosSdk.Msg{}

	fmt.Printf("Performing %s. Swap has %d routes\n", logNameSwapType, len(routes))

	for i, route := range routes {
		msg := ExecuteMsg{
			Swap: &Swap{
				InputToken:  token1Or2,
				InputAmount: amountTokenIn,
				MinOutput:   route.TokenOutAmount,
			},
		}
		b, _ := json.Marshal(msg)

		fmt.Printf("%s, route %d. Token in (denom): %s. Token in (junoswap pool ID): %s. Amount in: %s. Minimum amount out: %s. Pool: %s.\n",
			logNameSwapType, i, inDenom, token1Or2, amountTokenIn, route.TokenOutAmount, route.Pool)

		amt, _ := cosmosSdk.NewIntFromString(amountTokenIn)
		req := &wasm.MsgExecuteContract{
			Sender:   userAddr,
			Contract: route.Pool,
			Msg:      b,
			Funds: []cosmosSdk.Coin{
				cosmosSdk.NewCoin("ujuno", amt),
			},
		}

		msgs = append(msgs, req)

		//there are two tokens per junoswap pool. therefore the token in for the next route is the opposite of the token out for the next route.
		if i < len(routes)-1 {
			if routes[i+1].TokenOutJunoswapID == "Token1" {
				token1Or2 = "Token2"
			} else {
				token1Or2 = "Token1"
			}
		}

		inDenom = route.TokenOutDenom
		amountTokenIn = route.TokenOutAmount
	}

	return msgs
}
