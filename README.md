# Audit the Zed

This small Go program audits the Zcash supply by using naive RPC calls. Effectively reconstructs the UTXO set in-memory in a map by iterating through every transaction in every block through zcashd, andkeeps track of each shielded pool based on entry and exit into transparent addresses.

## Requirements

You need Golang installed (tested on version 1.11.5), a fully synced zcashd node, zcashd to be in your $PATH, and the `txindex=1` option set in your `zcash.conf`

Also note that it is _slow_ and _memory intensive_. This is an extremely naive implementation; based on my testing it may use up to ~8GB of free RAM to conduct the audit (in addition to what memory zcashd is using). YMMV.

## Instructions

Have `zcashd` running in another process, then run the main script:

`go run main.go`

Then grab a coffee or two...or ten. It'll be a while.

## Contribution Guidelines

PRs are welcome, particularly for people who make it more efficient. :) 

### Sample Output to Block 12000

```
zcashd says current height is 12000, auditing to that height.
At height 10000:
Maximum Allowed Zatoshis: 3125312562500
Public + Shielded: 3113242426967
Public UTXO Zatoshis: 3089689326427
Sprout Zatoshis: 23553100540
Sapling Zatoshis: 0
All good! Everything checks out ok 👍
Haven't reached tip of blockchain, continuing...
At height 12000:
Maximum Allowed Zatoshis: 4500500062500
Public + Shielded: 4487975025225
Public UTXO Zatoshis: 4456226669316
Sprout Zatoshis: 31748355909
Sapling Zatoshis: 0
All good! Everything checks out ok 👍
```
