package main

import (
	"context"
	"fmt"
	"log"
	"netcheck/lib"
	"os"

	"github.com/jackpal/gateway"
	"github.com/urfave/cli/v3"

	netmon "tailscale.com/net/netmon"

	"net"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func GetDefaultInterface(includeIPv6 bool) (string, error) {
	defaultInterface, err := netmon.DefaultRouteInterface()
	if err != nil {
		return "", fmt.Errorf("failed to get default interface: %w", err)
	}

	return defaultInterface, nil
}

func localAddresses(showIPv6 bool, showVirtual bool) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		return
	}

	defaultInterface, err := GetDefaultInterface(showIPv6)
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		return
	}

	var rows [][]string

	var defaultInterfaceIndex int

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		}

		for _, a := range addrs {

			if !showIPv6 {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() == nil {
					continue
				}
			}

			isDefaultInterface := i.Name == defaultInterface

			if isDefaultInterface {
				defaultInterfaceIndex = len(rows)
			}

			if !showVirtual && !isDefaultInterface {
				if lib.IsLikelyVirtual(i.Name) {
					continue
				}
			}

			interfaceName := i.Name

			if !showVirtual && lib.IsLikelyVirtual(i.Name) {
				interfaceName = fmt.Sprintf("%s (virtual)", i.Name)
			}

			rows = append(rows, []string{
				interfaceName,
				a.String(),
				i.HardwareAddr.String(),
			})
		}
	}

	var (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")

		headerStyle    = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		cellStyle      = lipgloss.NewStyle().Padding(0, 1)
		oddRowStyle    = cellStyle.Foreground(gray)
		evenRowStyle   = cellStyle.Foreground(lightGray)
		highlightStyle = cellStyle.Foreground(lipgloss.Color("10")).Bold(true)
	)

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row == defaultInterfaceIndex:
				return highlightStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("Interface", "Address", "MAC").
		Rows(rows...)

	fmt.Println(t)

}

func printGateway() {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Println(fmt.Errorf("gateway error: %w", err))
		return
	}
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Align(lipgloss.Center)
	fmt.Println(boxStyle.Render("Gateway: " + gw.String()))
}

func main() {
	cmd := &cli.Command{
		Name:                   "netcheck",
		Version:                "0.0.1",
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			localAddresses(cmd.Bool("ipv6"), cmd.Bool("virtual"))
			printGateway()
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
