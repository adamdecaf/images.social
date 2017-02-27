package blacklist

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"path"
	"os"
)

// blacklist based on certain conditions
// - ip, read from file every N minutes
// - upload frequency, each upload is recorded

// Blacklist inspects a request to determine if the request needs to be blocked.
// It is expected that the `Blocked` method WILL NOT modify the request in any way.
type Blacklist interface {
	Blocked(http.Request) bool
}

// New creates a new blacklist object by reading from a filepath, right now
// the format is simple in that it's matching on blocked ip addresses.
func New(filepath string) (Blacklist, error) {
	filepath = expandPath(filepath)
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

func expandPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path.Join(dir, p)
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
