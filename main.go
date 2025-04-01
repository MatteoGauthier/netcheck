package main

import (
	"context"
	"fmt"
	"log"
	"netcheck/lib"
	"os"

	"github.com/urfave/cli/v3"

	"net"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type DefaultInterface struct {
	Name       string
	IPAddress  string
	MACAddress string
}

func GetDefaultInterface(includeIPv6 bool) (*DefaultInterface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if !includeIPv6 {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() == nil {
					continue
				}
			}

			return &DefaultInterface{
				Name:       iface.Name,
				IPAddress:  addr.String(),
				MACAddress: iface.HardwareAddr.String(),
			}, nil
		}
	}

	return nil, fmt.Errorf("no default interface found")
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

			isDefaultInterface := i.Name == defaultInterface.Name

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
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
