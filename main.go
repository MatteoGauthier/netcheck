package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"net"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func localAddresses(showIPv6 bool) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		return
	}

	var rows [][]string

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		}

		for _, a := range addrs {
			// Skip IPv6 addresses if the flag is not set
			if !showIPv6 {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() == nil {
					continue
				}
			}

			rows = append(rows, []string{
				fmt.Sprintf("%d", i.Index),
				i.Name,
				a.String(),
				i.HardwareAddr.String(),
			})
		}
	}

	var (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")

		headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		cellStyle    = lipgloss.NewStyle().Padding(0, 1)
		oddRowStyle  = cellStyle.Foreground(gray)
		evenRowStyle = cellStyle.Foreground(lightGray)
	)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("Index", "Interface", "Address", "MAC").
		Rows(rows...)

	fmt.Println(t)

}

func main() {
	cmd := &cli.Command{
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "ipv6",
				Usage:   "Show IPv6 addresses",
				Aliases: []string{"6"},
				Value:   false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			localAddresses(cmd.Bool("ipv6"))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
