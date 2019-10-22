package tcp

import (
	"time"

	"github.com/mschneider82/health"
	tcpshaker "github.com/tevino/tcp-shaker"
)

// Checker is a tcp checker that check a given addr
type Checker struct {
	Addr                string
	Timeout            time.Duration
}

// NewChecker returns a new url.Checker with the given URL
func NewChecker(addr string, timeout time.Duration) Checker {
	return Checker{Addr: addr, Timeout: timeout}
}

// Check makes a HEAD request to the given URL
// If the request returns something different than StatusOK,
// The status is set to StatusBadRequest and the URL is considered Down
func (u Checker) Check() health.Health {
	health := health.NewHealth()
	health.AddInfo("addr", u.Addr)
	t := tcpshaker.NewChecker()
	err := t.CheckAddr(u.Addr, u.Timeout)
	if err != nil {
		health.Down()
		return health
	}
	health.Up()
	return health
}
