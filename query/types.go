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

// The swap that we are simulating. These correspond to the options the user selects with the CLI tool
// Basically, how much the user wants to trade (e.g. 100 JUNO for X OSMO)
// swagger:parameters swapRequest
type SimulatedSwapRequest struct {
	// e.g. OSMO or JUNO; must match the 'symbol' at https://github.com/osmosis-labs/assetlists/blob/main/osmosis-1/osmosis-1.assetlist.json.
	//
	// example: JUNO
	TokenInSymbol string
	// Amount corresponding to the given symbol (NOT the base denomination). This is what a user would type into the DEX frontend.
	//
	// example: 100
	TokenInAmount string
	// e.g. OSMO or JUNO; must match the 'symbol' at https://github.com/osmosis-labs/assetlists/blob/main/osmosis-1/osmosis-1.assetlist.json.
	//
	// example: OSMO
	TokenOutSymbol string
	// If trading 100 JUNO and you want at least 165 OSMO back, this would simply be set to "165".
	//
	// example: 165
	TokenOutMinAmount string
	// The swap simulation will be performed as if it is for this wallet.
	//
	// example: osmo1njsthsamgkzdqqg4awj7y6rxtuk3q26t2t8md2
	UserWallet string
	// On Osmosis, you must specify a wallet with at least 100 OSMO in it. If this param is left empty, the UserWallet will be used instead.
	// NOTE: At present, on Osmosis, only arbing with the OSMO token is supported. As we will support additional tokens in the future, we use
	// this parameter so that callers can indicate which arbitrage tokens they want to use.
	//
	// example: osmo14mmus5h7m6vkp0pteks8wawaj4wf3sx7fy3s2r
	ArbitrageWallet string
	// Setting SkipWalletFundsCheck to true will run the simulation regardless of whether the user's wallet has sufficient funds to perform the TX.
	// If this param is false, and the user wants to simulate a trade of 100 JUNO for X OSMO, the user's wallet must have at least 100 JUNO in it.
	//
	// example: true
	SkipWalletFundsCheck bool
}

// The swap that we are simulating, with the exact pools to trade through
type SimulatedSwapExactPoolsRequest struct {
	// e.g. uosmo or ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2; must match the 'denom' at https://github.com/osmosis-labs/assetlists/blob/main/osmosis-1/osmosis-1.assetlist.json.
	//
	// example: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2
	TokenInDenom string
	// Amount corresponding to the given denom (NOT the symbol). This is the exact amount that the trade will be simulated for.
	// For example, if the user wants to trade 100 OSMO, this amount would be 100000000.
	//
	// example: 100000000
	TokenInAmount string
	// The exact routes to use for the trade. These are the gamm routes used by Osmosis DEX.
	// example: [{"pool_id":1,"token_out_denom":"uosmo"}]
	Routes []SwapAmountInRoute
	// Will be the exact amount to submit on-chain as the minimum amount out for the trade
	//
	// example: 1
	TokenOutMinAmount string
	// The swap simulation will be performed as if it is for this wallet.
	//
	// example: osmo1njsthsamgkzdqqg4awj7y6rxtuk3q26t2t8md2
	UserWallet string
	// On Osmosis, you must specify a wallet with at least 100 OSMO in it. If this param is left empty, the UserWallet will be used instead.
	// NOTE: At present, on Osmosis, only arbing with the OSMO token is supported. As we will support additional tokens in the future, we use
	// this parameter so that callers can indicate which arbitrage tokens they want to use.
	//
	// example: osmo14mmus5h7m6vkp0pteks8wawaj4wf3sx7fy3s2r
	ArbitrageWallet string
	// Setting SkipWalletFundsCheck to true will run the simulation regardless of whether the user's wallet has sufficient funds to perform the TX.
	// If this param is false, and the user wants to simulate a trade of 100 JUNO for X OSMO, the user's wallet must have at least 100 JUNO in it.
	//
	// example: true
	SkipWalletFundsCheck bool
}

type SwapAmountInRoute struct {
	Pool          string `json:"pool_id,omitempty" yaml:"pool_id"`
	TokenOutDenom string `json:"token_out_denom,omitempty" yaml:"token_out_denom"`
}

// Results of the simulation
// swagger:model SimulatedSwapResult
type SimulatedSwapResult struct {
	// the user's swap including the most efficient routes (pools) to use
	SimulatedUserSwap *SimulatedSwap `json:"userSwap,omitempty"`
	// how much arbitrage the user's swap will cause, routes to use, etc
	ArbitrageSwap *ArbitrageSwap `json:"arbitrageSwap,omitempty"`
	// Whether or not the user's swap will cause arbitrage once executed on chain
	HasArbitrageOpportunity bool
	// if there was some issue detected on the server
	Error string `json:"error,omitempty"`
}

type ArbitrageSwap struct {
	// the arbitrage swap including the most efficient routes (pools) to use
	SimulatedSwap *SimulatedSwap
	// e.g. 165.1269 OSMO
	//
	// example: 165.1269 OSMO
	EstimatedProfitHumanReadable string
	// e.g. 165.1269
	//
	// example: 165.1269
	EstimatedProfitBaseAmount string
}

type SimulatedSwap struct {
	// Will be the exact amount/denomination to submit on-chain for the trade
	TokenIn cosmosTypes.Coin
	// Will be the exact amount to submit on-chain as the minimum amount out for the trade
	TokenOutMinAmount cosmosTypes.Int
	//Comma separated list of pools that will be traded through (only for human readable info)
	Pools string
	// The exact routes to use for the trade. These are the gamm routes used by Osmosis DEX.
	// example: [{"pool_id":1,"token_out_denom":"uosmo"}]
	Routes types.SwapAmountInRoutes `json:"routes,omitempty"`
	// Will be the simulated amount that will be received when this trade is submitted to the chain.
	TokenOutAmount cosmosTypes.Int
	// One of the 'denom' from asset lists at https:// github.com/osmosis-labs/assetlists/tree/main/osmosis-1
	TokenOutDenom string
	// One of the 'symbol' from asset lists at https:// github.com/osmosis-labs/assetlists/tree/main/osmosis-1
	TokenInSymbol string
	// example: 165.1269 OSMO
	AmountOutHumanReadable string
	// One of the 'symbol' from asset lists at https:// github.com/osmosis-labs/assetlists/tree/main/osmosis-1
	TokenOutSymbol string
	// example: 165.1269
	BaseAmount string
	// Amount this trade impacts the pool prices. For example, .025 would mean a 2.5% impact.
	// example: .025
	PriceImpact float64
}
