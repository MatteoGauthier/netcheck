package lib

import (
	"strings"
)

func IsLikelyVirtual(name string) bool {
	// Handle potential nil or empty names gracefully
	if name == "" {
		return false
	}

	// --- Exact Matches ---
	// 'any' is a special pseudo-device for capturing on all interfaces (virtual in concept)
	if name == "any" {
		return true
	}
	// Loopback is always virtual
	if name == "lo" || name == "lo0" {
		return true
	}

	// --- Prefix Matches ---
	// List of common prefixes for virtual interfaces across OSes
	virtualPrefixes := []string{
		"tun",     // TUN device (VPNs, tunnels)
		"tap",     // TAP device (VPNs, tunnels, virtualization)
		"veth",    // Virtual Ethernet Pair (containers, namespaces)
		"br",      // Bridge interface
		"docker",  // Docker-managed interface (often a bridge)
		"virbr",   // libvirt bridge (KVM/QEMU VMs)
		"vmnet",   // VMware virtual network interface
		"vboxnet", // VirtualBox host-only network
		"utun",    // macOS tunnel interface (VPNs)
		"bond",    // Bonded/teamed interface master
		"team",    // Teaming interface master
		"gre",     // GRE tunnel interface
		"ipsec",   // IPsec interface (often VPN related)
		"ppp",     // Point-to-Point Protocol (dial-up, DSL, some VPNs)
		"nas",     // Linux ATM bridging (often xDSL)
		"awdl",    // Apple Wireless Direct Link (macOS)
		"llw",     // Apple Low Latency WLAN (macOS)
		"gif",     // Generic tunnel interface (macOS/BSD IPv6 tunneling)
		"stf",     // 6to4 tunnel interface (macOS/BSD IPv6 tunneling)
		"p2p",     // Peer-to-peer (often Wi-Fi Direct related, macOS/Linux)
		"ap",      // Access Point mode interface (e.g., ap0, ap1 on macOS/Linux)
		"anpi",    // Apple Network Protocol Interface (macOS)
		"faith",   // IPv6 transition mechanism interface (BSD)
		"wg",      // WireGuard VPN interface often named wgX
		"ip_vti",  // Linux Virtual Tunnel Interface
		// Add other specific prefixes you encounter
	}

	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(name, prefix) {
			// Small refinement: ensure it's not just a coincidental prefix for 'lo'
			if prefix == "lo" && !(name == "lo" || (len(name) > 2 && name[2] >= '0' && name[2] <= '9')) {
				continue // Avoid matching 'longname' based on 'lo' prefix
			}
			return true
		}
	}

	// --- Pattern Matches ---
	// Check for VLAN interfaces (e.g., eth0.100, enp3s0.50)
	// These are virtual sub-interfaces on a physical parent.
	if dotIndex := strings.LastIndexByte(name, '.'); dotIndex != -1 && dotIndex < len(name)-1 {
		suffix := name[dotIndex+1:]
		// Check if the suffix consists only of digits
		allDigits := true
		for _, r := range suffix {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			// It looks like a VLAN sub-interface name (*.VLANID)
			return true
		}
		// Alternative precise check using strconv (import "strconv")
		// if _, err := strconv.Atoi(suffix); err == nil {
		//	 return true // Suffix is a number, likely a VLAN ID
		// }
	}

	// If none of the above matched, assume it's likely physical or unknown
	return false
}

// Optional: Helper for Windows using the Description field
// This is less reliable than name patterns but necessary for some Windows virtual adapters.
func isWindowsVirtualByDescription(description string) bool {
	if description == "" {
		return false
	}
	lowerDesc := strings.ToLower(description)
	windowsVirtualSubstrings := []string{
		"virtual",  // General catch-all
		"loopback", // Windows Loopback Adapter
		"tap-",     // Common prefix for TAP drivers (e.g., OpenVPN)
		"tap adapter",
		"wan miniport", // Used for VPN/dial-up connections
		"hyper-v",      // Hyper-V virtual switch/adapter
		"vmware",       // VMware virtual adapter
		"virtualbox",   // VirtualBox adapter
		"vbox",
		"pptp",         // Point-to-Point Tunneling Protocol (VPN)
		"l2tp",         // Layer 2 Tunneling Protocol (VPN)
		"ikev2",        // Internet Key Exchange v2 (VPN)
		"sstp",         // Secure Socket Tunneling Protocol (VPN)
		"anchorfree",   // Some VPN providers install drivers with this name
		"kernel debug", // Windows Kernel Debug Network Adapter
	}

	for _, sub := range windowsVirtualSubstrings {
		if strings.Contains(lowerDesc, sub) {
			return true
		}
	}
	return false
}
