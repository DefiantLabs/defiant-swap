package query

import (
	"github.com/cosmos/cosmos-sdk/client"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
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
