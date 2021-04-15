package netinfo

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type HttpTest struct {
	Port int `json:"port"`
}

// Get returns network information.
// When test != nil, it will run HTTP tests.
func Get(test *HttpTest) NetworkInfo {
	n := NetworkInfo{}
	n.IPv4Method, n.IPv6Method = "https://v4.ident.me/", "https://v6.ident.me/"
	n.IPv4, n.IPv6 = getIP(n.IPv4Method), getIP(n.IPv6Method)

	// Grab all network interfaces.
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		var addresses []string

		// Grab all iface addresses.
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}

		n.Ifaces = append(n.Ifaces, iface{
			MTU:          i.MTU,
			Name:         i.Name,
			PhysicalAddr: i.HardwareAddr.String(),
			Addresses:    addresses,
		})
	}

	// Perform HTTP tests.
	if test != nil {
		n.HttpTest = test.perform(n.IPv4, n.IPv6)
	}

	return n
}

type httpTestResult struct {
	IPv4Passed bool `json:"ipv4_passed"`
	IPv6Passed bool `json:"ipv6_passed"`
}

// perform runs HTTP tests and returns the httpTestResult.
func (t HttpTest) perform(ipv4, ipv6 string) httpTestResult {
	http.HandleFunc("/ping", pingHandler)

	// Run HTTP server in goroutine.
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", t.Port), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	result := httpTestResult{}

	// Perform IPv4 tests.
	_, err := http.Get(fmt.Sprintf("http://%s:%v/ping", ipv4, t.Port))
	result.IPv4Passed = err == nil

	// Perform IPv6 tests.
	_, err = http.Get(fmt.Sprintf("http://[%s]:%v/ping", ipv6, t.Port))
	result.IPv6Passed = err == nil

	return result
}

// NetworkInfo information about network.
type NetworkInfo struct {
	IPv4       string         `json:"ipv4"`
	IPv4Method string         `json:"ipv4_method"`
	IPv6       string         `json:"ipv6"`
	IPv6Method string         `json:"ipv6_method"`
	Ifaces     []iface        `json:"ifaces"`
	HttpTest   httpTestResult `json:"http_test"`
}

type iface struct {
	MTU          int      `json:"mtu"`
	Name         string   `json:"name"`
	PhysicalAddr string   `json:"physical_addr"`
	Addresses    []string `json:"addresses"`
}

// getIP make GET request to the url, which must return public IP address of the current machine.
func getIP(url string) string {
	res, err := http.Get(url)
	if err != nil {
		return ""
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return string(b)
}

// pingHandler handles ping request.
func pingHandler(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("pong"))
}
