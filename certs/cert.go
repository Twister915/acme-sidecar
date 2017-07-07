package certs

import "time"

type Cert struct {
	Key         []byte
	Certificate []byte
	Domain      string
	Expires     time.Time
}
