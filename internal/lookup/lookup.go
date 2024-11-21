package lookup

import (
	"regexp"
)

type Lookup struct {
	ValidIpv4  *regexp.Regexp
	ValidIpv6  *regexp.Regexp
	ValidMac   *regexp.Regexp
	ValidEmail *regexp.Regexp
	ValidUrl   *regexp.Regexp
}

func New() *Lookup {
	return &Lookup{
		ValidIpv4:  regexp.MustCompile(Ipv4Pattern),
		ValidIpv6:  regexp.MustCompile(Ipv6Pattern),
		ValidMac:   regexp.MustCompile(MacPattern),
		ValidEmail: regexp.MustCompile(EmailPattern),
		ValidUrl:   regexp.MustCompile(UrlPattern),
	}
}
