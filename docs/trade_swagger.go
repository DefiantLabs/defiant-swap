package docs

import "github.com/DefiantLabs/JunoswapArbitrageCLI/query"

// swagger:route POST /foobar foobar-tag swapRequest
// Foobar does some amazing stuff.
// responses:
//   200: result

// This text will appear as description of your response body.
// swagger:response result
type swapResponseWrapper struct {
	// in:body
	// required: true
	Body struct{ query.SimulatedSwapResult }
}

// swagger:parameters swapRequest
type swapParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	// required: true
	Body struct {
		query.SimulatedSwapRequest
	}
}
