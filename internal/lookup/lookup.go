package lookup

import "regexp"

type Lookup struct {
	ValidIp *regexp.Regexp
}

func New() *Lookup {
	return &Lookup{
		ValidIp: regexp.MustCompile(`((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}`),
	}
}
