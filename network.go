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

// GetFreePort returns a free TCP port on the local machine by binding to a random port and immediately releasing the connection
func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	addr := listener.Addr().String()
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, err
	}
	return port, nil
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
