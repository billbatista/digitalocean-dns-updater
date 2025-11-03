package ipgrabber

import (
	"fmt"
	"io"
	"net/http"
)

type IPGrabber interface {
	GetPublicIP() (string, error)
}

type IFConfig struct{}

func (i *IFConfig) GetPublicIP() (string, error) {
	// Implementation to get public IP from https://ifconfig.co/
	resp, err := http.Get("https://ifconfig.co/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get public IP: %s", resp.Status)
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// ip has a new line character at the end, so we trim it
	ipStr := string(ip)
	ipStr = ipStr[:len(ipStr)-1]

	return ipStr, nil
}
