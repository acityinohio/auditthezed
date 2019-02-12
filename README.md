# Audit the Zed

This small Go program audits the Zcash supply by using naive RPC calls. Effectively reconstructs the UTXO set in-memory in a map by iterating through every transaction in every block through `zcashd`, and keeps track of each shielded pool based on entry and exit into transparent addresses.

## Requirements

You need Golang installed (tested on version 1.11.5), a fully synced `zcashd` node, `zcashd` to be in your `$PATH`, and the `txindex=1` option set in your `zcash.conf`

Also note that it is _slow_ and _memory intensive_. This is an extremely naive implementation; based on my testing for around ~470,000 blocks it may use up to ~12GB of free RAM to conduct the audit (in addition to whatever memory `zcashd` is using). The sample audit added to this repo took 135m to compute. YMMV.

## Instructions

Have `zcashd` running in another process, then run the main script:

`go run main.go`

Then grab a coffee or two...or ten. It'll be a while.

## Contribution Guidelines

PRs are welcome, particularly for people who make it more efficient. :) 

### Sample Output

```
zcashd says current height is 479827, auditing to that height.
At height 10000:
Maximum Allowed Zatoshis: 3125312562500
Public + Shielded: 3113242426967
Public UTXO Zatoshis: 3089689326427
Sprout Zatoshis: 23553100540
Sapling Zatoshis: 0
All good! Everything checks out ok 👍
Haven't reached tip of blockchain, continuing...
At height 20000:
Maximum Allowed Zatoshis: 12501250000000
Public + Shielded: 12485006934974
Public UTXO Zatoshis: 12419716255402
Sprout Zatoshis: 65290679572
Sapling Zatoshis: 0
All good! Everything checks out ok 👍
...
...
At height 479827:
Maximum Allowed Zatoshis: 587285000000000
Public + Shielded: 587269844155807
Public UTXO Zatoshis: 552368082693859
Sprout Zatoshis: 25738780979868
Sapling Zatoshis: 9162980482080
All good! Everything checks out ok 👍
```
