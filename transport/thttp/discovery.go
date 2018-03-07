package thttp

import (
	"net"
	"github.com/miekg/dns"
	"fmt"
	"sync"
	"time"
	"math/rand"
)

type Discovery interface {
	Url() (string, error)
}

type DiscoveryDefault struct {
	Protocol string
	Host     string
	Port     string
}

func (d *DiscoveryDefault) Url() (string, error) {
	return stringOr(d.Protocol, STD_PROTOCOL) + "://" + d.Host + ":" + stringOr(d.Port, STD_PORT), nil
}

type DiscoveryDnsA struct {
	Protocol string
	Host     string
	Port     string

	Resolv   string

	lock     sync.Mutex
	pos      int
	ips      []string
	expires  int64
}

func (d *DiscoveryDnsA) refresh() {

	d.lock.Lock()
	defer func() {
		d.lock.Unlock()
	}()

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
	ips := []string{}
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
		d.ips = ips
		d.expires = time.Now().Unix() + int64(ttl)
	}

}

func (d *DiscoveryDnsA) Url() (string, error) {

	if len(d.ips) == 0 || time.Now().Unix() > d.expires{
		d.refresh()
	}


	return stringOr(d.Protocol, STD_PROTOCOL) + "://" + d.ips[rand.Intn(len(d.ips))] + ":" + stringOr(d.Port, STD_PORT), nil
}

