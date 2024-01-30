package cmd

import (
	"context"
	"fmt"
	"strconv"

	"math/big"

	"github.com/dominant-strategies/go-quai/common/hexutil"
	"github.com/dominant-strategies/go-quai/quaiclient"
	"github.com/urfave/cli"
)

// PowCmd ...
var PowCmd = cli.Command{
	Name:  "pow",
	Usage: "pow actions",
	Subcommands: []cli.Command{
		hashRateCmd,
		listHashRateCmd,
	},
}

var hashRateCmd = cli.Command{
	Name:   "hash_rate",
	Usage:  "hash_rate",
	Action: hashRate,
}

var listHashRateCmd = cli.Command{
	Name:   "list_hash_rate",
	Usage:  "list_hash_rate",
	Action: listHashRate,
}

func hashRate(ctx *cli.Context) (err error) {

	if len(ctx.Args()) != 2 {
		err = fmt.Errorf("hash_rate diff time")
		return
	}
	diff, ok := big.NewInt(0).SetString(ctx.Args()[0], 16)
	if !ok {
		diff, ok = big.NewInt(0).SetString(ctx.Args()[0], 0)
		if !ok {
			err = fmt.Errorf("invalid diff:%s", ctx.Args()[0])
			return
		}

	}
	if diff.BitLen() == 0 {
		err = fmt.Errorf("invalid diff")
		return
	}

	time, err := strconv.Atoi(ctx.Args()[1])
	if err != nil {
		return
	}

	// fmt.Println("bitlen", diff.BitLen(), "time", time, "count", count)
	fmt.Println(computeHashRate(diff, uint64(time)))
	return
}

var big2e256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

func computeHashRate(diff *big.Int, time uint64) *big.Int {
	if diff.Sign() <= 0 {
		panic("diff negative")
	}
	target := new(big.Int).Div(big2e256, diff)
	count := big.NewInt(0).Lsh(big.NewInt(1), uint(256-target.BitLen()))
	count.Mul(count, big.NewInt(0).Lsh(big.NewInt(1), uint(target.BitLen())))
	count.Div(count, target)
	return big.NewInt(0).Div(count, big.NewInt(int64(time)))
}

func listHashRate(ctx *cli.Context) (err error) {
	zoneRPCs := []string{
		"https://rpc.cyprus1.colosseum.quaiscan.io",
		"https://rpc.cyprus2.colosseum.quaiscan.io",
		"https://rpc.cyprus3.colosseum.quaiscan.io",
		"https://rpc.paxos1.colosseum.quaiscan.io",
		"https://rpc.paxos2.colosseum.quaiscan.io",
		"https://rpc.paxos3.colosseum.quaiscan.io",
		"https://rpc.hydra1.colosseum.quaiscan.io",
		"https://rpc.hydra2.colosseum.quaiscan.io",
		"https://rpc.hydra3.colosseum.quaiscan.io",
	}

	anchor := uint64(100)
	totalRate := big.NewInt(0)
	for i, rpc := range zoneRPCs {
		client, err := quaiclient.Dial(rpc)
		if err != nil {
			return err
		}

		currentHeader := client.HeaderByNumber(context.Background(), "latest")
		if currentHeader == nil {
			panic(fmt.Sprintf("failed to fetch currentHeader from %s", rpc))
		}

		anchorHeader := client.HeaderByNumber(context.Background(), hexutil.EncodeUint64(currentHeader.NumberU64(2 /*ZONE_CTX*/)-anchor))
		if anchorHeader == nil {
			panic(fmt.Sprintf("failed to fetch anchorHeader from %s", rpc))
		}

		avgDiff := big.NewInt(0)
		avgDiff.Add(currentHeader.Difficulty(), anchorHeader.Difficulty())
		avgDiff.Div(avgDiff, big.NewInt(2))
		// fmt.Printf("cur diff:%v anchor diff:%v avg diff:%v\n", currentHeader.Difficulty(), anchorHeader.Difficulty(), avgDiff)

		avgTime := (currentHeader.Time() - anchorHeader.Time()) / (anchor - 1)

		rate := computeHashRate(avgDiff, avgTime)
		totalRate.Add(totalRate, rate)
		fmt.Printf("zone %d rpc:%v rate:%v\n", i, rpc, rate)

	}

	fmt.Println("totalRate", totalRate)

	return
}
