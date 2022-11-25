package query

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
)

type JWTRequest struct {
	Address string `json:"address"`
}

type JWT struct {
	Token string `json:"token"`           // If the request succeeded, this is your actual JWT
	Error string `json:"error,omitempty"` // server will respond with error information (if the request failed) in this JSON field
}

// The swap that we are simulating. These correspond to the options the user selects with the CLI tool
// Basically, how much the user wants to trade (e.g. 100 JUNO for X OSMO)
// swagger:parameters swapRequest
type SimulatedSwapRequest struct {
	TokenInSymbol     string
	TokenInAmount     string
	TokenOutSymbol    string
	TokenOutMinAmount string
	UserWallet        string
	ArbitrageWallet   string // By default this is just the user's own wallet. See README for detailed explanation.
	// If you ask to simulate a trade of 100 JUNO for X OSMO, but your wallet doesn't have at least 100 JUNO,
	// it will cause an error. But if you specify "--skip-checks" it will simulate for you anyway instead of causing an error.
	SkipWalletFundsCheck bool
}

// The swap that we are simulating, with the exact pools to trade through
type SimulatedSwapExactPoolsRequest struct {
	TokenInDenom         string
	TokenInAmount        string
	Routes               []SwapAmountInRoute
	TokenOutMinAmount    string
	UserWallet           string
	ArbitrageWallet      string
	SkipWalletFundsCheck bool
}

type ExecuteMsg struct {
	Swap *Swap `json:"swap,omitempty"`
}

type JunoswapInstantiateContract struct {
	Token1Denom          *Token `json:"token1_denom"`
	Token2Denom          *Token `json:"token2_denom"`
	LpTokenCodeID        int    `json:"lp_token_code_id"`
	ProtocolFeeRecipient string `json:"protocol_fee_recipient"`
	ProtocolFeePercent   string `json:"protocol_fee_percent"`
	LpFeePercent         string `json:"lp_fee_percent"`
}

// Only one of "Native,CW20" will be set.
type Token struct {
	Native   string `json:"native,omitempty"`
	CW20     string `json:"cw20,omitempty"`
	Decimals int
}

//{ "token1_denom": {"native": "ujunox" }, "token2_denom": {"native": "uusdcx" }, "lp_token_code_id": 2658,
//"protocol_fee_recipient": "juno15eft7zal7mg3e8enpg94el6lwnlhsxrcu8025v", "protocol_fee_percent": "0.2", "lp_fee_percent": "0.1"}

type Swap struct {
	InputToken  string `json:"input_token"` //Token1 or Token2
	InputAmount string `json:"input_amount"`
	MinOutput   string `json:"min_output"`
}

type SwapAmountInRoute struct {
	Pool          string `json:"pool_id,omitempty" yaml:"pool_id"`
	TokenOutDenom string `json:"token_out_denom,omitempty" yaml:"token_out_denom"`
	//When you swap on Junoswap, the smart contracts require that you specify either "Token1" or "Token2"
	//as the input token (instead of the token's denomination).
	TokenOutJunoswapID string `json:"token_out_junoswap_id,omitempty" yaml:"token_out_junoswap_id"`
	TokenOutAmount     string `json:"token_out_amount,omitempty"` //the amount out (estimated) for the results of the swap for this particular route
}

// Results of the simulation
// swagger:model result
type SimulatedSwapResult struct {
	SimulatedUserSwap       *SimulatedSwap `json:"userSwap,omitempty"`      // the user's swap including the most efficient routes (pools) to use
	ArbitrageSwap           *ArbitrageSwap `json:"arbitrageSwap,omitempty"` // how much arbitrage the user's swap will cause, routes to use, etc
	HasArbitrageOpportunity bool           // whether or not the user's swap would cause arbitrage
	Error                   string         `json:"error,omitempty"` // if there was some issue detected on the server
}

type ArbitrageSwap struct {
	SimulatedSwap                *SimulatedSwap // the arbitrage swap including the most efficient routes (pools) to use
	EstimatedProfitHumanReadable string         // e.g. 11.64 OSMO
	EstimatedProfitBaseAmount    string         // e.g. 11.64
}

// FooBarResponse represents body of FooBar response.
type FooBarResponse struct {
	Baz struct {
		Prop string `json:"prop"`
	} `json:"baz"`
}

type SimulatedSwap struct {
	TokenIn                cosmosTypes.Coin
	TokenInJunoswapID      string `json:"token_in_junoswap_id,omitempty"`
	TokenOutMinAmount      string `json:"token_out_min_amount,omitempty"`
	Pools                  string
	Routes                 []SwapAmountInRoute `json:"routes,omitempty"`
	TokenOutAmount         cosmosTypes.Int
	TokenOutDenom          string // one of the 'denom' from asset lists at https:// github.com/osmosis-labs/assetlists/tree/main/osmosis-1
	TokenInSymbol          string
	AmountOutHumanReadable string // e.g. 11.64 OSMO
	TokenOutSymbol         string // e.g. OSMO
	BaseAmount             string // e.g. 11.64
	PriceImpact            float64
}
