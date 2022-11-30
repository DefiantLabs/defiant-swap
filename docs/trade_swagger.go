package docs

import "github.com/DefiantLabs/OsmosisArbitrageCLI/query"

// swagger:route POST /api/token authentication jwt
// Get a JWT that must be included with future requests
// responses:
//
//	200: JWT
//
// This response will contain the recommended trade to execute on the user's behalf, plus a recommended second arbitrage trade to execute.
// Both trades should be executed within a single transaction (as separate messages) to guarantee profits go to the caller.
// swagger:response JWT
type jwtResponseWrapper struct {
	// in:body
	// required: true
	Body struct{ query.JWT }
}

// swagger:parameters jwt
type jwtParamsWrapper struct {
	// JWTRequest containing arbitrage wallet address (can be set to anything, is not verified)
	// in:body
	// required: true
	Body struct {
		query.JWTRequest
	}
}

// swagger:route POST /api/secured/estimateswap simulations swapRequest
// Simulate the user's swap, providing the input and output token symbols and the amount the user wants to trade.
// responses:
//
//	200: SimulatedSwapResult
//
// This response will contain the recommended trade to execute on the user's behalf, plus a recommended second arbitrage trade to execute.
// Both trades should be executed within a single transaction (as separate messages) to guarantee profits go to the caller.
// swagger:response SimulatedSwapResult
type swapResponseWrapper struct {
	// in:body
	// required: true
	Body struct{ query.SimulatedSwapResult }
}

// swagger:parameters swapRequest
type swapParamsWrapper struct {
	// The swap that will be simulated (no on chain transactions will be generated).
	// in:body
	// required: true
	Body struct {
		query.SimulatedSwapRequest
	}
}

// swagger:route POST /api/secured/estimatewithpools simulations swapExactRequest
// Simulate the user's swap, providing the exact routes (pools), denoms, and amounts.
// responses:
//
//	200: SimulatedSwapResult
//
// This response will contain the recommended trade to execute on the user's behalf, plus a recommended second arbitrage trade to execute.
// Both trades should be executed within a single transaction (as separate messages) to guarantee profits go to the caller.
// swagger:response SimulatedSwapResult
type swapExactResponseWrapper struct {
	// in:body
	// required: true
	Body struct{ query.SimulatedSwapResult }
}

// swagger:parameters swapExactRequest
type swapExactParamsWrapper struct {
	// The swap that will be simulated (no on chain transactions will be generated).
	// in:body
	// required: true
	Body struct {
		query.SimulatedSwapExactPoolsRequest
	}
}
