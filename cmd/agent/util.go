package main

import (
	"net"
)

func resolveIP(address string) (net.IP, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err2 := net.LookupIP(host)
		if err2 != nil {
			return nil, err
		}
		ip = ips[0]
	}
	return ip, nil
}
