package lib

import (
	"fmt"
	"net"

	"github.com/jackpal/gateway"
	netmon "tailscale.com/net/netmon"
)

type DefaultInterfaceInfo struct {
	Name       string
	SubnetMask string
}

func GetDefaultInterface(includeIPv6 bool) (string, error) {
	defaultInterface, err := netmon.DefaultRouteInterface()
	if err != nil {
		return "", fmt.Errorf("failed to get default interface: %w", err)
	}

	return defaultInterface, nil
}

func GetDefaultInterfaceInfo() (*DefaultInterfaceInfo, error) {
	defaultInterfaceName, err := GetDefaultInterface(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get default interface: %w", err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Name != defaultInterfaceName {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				// Get the subnet mask
				mask := ipnet.Mask
				if len(mask) == 4 {
					subnetMask := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
					return &DefaultInterfaceInfo{
						Name:       defaultInterfaceName,
						SubnetMask: subnetMask,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("subnet mask not found for default interface")
}

func GetSubnetMask() (string, error) {
	info, err := GetDefaultInterfaceInfo()
	if err != nil {
		return "", err
	}
	return info.SubnetMask, nil
}

func LocalAddresses(showIPv6 bool, showVirtual bool) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		return
	}

	defaultInterfaceInfo, err := GetDefaultInterfaceInfo()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
		return
	}
	defaultInterface := defaultInterfaceInfo.Name

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
				if IsLikelyVirtual(i.Name) {
					continue
				}
			}

			interfaceName := i.Name

			if !showVirtual && IsLikelyVirtual(i.Name) {
				interfaceName = fmt.Sprintf("%s (virtual)", i.Name)
			}

			rows = append(rows, []string{
				interfaceName,
				a.String(),
				i.HardwareAddr.String(),
			})
		}
	}

	PrintLocalAddresses(rows, defaultInterfaceIndex)
}

func GetNetwork(shouldPing bool) {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Println(fmt.Errorf("gateway error: %w", err))
		return
	}

	subnet, err := GetSubnetMask()
	if err != nil {
		fmt.Println(fmt.Errorf("subnet mask error: %w", err))
		return
	}

	pingAvg := ""
	if shouldPing {
		pingTime, pingErr := Ping(gw.String())
		if pingErr == nil {
			pingAvg = pingTime.String()
		}
	}

	PrintNetworkInfo(gw.String(), subnet, pingAvg)
}
