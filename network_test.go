package guti

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func TestGetPublicIP(t *testing.T) {
	ip, err := GetPublicIP()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check if the IP address is a valid IPv4 or IPv6 address
	if ipAddr := net.ParseIP(ip); ipAddr == nil {
		t.Errorf("invalid IP address: %s", ip)
	}
}

func TestGetFreePort(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		{
			name:     "Test_case_1",
			expected: 0,
		},
		{
			name:     "Test_case_2",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, err := GetFreePort()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if port <= 0 {
				t.Errorf("expected port greater than 0, but got %v", port)
			}
		})
	}
}

func TestGetLocalIPs(t *testing.T) {
	ps, err := GetLocalIPs()
	if err != nil {
		return
	}

	assert.Equal(t, true, len(ps) > 0)
}

func TestIsPortOpen(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		expected bool
	}{
		{
			name:     "Test_case_1",
			host:     "google.com",
			port:     80,
			expected: true,
		},
		{
			name:     "Test_case_2",
			host:     "example.com",
			port:     8080,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsPortOpen(tt.host, tt.port)
			if actual != tt.expected {
				t.Errorf("unexpected result: expected=%v, actual=%v", tt.expected, actual)
			}
		})
	}
}

func TestGetRemoteIP(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected string
		foundIp  bool
	}{
		{
			name:    "Test_case_1",
			host:    "google.com",
			foundIp: true,
		},
		{
			name:    "Test_case_2",
			host:    "somethningelse",
			foundIp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := GetRemoteIP(tt.host)
			if !tt.foundIp {
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, ip)
			}
		})
	}
}

func TestGetHTTPStatusCode(t *testing.T) {
	// Create a new test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Test the function with the temporary test server's URL
	statusCode, err := GetHTTPStatusCode(ts.URL)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("unexpected status code: expected=%v, actual=%v", http.StatusOK, statusCode)
	}
}

func TestHttpRequestWithRetry(t *testing.T) {
	// Create a new HTTP test server to mock the endpoint
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	}))
	defer testServer.Close()

	// Create a new HTTP client
	httpClient := &http.Client{Timeout: 2 * time.Second}

	// Create a new GET request to the test server's URL
	req, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Make the HTTP request with retries
	resp, err := HttpRequestWithRetry(httpClient, req, 3, 1*time.Second)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: expected=%d, actual=%d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expectedBody := "Hello, World!\n"
	if string(body) != expectedBody {
		t.Errorf("unexpected body: expected=%q, actual=%q", expectedBody, string(body))
	}
}

func TestGenerateRandomMacAddress(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool // expected return value
	}{
		{
			name:     "Test_case_1",
			expected: true,
		},
		{
			name:     "Test_case_2",
			expected: true,
		},
		{
			name:     "Test_case_3",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			macAddress, err := GenerateRandomMacAddress()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// The MAC address should be in the format XX:XX:XX:XX:XX:XX
			matched, err := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, macAddress)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !matched {
				t.Errorf("unexpected MAC address format: %s", macAddress)
			}
		})
	}
}

func TestSendUDPPacket(t *testing.T) {
	// Start mock server
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Start goroutine to read UDP packets from the mock server
	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}
			// Print the received message for debugging purposes
			fmt.Printf("Received message: %s\n", string(buf[:n]))
		}
	}()

	// Send UDP packet to the mock server
	payload := []byte("hello world")
	err = SendUDPPacket(payload, conn.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}

	// Wait for a short time to ensure the mock server receives the packet
	time.Sleep(100 * time.Millisecond)
}
