package lib

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/lipgloss/table"
)

func PrintNetworkInfo(gateway string, subnet string, pingAvg string) {

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Align(lipgloss.Center)

	pingInfo := ""
	if pingAvg != "" {
		pingInfo = fmt.Sprintf(" (%s)", pingAvg)
	}

	gatewayBox := boxStyle.Render("Gateway: " + gateway + pingInfo)
	subnetBox := boxStyle.Render("Subnet Mask: " + subnet)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, gatewayBox, " ", subnetBox))

}

func PrintInternetConnectivity(publicIP string, ping string, dns string) {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("76")).
		Padding(0, 1).
		Align(lipgloss.Center)

	publicIPBox := boxStyle.Render("Public IP: " + publicIP)
	internetPingBox := boxStyle.Render("Internet Ping: " + ping)
	dnsBox := boxStyle.Render("DNS: " + dns)

	fmt.Println()
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, publicIPBox, " ", internetPingBox, " ", dnsBox))
}

func PrintLocalAddresses(rows [][]string, highlightIndex int) {

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
			case row == highlightIndex:
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
