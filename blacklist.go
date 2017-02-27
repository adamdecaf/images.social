package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"os"
)

// blacklist based on certain conditions
// - ip, read from file every N minutes
// - upload frequency, each upload is recorded

type Blacklist interface {
	Blocked(http.Request) bool
}

func NewBlacklist(filepath string) (Blacklist, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		return ipBlacklist{}, err
	}

	// Grab each line and parse it
	f, err := os.Open(filepath)
	if err != nil {
		return ipBlacklist{}, err
	}
	r := bufio.NewReader(f)

	ips := make([]net.IP, 0)
	cidrs := make([]net.IPNet, 0)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		ip, cidr, err := net.ParseCIDR(line)
		if err == nil {
			if cidr != nil {
				cidrs = append(cidrs, *cidr)
			} else {
				ips = append(ips, ip)
			}
		}
	}

	return ipBlacklist{
		filepath: filepath,
		ips:      ips,
		cidrs:    cidrs,
	}, nil
}

// IP blacklists are based on cidr ranges or specific ips
type ipBlacklist struct {
	Blacklist

	filepath string

	ips   []net.IP
	cidrs []net.IPNet
}

func (b ipBlacklist) Blocked(req http.Request) bool {
	addr := net.ParseIP(req.RemoteAddr)

	// Check specific ips
	for i := range b.ips {
		if b.ips[i].Equal(addr) {
			return true
		}
	}

	// Check ranges
	for i := range b.cidrs {
		if b.cidrs[i].Contains(addr) {
			return true
		}
	}

	return false
}
