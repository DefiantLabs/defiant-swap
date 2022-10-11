Author: Kyle Moser, Defiant Labs

CLI tool for swapping tokens on Osmosis blockchain. Maximizes rates by capturing arbitrage (users capture arbitrage profits instead of bots).
The larger your trade is, the more beneficial this tool is. Generally, high slippage = worse rates for users = more arbitrage profit for bots. 
This tool captures arbitrage revenue at the time of the trade, giving it back to users instead of leaving it on-chain for bots to capture. 

This is a FREE service. No fees, profits, etc. go to Defiant Labs. In the future Defiant Labs will be looking to partner with DEXs
and alternate trading platforms that wish to maximize returns/rates for their users. Contact us at info@defiantlabs.net.

This service is supported by a proprietary backend REST API (and trading algorithm). Please do not abuse this tool by invoking our REST API
in an automated fashion e.g. with code/bots. If you wish to use our API for automation purposes, please contact us first. 

USAGE:
 ./defiant-swap swap --in AKT --out OSMO --amount-in 100000 --min-amount-out 1 --from arb --keyring-backend test --verify-funds=false

The verify funds param does not need to be set to false. However, if true (default), it will check that the user's wallet has enough funds to perform the swap.
If false, it won't check and it will just do the simulation as if you have enough funds.