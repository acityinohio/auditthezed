// Audits the Zcash blockchain using naive RPC calls
// Requires txindex=1 option set in your zcash.conf

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

//a very limited TX struct for use in JSON unmarshaling
type TX struct {
	Txid         string
	Vin          []TXInput
	Vout         []TXOutput
	Vjoinsplit   []TXJoinSplit
	ValueBalance float64
}

type TXInput struct {
	Txid string
	Vout int
}

type TXOutput struct {
	ValueZat int64
	N        int
}

type TXJoinSplit struct {
	Vpub_old float64
	Vpub_new float64
}

func main() {
	//Iterate through blocks, tracking public UTXOs and amounts in shielded pool, summing them
	height := readHeight()
	fmt.Printf("zcashd says current height is %v, auditing to that height.\n", height)
	scanThePlanet(height)
	return
}

//pretty-print state of audit
func printAudit(pubZats int64, sproutZats int64, saplingZats int64, height int) {
	maxZats := calcMaxZats(height)
	fmt.Printf("At height %v:\n", height)
	fmt.Println("Maximum Allowed Zatoshis:", maxZats)
	fmt.Println("Public + Shielded:", pubZats+sproutZats+saplingZats)
	fmt.Println("Public UTXO Zatoshis:", pubZats)
	fmt.Println("Sprout Zatoshis:", sproutZats)
	fmt.Println("Sapling Zatoshis:", saplingZats)
	if pubZats+sproutZats <= maxZats {
		//https://twitter.com/taoeffect/status/1094716402991132672
		fmt.Println("All good! Everything checks out ok 👍")
	} else {
		fmt.Println("Ruh roh")
	}
}

//reads current height from zcashd
func readHeight() (height int) {
	var readOut struct {
		Blocks int
	}
	//get height from zcashd
	out, err := exec.Command("zcash-cli", "getblockchaininfo").Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(out, &readOut)
	if err != nil {
		log.Fatal(err)
	}
	height = readOut.Blocks
	return
}

//Calculates maximum possible Zatoshis for a given height
func calcMaxZats(height int) (maxZats int64) {
	blockReward := int64(62500)
	for n := 1; n <= height; n++ {
		maxZats += blockReward
		//slow start, skip 6.25 ZEC reward per https://github.com/zcash/zcash/issues/762#issuecomment-211884324
		if n == 9999 {
			blockReward += 125000
		} else if n < 19999 {
			blockReward += 62500
		} else {
			//halvenings
			if (n-10000)%840000 == 0 {
				blockReward = blockReward / 2
			}
		}
	}
	return
}

//scans the Zcash blockchain up to height, keeps track of UTXOs and shielded pool, adds them
//lets hope this doesn't blow up my computer
func scanThePlanet(height int) (pubZats, sproutZats, saplingZats int64) {
	UTXOs := make(map[string]int64)
	type Block struct {
		Tx []TX
	}
	//brief test
	for i := 1; i <= height; i++ {
		var currentBlock Block
		//periodic pause for progress
		if i%10000 == 0 {
			//in-progress audit
			printAudit(pubZats, sproutZats, saplingZats, i)
			if i < height {
				fmt.Println("Haven't reached tip of blockchain, continuing...")
			} else {
				return
			}
		}
		out, err := exec.Command("zcash-cli", "getblock", strconv.Itoa(i), "2").Output()
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(out, &currentBlock)
		if err != nil {
			log.Fatal(err)
		}
		//iterate through every transaction in a block and update totals
		for _, tx := range currentBlock.Tx {
			//For each Public TXInput, delete from UTXO map and subtract value from pubZats
			for _, vin := range tx.Vin {
				lookup := vin.Txid + "+" + strconv.Itoa(vin.Vout)
				pubZats -= UTXOs[lookup]
				delete(UTXOs, lookup)
			}
			//For each Public TXOutput, add to UTXO map and add value to pubZats
			for _, vout := range tx.Vout {
				lookup := tx.Txid + "+" + strconv.Itoa(vout.N)
				UTXOs[lookup] = vout.ValueZat
				pubZats += vout.ValueZat
			}
			//For each JoinSplit, add or subtract to sprout pool
			for _, vjoin := range tx.Vjoinsplit {
				sproutZats += (int64(vjoin.Vpub_old*1e8) - int64(vjoin.Vpub_new*1e8))
			}
			//For valueBalance, add or subtract to Sapling pool
			if tx.ValueBalance != 0 {
				saplingZats -= int64(tx.ValueBalance * 1e8)
			}
		}
	}
	//print results
	printAudit(pubZats, sproutZats, saplingZats, height)
	return
}
