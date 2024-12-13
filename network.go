// Package guti contains packages
package guti

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

// GetPublicIP  uses an HTTP client to make a request to an API that returns the public IP address of the client. It then parses the response to extract the IP address and returns it as a string
func GetPublicIP() (string, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

// IsPortAvailable checks if a given port is available for use.
//
// This function attempts to bind to the specified port and then immediately
// releases it. It also performs a secondary check by trying to establish a
// connection to the port.
//
// Parameters:
//   - port: The port number to check for availability.
//
// Returns:
//   - bool: true if the port is available, false otherwise.
//
// Note: This function may return false positives in rare cases where the port
// becomes occupied immediately after the check. For more reliable results in
// concurrent environments, consider using port locking mechanisms.
func IsPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()

	// Double-check if the port is really available
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
	if err != nil {
		return true
	}
	if conn != nil {
		conn.Close()
	}
	return false
}

// GetFreePort returns a free TCP port on the local machine.
//
// This function attempts to bind to a random port assigned by the operating system,
// then checks its availability using IsPortAvailable.
//
// Returns:
//   - int: The available port number.
//   - error: An error if no port could be obtained or if the port is not available.
//
// Note: There's a small chance that the port might become unavailable immediately
// after this function returns. In concurrent environments, implement additional
// synchronization if necessary.
func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to bind to a port: %w", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, fmt.Errorf("failed to split host and port: %w", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert port to integer: %w", err)
	}

	if IsPortAvailable(port) {
		return port, nil
	}

	return 0, fmt.Errorf("port %d is not available", port)
}

// GetFreePortFromPortRange returns a free TCP port within the specified range.
//
// This function iterates through the given port range and returns the first
// available port it finds.
//
// Parameters:
//   - minPort: The lower bound of the port range to search (inclusive).
//   - maxPort: The upper bound of the port range to search (inclusive).
//
// Returns:
//   - int: An available port number within the specified range.
//   - error: An error if no free ports are found in the given range.
//
// Note: This function may be slower than GetFreePort for large ranges. Consider
// using GetFreePort if you don't need a port in a specific range.
func GetFreePortFromPortRange(minPort, maxPort int) (int, error) {
	for port := minPort; port <= maxPort; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no free ports found in range %d-%d", minPort, maxPort)
}

// GetLocalIPs returns a slice of IP addresses for all network interfaces on the local machine
func GetLocalIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	ips := []string{}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips, nil
}

// IsPortOpen checks if a given TCP port is open on a given host
func IsPortOpen(host string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetRemoteIP retrieves the IP address of a remote host
func GetRemoteIP(host string) (string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	return addrs[0], nil
}

// GetHTTPStatusCode  retrieves the HTTP status code of a given URL.
func GetHTTPStatusCode(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

// HTTPRequestWithRetry sends an HTTP request with retries
func HTTPRequestWithRetry(client *http.Client, req *http.Request, maxRetries int, sleepDuration time.Duration) (*http.Response, error) {
	var err error
	var resp *http.Response

	for i := 0; i < maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(sleepDuration)
	}

	return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, err)
}

// SendUDPPacket sends a UDP packet
func SendUDPPacket(payload []byte, address string) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(payload)
	if err != nil {
		return err
	}

	return nil
}

// GenerateRandomMacAddress generates a random MAC address
func GenerateRandomMacAddress() (string, error) {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	// Set the MAC address's multicast bit
	buf[0] |= 2

	// Set the MAC address's local assignment bit
	buf[0] &= 254

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5]), nil
}
