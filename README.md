Author: Kyle Moser, Defiant Labs

CLI tool for swapping tokens on Osmosis blockchain. Maximizes rates by capturing arbitrage (users capture arbitrage profits instead of bots).
The larger your trade is, the more beneficial this tool is. Generally, high slippage = worse rates for users = more arbitrage profit for bots. 
This tool captures arbitrage revenue at the time of the trade, giving it back to users instead of leaving it on-chain for bots to capture. This is a FREE service for users. No fees, profits, etc. go to Defiant Labs. 

The tool will simulate your trade ahead of time (before performing the trade) and tell you how much extra $$ you will get through arbitrage. You do not have to perform the trade, if you wish, you are free to test out the simulator. 

This service is supported by a proprietary backend REST API (and trading algorithm). Please do not abuse this tool by invoking our REST API in an automated fashion unless you have been approved by our team (info@defiantlabs.net). 


BUILD:
You do not HAVE to build. You can just run the Linux binary provided (defiant-swap).
However, if you wish to compile with ledger support, run:
go build -tags "netgo ledger" -ldflags '-X "github.com/cosmos/cosmos-sdk/version.BuildTags=netgo ledger," -w -s' -o defiant-swap main.go
Or without ledger support:
go build -ldflags '-w -s' -o defiant-swap main.go

USAGE (setting up ledger):
You only need to follow these steps if you have never set up a Ledger before. If you already have a Ledger set up you can use it as normal.
To list your ledger keys run './defiant-swap ledger' with no other arguments.
To add a ledger key run './defiant-swap ledger kyle' (or your key name in place of kyle).
To remove a ledger key run './defiant-swap ledger kyle --delete'.

USAGE (with ledger):
./defiant-swap swap --in OSMO --out IXO --amount-in 1 --min-amount-out 1 --from kyle --ledger=true

USAGE (no ledger):
 ./defiant-swap swap --in AKT --out OSMO --amount-in 100000 --min-amount-out 1 --from arb --keyring-backend test --verify-funds=false

By default, the web server (arbitrage/swap estimator) will check that the user's wallet has enough funds to perform the swap.
Setting the verify funds param to false will estimate the swap without checking the user's wallet funds (useful for simulating swaps). 
