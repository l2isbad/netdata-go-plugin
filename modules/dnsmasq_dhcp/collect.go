package dnsmasq_dhcp

import (
	"bufio"
	"io"
	"math/big"
	"net"
	"os"
	"strings"
)

func (d *DnsmasqDHCP) collect() (map[string]int64, error) {
	f, err := os.Open(d.LeasesPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	notChanged := d.modTime.Equal(fi.ModTime())
	if notChanged {
		return d.mx, nil
	}

	d.modTime = fi.ModTime()
	d.mx = d.collectRangesStats(findIPs(f))

	return d.mx, nil
}

func (d *DnsmasqDHCP) collectRangesStats(ips []net.IP) map[string]int64 {
	mx := make(map[string]int64)

	for _, ip := range ips {
		for _, r := range d.ranges {
			if !r.Contains(ip) {
				continue
			}
			mx[r.String()]++
			break
		}
	}

	for _, r := range d.ranges {
		name := r.String()
		numOfIps, ok := mx[name]
		if !ok {
			mx[name] = 0
		}

		hosts := r.Hosts()
		if !hosts.IsInt64() {
			continue
		}

		mx[name+"_percent"] = int64(calcPercent(numOfIps, hosts) * 1000)
	}

	return mx
}

func findIPs(r io.Reader) []net.IP {
	/*
		1560300536 08:00:27:61:3c:ee 2.2.2.3 debian8 *
		duid 00:01:00:01:24:90:cf:5b:08:00:27:61:2e:2c
		1560300414 660684014 1234::20b * 00:01:00:01:24:90:cf:a3:08:00:27:61:3c:ee
	*/
	var ips []net.IP
	s := bufio.NewScanner(r)

	for s.Scan() {
		parts := strings.Fields(s.Text())
		if len(parts) != 5 {
			continue
		}

		ip := net.ParseIP(parts[2])
		if ip == nil {
			continue
		}
		ips = append(ips, ip)
	}

	return ips
}

func calcPercent(ips int64, hosts *big.Int) float64 {
	if ips == 0 || hosts.Int64() == 0 {
		return 0
	}
	return float64(ips) * 100 / float64(hosts.Int64())
}
