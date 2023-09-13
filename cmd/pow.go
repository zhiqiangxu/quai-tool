package cmd

import (
	"fmt"
	"strconv"

	"math/big"

	"github.com/urfave/cli"
)

// PowCmd ...
var PowCmd = cli.Command{
	Name:  "pow",
	Usage: "pow actions",
	Subcommands: []cli.Command{
		hashRateCmd,
	},
}

var hashRateCmd = cli.Command{
	Name:   "hash_rate",
	Usage:  "hash_rate",
	Action: hashRate,
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

	count := big.NewInt(0).Lsh(big.NewInt(1), uint(diff.BitLen()-1))
	// fmt.Println("bitlen", diff.BitLen(), "time", time, "count", count)
	fmt.Println(big.NewInt(0).Div(count, big.NewInt(int64(time))))
	return
}
