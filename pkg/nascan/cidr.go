package nascan

import (
	"errors"
	"math"
	"net"
	"strings"
)

func DealCIDR(cidr string) ([]string, error) {
	var ips []string
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ips, err
	}
	for ip = ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		if len(ip) > 0 && ip[len(ip)-1] == 0 || ip[len(ip)-1] == 255 {
			continue
		}
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func DealAsterisk(ip string) ([]string, error) {
	ipSplit := strings.Split(ip, ".")
	switch len(ipSplit) {
	case 2:
		return DealCIDR(strings.Join(ipSplit, ".") + ".1.1/24")
	case 3:
		return DealCIDR(strings.Join(ipSplit, ".") + ".1/16")
	case 4:
		return DealCIDR(strings.Join(ipSplit, ".") + "/32")
	default:
		return []string{}, errors.New(ip)
	}
}

func DealHyphen(ip string) ([]string, error) {
	var ips []string
	ipRange := strings.Split(ip, "-")
	//TODO 异常处理
	if len(ipRange) == 2 {
		ipStart, err := IpString2Long(strings.TrimSpace(ipRange[0]))
		if err != nil {
			return ips, err
		}
		ipEnd, err := IpString2Long(strings.TrimSpace(ipRange[1]))
		if err != nil {
			return ips, err
		}
		for i := ipStart; i <= ipEnd; i++ {
			if ipTmp, err := Long2IpString(i); err == nil {
				ips = append(ips, ipTmp)
			}
		}
		return ips, nil
	}
	return nil, errors.New("IP format error")
}

func IpString2Long(ip string) (uint, error) {
	ipTmp := net.ParseIP(ip).To4()
	if ipTmp == nil {
		return 0, errors.New("IP format error")
	}
	return uint(ipTmp[3]) | uint(ipTmp[2])<<8 | uint(ipTmp[1])<<16 | uint(ipTmp[0])<<24, nil
}

func Long2IpString(ip uint) (string, error) {
	if ip > math.MaxUint32 {
		return "", errors.New("IP format error")
	}
	ipTmp := make(net.IP, net.IPv4len)
	ipTmp[0] = byte(ip >> 24)
	ipTmp[1] = byte(ip >> 16)
	ipTmp[2] = byte(ip >> 8)
	ipTmp[3] = byte(ip)
	return ipTmp.String(), nil
}

func Handler(s string) ([]string, error) {
	ipStrings := strings.TrimSpace(s)
	var ips []string
	var err error
	if strings.Contains(ipStrings, "/") {
		//TODO 192.168.0.1/24
		ips, err = DealCIDR(ipStrings)
	} else if strings.Contains(ipStrings, "-") {
		//TODO 192.668.0.1-192.668.0.255
		ips, err = DealHyphen(ipStrings)
	} else {
		//TODO 192.168
		ips, err = DealAsterisk(ipStrings)
	}
	return ips, err
}
