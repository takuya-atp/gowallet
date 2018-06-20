package main

import (
	"log"
	"os"

	"gowallet/ethereum"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gowallet"

	ethcl, err := ethereum.NewEthClient()
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "transfer",
			Aliases: []string{"t"},
			Usage:   "transfer wei to address",
			Action: func(c *cli.Context) error {
				pass := c.GlobalString("password")
				wei := c.GlobalInt("wei")
				addr := c.GlobalString("address")
				keyPath := c.GlobalString("keypath")
				ethcl.TransferWei(pass, wei, addr, keyPath)
				return nil
			},
		},
		{
			Name:    "balance",
			Aliases: []string{"b"},
			Usage:   "check wallet balance",
			Action: func(c *cli.Context) error {
				pass := c.GlobalString("password")
				keyPath := c.GlobalString("keypath")
				ethcl.ConfirmBalance(pass, keyPath)
				return nil
			},
		},
		{
			Name:    "account",
			Aliases: []string{"a"},
			Usage:   "create new wallet",
			Action: func(c *cli.Context) error {
				pass := c.GlobalString("password")
				ethcl.GetAccount(pass)
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "keypath, k",
			Usage: "pass for confirm your wallet balance and transfer wei from your wallet",
		},
		cli.StringFlag{
			Name:  "address, a",
			Usage: "address for the transfer",
		},
		cli.IntFlag{
			Name:  "wei, w",
			Usage: "transfering wei quantity from your wallet",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "password for new wallet",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
