package thttp

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Discovery defines the interface needed to support http discovery
type Discovery interface {
	URL() (string, error)
}

// DiscoveryDefault defines the default discovery using the provided host name
type DiscoveryDefault struct {
	Protocol string
	Host     string
	Port     string
}

// URL implements the Discovery interface
func (d *DiscoveryDefault) URL() (string, error) {
	return stringOr(d.Protocol, StdProtocol) + "://" + d.Host + ":" + stringOr(d.Port, StdPort), nil
}

// DiscoveryDNSA defines a discover using dns A records to round robin
type DiscoveryDNSA struct {
	Protocol string
	Host     string
	Port     string

	Resolv string

	lock       sync.Mutex
	updatelock sync.Mutex
	pos        int
	ips        []string
	expires    int64
}

func (d *DiscoveryDNSA) refresh() {

	d.updatelock.Lock()
	defer d.updatelock.Unlock()

	config, err := dns.ClientConfigFromFile(stringOr(d.Resolv, "/etc/resolv.conf"))
	if err != nil {
		fmt.Println("Could not connect to resolver", err)
		return
	}

	c := new(dns.Client)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(d.Host), dns.TypeA)
	//m.SetQuestion(dns.Fqdn("mf.strictlog.se"), dns.TypeA)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))

	if err != nil {
		fmt.Println("Could not query to resolver", err)
		return
	}

	if r == nil {
		fmt.Println("Could not query to resolver", err)
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		fmt.Println("Invalid answer from resolver", err)
		return
	}
	// Stuff must be in the answer section
	var ips []string
	var ttl uint32
	for _, rec := range r.Answer {
		if a, ok := rec.(*dns.A); ok {

			//fmt.Println(" IP", a.A.String())
			//fmt.Println(" Ttl", a.Header().Ttl)
			ips = append(ips, a.A.String())
			ttl = a.Header().Ttl
		}
	}

	if len(ips) > 0 {
		//fmt.Println("Setting", ips, len(ips))
		d.ips = ips
		d.expires = time.Now().Unix() + int64(ttl)
	}

}

// URL implements the Discovery interface
func (d *DiscoveryDNSA) URL() (string, error) {

	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.ips) == 0 || time.Now().Unix() > d.expires {
		//fmt.Println("Refreshing", d.ips)
		d.refresh()
	}

	if len(d.ips) == 0 {
		return "", errors.New("could not resolve ip for url")
	}

	return stringOr(d.Protocol, StdProtocol) + "://" + d.ips[rand.Intn(len(d.ips))] + ":" + stringOr(d.Port, StdPort), nil
}
