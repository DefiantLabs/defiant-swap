# defiant-swap

## Authors
- [Kyle Moser](https://github.com/KyleMoser)

## TLDR

`defiant-swap` is a CLI tool for swapping tokens on Osmosis dex. It maximizes rates by capturing arbitrage at the time of trade. This arbitrage is awarded to the user instead of others (bots, protocols)

## How it works
- The larger your trade is, the more beneficial this tool is. Generally, high slippage = worse rates for users = more arbitrage profit for bots. 
This tool captures arbitrage revenue at the time of the trade, giving it back to users instead of leaving it on-chain for bots to capture. This is launching as a `FREE` service for users. Please do not abuse.

- The tool will simulate your trade ahead of time (before performing the trade) and tell you how much extra $$ you will get through arbitrage. You do not have to perform the trade, if you wish, you are free to test out the simulator. 

- This service is supported by a proprietary backend REST API (and trading algorithm). Please do not abuse this tool by invoking our REST API in an automated fashion unless you have been approved by our team (info@defiantlabs.net). 


## Download
Release are avilable [here](https://github.com/DefiantLabs/defiant-swap/releases)

## BUILD
To compile with ledger support, run:
```
go build -tags "netgo ledger" -ldflags '-X "github.com/cosmos/cosmos-sdk/version.BuildTags=netgo ledger," -w -s' -o defiant-swap main.go
```
Or without ledger support:
```
go build -ldflags '-w -s' -o defiant-swap main.go
```
## Ledger Info

By default, the web server (arbitrage/swap estimator) will check that the user's wallet has enough funds to perform the swap. Setting the `verify funds` param to false will estimate the swap without checking the user's wallet funds (useful for simulating swaps). 


To list your ledger keys run 
```
./defiant-swap ledger
``` 

To add a ledger key run
```
./defiant-swap ledger yourname
``` 
To remove a ledger key run 
```
./defiant-swap ledger yourname --delete
```

USAGE (with ledger):
```
./defiant-swap swap --in OSMO --out IXO --amount-in 1 --min-amount-out 1 --from kyle --ledger=true
```

USAGE (no ledger):
```
 ./defiant-swap swap --in AKT --out OSMO --amount-in 100000 --min-amount-out 1 --from arb --keyring-backend test --verify-funds=false
 ```

