package main

import (
	"fmt"
	"net"
	"sort"

	"github.com/miekg/dns"
)

func queryDNS(d, s string, t uint16, a bool) (ret *dns.Msg, err error) {
	q := new(dns.Msg)
	q.SetQuestion(dns.Fqdn(d), t)

	c := new(dns.Client)
	ret, _, err = c.Exchange(q, s+":53")
	if err != nil {
		return
	}

	if ret.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("DNS server responded with code %d", ret.Rcode)
		return
	}

	if a && !ret.Authoritative {
		err = fmt.Errorf("Response is NOT authoritative")
		return
	}

	if ret.Truncated {
		err = fmt.Errorf("Response was truncated")
		return
	}

	if len(ret.Answer) == 0 {
		err = fmt.Errorf("Got no answers")
		return
	}

	return
}

func getNS(d, s string, auth bool) (ret []NS, err error) {
	m, err := queryDNS(d, s, dns.TypeNS, auth)
	if err != nil {
		return
	}

	var tmp []NS

	for k, v := range m.Answer {
		a, ok := v.(*dns.NS)
		if !ok {
			err = fmt.Errorf("%d. answer is NOT a NS record", k)
			continue
		}
		var n NS
		n, err = nsToIPs(a.Ns)
		if err != nil {
			return
		}
		tmp = append(tmp, n)
	}

	if len(tmp) == 0 {
		err = fmt.Errorf("Got no valid NS records")
		return
	}

	var keys []string
	vals := make(map[string][]Entry)
	for _, v := range tmp {
		keys = append(keys, v.Name)
		vals[v.Name] = v.List
	}
	sort.Strings(keys)

	for _, v := range keys {
		ret = append(ret, NS{Name: v, List: vals[v]})
	}

	return
}

func getSOA(d, s string, auth bool) (ret *dns.SOA, err error) {
	m, err := queryDNS(d, s, dns.TypeSOA, auth)
	if err != nil {
		return
	}

	if len(m.Answer) != 1 {
		err = fmt.Errorf("Requested one response, got %d", len(m.Answer))
		return
	}

	ret, ok := m.Answer[0].(*dns.SOA)
	if !ok {
		err = fmt.Errorf("Answer is not a SOA record")
	}

	return
}

func getSerial(d, s string) (ret string, err error) {
	soa, err := getSOA(d, s, true)
	if err != nil {
		return
	}
	ret = fmt.Sprintf("%d", soa.Serial)
	return
}

func nsToIPs(ns string) (ret NS, err error) {
	h, err := net.LookupHost(ns)
	if err != nil {
		return
	}
	ret.Name = ns

	var tmp []string

	for _, v := range h {
		a := net.ParseIP(v)
		if a == nil {
			err = fmt.Errorf("%s is NOT an IP address", v)
			return
		}
		if a.To4() != nil {
			tmp = append(tmp, v)
		}
	}

	if len(tmp) == 0 {
		err = fmt.Errorf("No valid IPv4 addresses found")
		return
	}

	sort.Strings(tmp)
	for _, v := range tmp {
		ret.List = append(ret.List, Entry{IP: v})
	}

	return
}

func getResolver() (ret string, err error) {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return
	}

	if len(conf.Servers) == 0 {
		err = fmt.Errorf("No nameserver found in /etc/resolv.conf")
		return
	}

	ret = conf.Servers[0]
	return
}
