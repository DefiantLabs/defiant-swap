package query

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

type JWTRequest struct {
	Address string `json:"address"`
}

type JWT struct {
	Token string `json:"token"`           // If the request succeeded, this is your actual JWT
	Error string `json:"error,omitempty"` // server will respond with error information (if the request failed) in this JSON field
}

// The swap that we are simulating. These correspond to the options the user selects with the CLI tool.
// Basically, how much the user wants to trade (e.g. 100 JUNO for X OSMO).
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

// Results of the simulation.
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

type SimulatedSwap struct {
	TokenIn                cosmosTypes.Coin
	TokenOutMinAmount      cosmosTypes.Int
	Pools                  string
	Routes                 types.SwapAmountInRoutes `json:"routes,omitempty"`
	TokenOutAmount         cosmosTypes.Int
	TokenOutDenom          string // one of the 'denom' from asset lists at https:// github.com/osmosis-labs/assetlists/tree/main/osmosis-1
	TokenInSymbol          string
	AmountOutHumanReadable string // e.g. 11.64 OSMO
	TokenOutSymbol         string // e.g. OSMO
	BaseAmount             string // e.g. 11.64
	PriceImpact            float64
}
