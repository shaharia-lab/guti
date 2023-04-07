package gotil

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
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

// Ping sends an ICMP echo request to a given host and returns true if the host responds to the request within a specified timeout
func Ping(host string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("ip4:icmp", host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	msg := make([]byte, 64)
	msg[0] = 8 // echo
	msg[1] = 0 // code
	msg[2] = 0 // checksum
	msg[3] = 0 // checksum
	id := os.Getpid() & 0xffff
	msg[4] = byte(id >> 8)
	msg[5] = byte(id & 0xff)
	msg[6] = 0 // sequence number
	msg[7] = 0 // sequence number
	checksum := checkSum(msg)
	msg[2] = byte(checksum >> 8)
	msg[3] = byte(checksum & 0xff)
	conn.SetDeadline(time.Now().Add(timeout))
	_, err = conn.Write(msg[0:8])
	if err != nil {
		return false
	}
	_, err = conn.Read(msg[0:8])
	if err != nil {
		return false
	}

	return true
}

func checkSum(msg []byte) uint16 {
	sum := uint32(0)
	for i := 0; i < len(msg)-1; i += 2 {
		sum += uint32(msg[i])<<8 | uint32(msg[i+1])
	}
	if len(msg)%2 == 1 {
		sum += uint32(msg[len(msg)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return uint16(^sum)
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

// GetMACAddress retrieves the MAC address of a network interface
func GetMACAddress(iface string) (string, error) {
	netInterface, err := net.InterfaceByName(iface)
	if err != nil {
		return "", err
	}
	return netInterface.HardwareAddr.String(), nil
}

// GetHTTPStatusCode  retrieves the HTTP status code of a given URL.
func GetHTTPStatusCode(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

// HttpRequestWithRetry sends an HTTP request with retries
func HttpRequestWithRetry(client *http.Client, req *http.Request, maxRetries int, sleepDuration time.Duration) ([]byte, error) {
	var err error
	var resp *http.Response
	var body []byte

	for i := 0; i < maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			return body, err
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
