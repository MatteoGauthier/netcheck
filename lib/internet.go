package lib

import (
	"io"
	"net/http"
	"strings"
	"time"

	"sync"
	"tailscale.com/health"
	"tailscale.com/net/dns"
)

func GetPublicIP(includeIPv6 bool) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	var url string
	if includeIPv6 {
		url = "https://api64.ipify.org?format=text"
	} else {
		url = "https://api.ipify.org?format=text"
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

func GetDNSServers() ([]string, error) {
	logf := func(format string, args ...interface{}) {}

	healthTracker := &health.Tracker{}
	osConfig, err := dns.NewOSConfigurator(logf, healthTracker, nil, "")
	if err != nil {
		return []string{"Unable to detect"}, nil
	}
	defer osConfig.Close()

	// Get the base DNS configuration from the OS
	baseConfig, err := osConfig.GetBaseConfig()
	if err != nil {
		return []string{"Unable to detect"}, nil
	}

	var dnsServers []string
	for _, nameserver := range baseConfig.Nameservers {
		dnsServers = append(dnsServers, nameserver.String())
	}

	if len(dnsServers) == 0 {
		return []string{"No DNS servers found"}, nil
	}

	return dnsServers, nil
}

func GetInternetConnectivity(includeIPv6 bool) {

	var wg sync.WaitGroup
	var publicIP string
	var internetPing time.Duration
	var dnsServers []string
	var publicIPErr, pingErr, dnsErr error

	wg.Add(3)

	// Get public IP
	go func() {
		defer wg.Done()
		publicIP, publicIPErr = GetPublicIP(includeIPv6)
	}()

	// Ping internet (using Cloudflare DNS)
	go func() {
		defer wg.Done()
		internetPing, pingErr = Ping("1.1.1.1")
	}()

	// Get DNS servers
	go func() {
		defer wg.Done()
		dnsServers, dnsErr = GetDNSServers()
	}()

	wg.Wait()

	publicIPInfo := "N/A"
	if publicIPErr == nil {
		publicIPInfo = publicIP
	}

	pingInfo := "N/A"
	if pingErr == nil {
		pingInfo = internetPing.String()
	}

	dnsInfo := "N/A"
	if dnsErr == nil && len(dnsServers) > 0 {
		dnsInfo = strings.Join(dnsServers, ", ")
	}

	PrintInternetConnectivity(publicIPInfo, pingInfo, dnsInfo)
}
