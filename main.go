package main

import (
	"context"
	"log"
	"netcheck/lib"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:                   "netcheck",
		Version:                "0.0.1",
		Usage:                  "Check quickly your network configuration",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "ipv6",
				Usage:   "Show IPv6 addresses",
				Aliases: []string{"6"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "virtual",
				Usage:   "Show virtual interfaces",
				Aliases: []string{"x"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "ping",
				Usage:   "Ping the gateway to check connectivity",
				Aliases: []string{"p"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "internet",
				Usage:   "Check internet connectivity (ping, public IP, DNS servers)",
				Aliases: []string{"i"},
				Value:   false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			lib.LocalAddresses(cmd.Bool("ipv6"), cmd.Bool("virtual"))
			lib.GetNetwork(cmd.Bool("ping"))
			if cmd.Bool("internet") {
				lib.GetInternetConnectivity(cmd.Bool("ipv6"))
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
